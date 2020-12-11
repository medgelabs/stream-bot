package bot

import (
	"fmt"
	"time"
)

func (bot *Bot) RegisterSubsHandler(subMessgeFormat, giftSubMessageFormat string) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			var msg string

			if evt.IsSubEvent() {
				msg = fmt.Sprintf(subMessgeFormat, evt.Amount, evt.Sender)
			} else if evt.IsGiftSubEvent() {
				msg = fmt.Sprintf(giftSubMessageFormat, evt.Sender, evt.Recipient)
			} else {
				return // no messaging otherwise
			}

			time.Sleep(2 * time.Second)
			bot.SendMessage(msg)
		}),
	)
}
