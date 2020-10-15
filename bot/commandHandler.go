package bot

import (
	"log"
	"medgebot/internal/pkg/irc"
	"strings"
)

func (bot *Bot) HandleCommands() {
	bot.RegisterHandler(func(msg irc.Message) {
		if msg.Command == "PRIVMSG" {
			contents := strings.TrimPrefix(strings.Join(msg.Params[1:], " "), ":")

			// Command processing
			// TODO make better
			if strings.HasPrefix(contents, "!hello") {
				if err := bot.SendMessage("WORLD"); err != nil {
					log.Printf("ERROR: send failed: %s", err)
				}
			}

			// Sorcery Shoutout
			if strings.HasPrefix(contents, "!sorcery") {
				if err := bot.SendMessage("!so @SorceryAndSarcasm"); err != nil {
					log.Printf("ERROR: send failed: %s", err)
				}
			}
		}
	})
}
