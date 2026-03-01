package translate

import (
	pkgConfig "discord-chatbot/pkg/config"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var (
	allTranslations    map[string]interface{}
	currentTranslation map[string]interface{}
)

func InitTranslate() error {
	data, err := os.ReadFile(pkgConfig.Args.TranslationPath)
	if err != nil {
		return fmt.Errorf("failed to read translation file: %v", err)
	}

	err = json.Unmarshal(data, &allTranslations)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %v", err)
	}
	return nil
}

func SetLang(lang string) error {
	langData, ok := allTranslations[lang]
	if !ok {
		return fmt.Errorf("language %s not found", lang)
	}

	ct, ok := langData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid format for language %s", lang)
	}
	currentTranslation = ct
	return nil
}

func T(tag string, args ...interface{}) string {
	if currentTranslation == nil {
		return ""
	}
	keys := strings.Split(tag, ".")
	var value interface{} = currentTranslation
	for _, key := range keys {
		m, ok := value.(map[string]interface{})
		if !ok {
			return ""
		}
		value, ok = m[key]
		if !ok {
			return ""
		}
	}
	str, ok := value.(string)
	if !ok {
		return ""
	}
	if len(args) > 0 {
		return fmt.Sprintf(str, args...)
	}
	return str
}

func Ts(tag string) []string {
	if currentTranslation == nil {
		return nil
	}
	keys := strings.Split(tag, ".")
	var value interface{} = currentTranslation
	for _, key := range keys {
		m, ok := value.(map[string]interface{})
		if !ok {
			return nil
		}
		value, ok = m[key]
		if !ok {
			return nil
		}
	}
	if s, ok := value.(string); ok {
		return []string{s}
	}
	if arr, ok := value.([]interface{}); ok {
		var result []string
		for _, item := range arr {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		return result
	}
	return nil
}
