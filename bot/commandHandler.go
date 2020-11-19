package bot

import (
	"strings"
)

func (bot *Bot) HandleCommands() {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsChatEvent() {
				contents := evt.Message

				// Command processing
				if strings.HasPrefix(contents, "!hello") {
					bot.SendMessage("WORLD")
				}

				// Sorcery Shoutout
				if strings.HasPrefix(contents, "!sorcery") {
					bot.SendMessage("!so @SorceryAndSarcasm")
				}
			}
		}),
	)
}
