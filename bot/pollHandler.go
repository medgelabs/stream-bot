package bot

import (
	"medgebot/logger"
	"strconv"
	"strings"
)

// RegisterPollHandler collects Poll answers from Chat messages
func (bot *Bot) RegisterPollHandler() {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if strings.HasPrefix(evt.Message, "!poll") {
				bot.SendPollMessage()
			}

			if !bot.IsPollRunning() {
				return
			}

			// Message should just be a number (and within range of answers). Otherwise reject as a vote
			vote, err := strconv.Atoi(evt.Message)
			if err != nil {
				return // Assumed to not be a valid vote
			}

			alreadyVoted, err := bot.dataStore.GetOrDefault("voters", "")
			if err != nil {
				logger.Error(err, "fetch voters from bot.dataStore")
				return
			}

			if strings.Contains(alreadyVoted, evt.Sender) {
				return // Can't vote multiple times
			}

			// Valid vote - append their vote and note that they voted
			bot.dataStore.Append("voters", ",", evt.Sender)
			bot.AddPollVote(vote)
		}),
	)
}
