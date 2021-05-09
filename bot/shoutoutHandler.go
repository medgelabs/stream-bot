package bot

import (
	"fmt"
	"strings"
)

// HandleShoutoutCommand responds to the !so chat command
func (bot *Bot) HandleShoutoutCommand() {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsChatEvent() {
				contents := evt.Message

				if !strings.HasPrefix(contents, "!lso") {
					return
				}

				tokens := strings.Split(contents, " ")
				broadcaster := strings.TrimPrefix(tokens[1], "@")
				bot.SendMessage(
					fmt.Sprintf("Go check out @%s at https://twitch.tv/%s!", broadcaster, broadcaster),
				)

				// Grab channel being shouted out
				// Call Twitch API for URL
				/*
					curl -X GET 'https://api.twitch.tv/helix/users?login_name=camikazeey' \
					-H 'Authorization: Bearer FAKE' \
					-H 'Client-Id: FAKE'

					// id == broadcaster_id?
				*/
				// Stretch: Add last playing
				/*
					curl -X GET 'https://api.twitch.tv/helix/channels?broadcaster_id=44445592' \
					-H 'Authorization: Bearer FAKE' \
					-H 'Client-Id: FAKE'

					// game_name
				*/

			}
		}),
	)
}
