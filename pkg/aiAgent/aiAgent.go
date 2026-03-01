package aiAgent

import (
	"context"
	pkgConfig "discord-chatbot/pkg/config"
	"fmt"
	"slices"
	"sync"

	"github.com/sashabaranov/go-openai"
)

// Provider constants
const (
	ProviderOpenAI   = "openai"
	ProviderDeepSeek = "deepseek"
	ProviderGemini   = "gemini"
	ProviderGrok     = "grok"
)

// Provider base URLs
var providerBaseURLs = map[string]string{
	ProviderOpenAI:   "",
	ProviderDeepSeek: "https://api.deepseek.com/v1",
	ProviderGemini:   "https://generativelanguage.googleapis.com/v1beta/openai/",
	ProviderGrok:     "https://api.x.ai/v1",
}

// Provider supported models
var providerModels = map[string][]string{
	ProviderOpenAI: {
		"gpt-5", "gpt-5-mini", "gpt-5-nano",
		"gpt-5.2", "gpt-5-chat-latest", "gpt-5.2-pro",
		"gpt-4o", "gpt-4o-mini", "gpt-4o-nano",
		"o3", "o4-mini", "o1", "o1-mini",
	},
	ProviderDeepSeek: {
		"deepseek-chat", "deepseek-reasoner",
	},
	ProviderGemini: {
		"gemini-2.0-flash-exp", "gemini-1.5-flash", "gemini-1.5-pro",
		"gemini-1.0-pro",
	},
	ProviderGrok: {
		"grok-2", "grok-2-mini", "grok-beta",
	},
}

const MAX_QUEUE_LENGTH = 100

type AiBot struct {
	client      *openai.Client
	provider    string
	model       string
	queueLength int
	MessageDb   *MessageDataBase
	tools       []openai.Tool
	lock        sync.Mutex
}

type MessageDataBase struct {
	systemMessages []openai.ChatCompletionMessage
	messages       []openai.ChatCompletionMessage
}

func NewAiBot(config pkgConfig.AiAgentConfig, apiKey string) (*AiBot, error) {
	provider := config.Provider
	model := config.Model

	if !IsSupportedProvider(provider) {
		return nil, fmt.Errorf("unsupported provider: %s (supported: %v)", provider, GetSupportedProviders())
	}
	if !IsSupportedModel(provider, model) {
		return nil, fmt.Errorf("unsupported model: %s for provider %s (supported: %v)", model, provider, GetSupportedModels(provider))
	}

	cfg := openai.DefaultConfig(apiKey)
	if baseURL, ok := providerBaseURLs[provider]; ok && baseURL != "" {
		cfg.BaseURL = baseURL
	}

	bot := &AiBot{
		client:      openai.NewClientWithConfig(cfg),
		provider:    provider,
		model:       model,
		queueLength: config.QueueLength,
		MessageDb:   &MessageDataBase{},
		tools:       []openai.Tool{},
		lock:        sync.Mutex{},
	}
	return bot, nil
}

func (m *MessageDataBase) AddSystemMessages(messages []string) {
	for _, message := range messages {
		m.systemMessages = append(m.systemMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: message,
		})
	}
}

func (m *MessageDataBase) AddUserMessages(messages []string) {
	for _, message := range messages {
		m.messages = append(m.messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: message,
		})
	}
}

func command(aiBot *AiBot, messageDb *MessageDataBase, message string) ([]openai.ChatCompletionMessage, error) {
	if aiBot.client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	if !aiBot.lock.TryLock() {
		return nil, fmt.Errorf("ai is busy")
	}
	defer aiBot.lock.Unlock()

	// Add user message
	userMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	}
	messageDb.messages = append(messageDb.messages, userMessage)

	// Build messages list
	var messages []openai.ChatCompletionMessage
	messages = append(messages, aiBot.MessageDb.systemMessages...)
	messages = append(messages, messageDb.systemMessages...)
	messages = append(messages, messageDb.messages...)

	// Build request
	req := openai.ChatCompletionRequest{
		Model:    aiBot.model,
		Messages: messages,
	}
	if len(aiBot.tools) > 0 {
		req.Tools = aiBot.tools
	}

	// Call API
	resp, err := aiBot.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return nil, err
	}

	// Process response
	var aiResponse []openai.ChatCompletionMessage
	for _, choice := range resp.Choices {
		responseMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: choice.Message.Content,
		}
		if len(choice.Message.ToolCalls) > 0 {
			responseMessage.ToolCalls = choice.Message.ToolCalls
		}
		aiResponse = append(aiResponse, responseMessage)
		messageDb.messages = append(messageDb.messages, responseMessage)
	}

	messageDb.messages = dequeueMessages(messageDb.messages, aiBot.queueLength)
	return aiResponse, nil
}

func (a *AiBot) Command(message string) ([]openai.ChatCompletionMessage, error) {
	return command(a, a.MessageDb, message)
}

func (a *AiBot) CommandWithDatabase(messageDb *MessageDataBase, message string) ([]openai.ChatCompletionMessage, error) {
	return command(a, messageDb, message)
}

func (a *AiBot) SetQueueLength(length int) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if length < 0 {
		length = 0
	} else if length > MAX_QUEUE_LENGTH {
		length = MAX_QUEUE_LENGTH
	}
	a.queueLength = length
}

// IsSupportedProvider checks if a provider is supported
func IsSupportedProvider(provider string) bool {
	_, ok := providerBaseURLs[provider]
	return ok
}

// GetSupportedProviders returns a list of supported providers
func GetSupportedProviders() []string {
	providers := make([]string, 0, len(providerBaseURLs))
	for p := range providerBaseURLs {
		providers = append(providers, p)
	}
	return providers
}

// IsSupportedModel checks if a model is supported by the given provider
func IsSupportedModel(provider, model string) bool {
	models, ok := providerModels[provider]
	if !ok {
		return false
	}
	return slices.Contains(models, model)
}

// GetSupportedModels returns a list of supported models for the given provider
func GetSupportedModels(provider string) []string {
	if models, ok := providerModels[provider]; ok {
		return models
	}
	return nil
}

// SetAiModel sets the provider and model, returns error if provider or model is not supported
func (a *AiBot) SetAiModel(provider string, model string, apiKey string) error {
	if !IsSupportedProvider(provider) {
		return fmt.Errorf("unsupported provider: %s (supported: %v)", provider, GetSupportedProviders())
	}
	if !IsSupportedModel(provider, model) {
		return fmt.Errorf("unsupported model: %s for provider %s (supported: %v)", model, provider, GetSupportedModels(provider))
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	// Rebuild client if provider changed
	if a.provider != provider {
		cfg := openai.DefaultConfig(apiKey)
		if baseURL, ok := providerBaseURLs[provider]; ok && baseURL != "" {
			cfg.BaseURL = baseURL
		}
		a.client = openai.NewClientWithConfig(cfg)
		a.provider = provider
	}

	a.model = model
	return nil
}

func dequeueMessages(messages []openai.ChatCompletionMessage, length int) []openai.ChatCompletionMessage {
	if length > 0 && len(messages) > length {
		startIndex := len(messages) - length
		messages = messages[startIndex:]
	}
	return messages
}
