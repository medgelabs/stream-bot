package bot

import (
	"log"
	"medgebot/irc"
	"strings"
)

func (bot *Bot) HandleCommands() {
	bot.RegisterHandler(func(msg irc.Message) {
		if msg.Command == "PRIVMSG" {
			channel := msg.Params[0]
			contents := strings.TrimPrefix(strings.Join(msg.Params[1:], " "), ":")

			// Command processing
			// TODO make better
			if strings.HasPrefix(contents, "!hello") {
				if err := bot.client.PrivMsg(channel, "WORLD"); err != nil {
					log.Printf("ERROR: send failed: %s", err)
				}
			}

			// Sorcery Shoutout
			if strings.HasPrefix(contents, "!sorcery") {
				if err := bot.client.PrivMsg(channel, "!so @SorceryAndSarcasm"); err != nil {
					log.Printf("ERROR: send failed: %s", err)
				}
			}
		}
	})
}
