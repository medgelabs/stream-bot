package bot

import (
	"log"
	"strings"
	"text/template"
)

func (bot *Bot) RegisterBitsHandler(messageTemplate *template.Template) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsBitsEvent() {
				log.Printf("> %s cheered %d bits!", evt.Sender, evt.Amount)

				var msg strings.Builder
				err := messageTemplate.Execute(&msg, evt)
				if err != nil {
					log.Printf("ERROR: bits template execute - %v", err)
					return
				}

				bot.SendMessage(msg.String())
			}
		}),
	)
}
