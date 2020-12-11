package bot

import (
	"fmt"
	"time"
)

func (bot *Bot) RegisterRaidHandler(messageFormat string) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsRaidEvent() {
				time.Sleep(3 * time.Second)
				bot.SendMessage(fmt.Sprintf(messageFormat, evt.Sender))
				// log.Println(fmt.Sprintf(messageFormat, evt.Sender))
			}
		}),
	)
}
