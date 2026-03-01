package discordlogger

import (
	pkgConfig "discord-chatbot/pkg/config"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// DiscordHook is a logrus hook that sends error-level logs to Discord channels
type DiscordHook struct {
	session    *discordgo.Session
	channelIDs []string
	sending    atomic.Bool
}

// NewDiscordHook creates a new Discord logrus hook
func NewDiscordHook(session *discordgo.Session, logChannels []pkgConfig.LogChannelConfig) *DiscordHook {
	channelIDs := make([]string, len(logChannels))
	for i, ch := range logChannels {
		channelIDs[i] = ch.Id
	}
	return &DiscordHook{
		session:    session,
		channelIDs: channelIDs,
	}
}

// Levels returns the log levels this hook applies to
func (h *DiscordHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

// Fire sends the log entry to all configured Discord channels
func (h *DiscordHook) Fire(entry *logrus.Entry) error {
	// Anti-recursion guard: if ChannelMessageSend itself triggers an error log, skip
	if !h.sending.CompareAndSwap(false, true) {
		return nil
	}
	defer h.sending.Store(false)

	msg := formatDiscordMessage(entry)

	for _, channelID := range h.channelIDs {
		h.session.ChannelMessageSend(channelID, msg)
	}

	return nil
}

func formatDiscordMessage(entry *logrus.Entry) string {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())

	packageName := ""
	if pkg, ok := entry.Data["package"]; ok {
		packageName = fmt.Sprintf(" [%v]", pkg)
	}

	caller := ""
	if file, ok := entry.Data["file"]; ok {
		if line, ok := entry.Data["line"]; ok {
			if fn, ok := entry.Data["function"]; ok {
				caller = fmt.Sprintf(" %v:%v %v", file, line, fn)
			}
		}
	}

	msg := fmt.Sprintf("```\n[%s] [%-5s]%s%s\n%s\n```",
		timestamp,
		level,
		packageName,
		caller,
		entry.Message,
	)

	// Discord message limit is 2000 characters
	if len(msg) > 2000 {
		msg = msg[:1993] + "\n```"
	}

	return msg
}
