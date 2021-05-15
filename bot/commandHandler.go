package bot

import (
	"fmt"
	"math/rand"
	log "medgebot/logger"
	"strings"
	"time"
)

// Command represents a known Command from Config that the Bot can respond to
type Command struct {
	Prefix          string
	IsAlias         bool
	AliasFor        string
	MessageTemplate HandlerTemplate
}

// ParsedMessage Return the interpolated Message for the given command
func (c *Command) ParsedMessage(evt Event) string {
	return c.MessageTemplate.Parse(evt)
}

// HandleCommands is the chat component of KnownCommands
func (bot *Bot) HandleCommands(knownCommands []Command) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsChatEvent() {
				contents := evt.Message

				log.Info("Checking for commands in: %+v", evt)

				// 1: check if the chat event is a Command message
				// 2: Is the command an Alias? If so - alter contents and send back through the bot
				// 3: Otherwise - map command.Prefix to desired message contents

				// For commands that are simple message responders
				for _, command := range knownCommands {
					if strings.HasPrefix(contents, command.Prefix) {

						// If the Command is an alias for another command, change message contents and send back to the Bot
						if command.IsAlias {
							evt.Message = command.AliasFor
							bot.receiveEvent(evt)
							break
						}

						// Otherwise, if it's a known simple Message Command
						bot.SendMessage(command.ParsedMessage(evt))
					}
				}

				// Derived commands lists
				if strings.HasPrefix(contents, "!commands") {
					var buf strings.Builder
					buf.WriteString("Commands: ")
					for _, command := range knownCommands {
						buf.WriteString(command.Prefix)
						buf.WriteString(" ")
					}

					buf.WriteString("!cthulhu")
					buf.WriteString(" ")

					buf.WriteString("!coin")
					buf.WriteString(" ")

					bot.SendMessage(buf.String())
				}

				// Special case because fitting this in config.yaml is :spooky127Concern:
				// Fjoell Feature Request: ASCII Cthulu
				if strings.HasPrefix(contents, "!cthulhu") {
					msg := `⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
							⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠋⠉⠉⠉⠙⢿⣷⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
							⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠏⠀⠀⠀⠀⠀⠀⠀⢹⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
							⣿⣿⣿⣿⣿⠿⠿⢿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⢻⣿⣿⣿⣿⣿⣿⣿⣿⣿
							⣿⣿⣿⡿⠋⣴⣴⣸⢿⠋⠈⠂⠀⠀⠀⠀⠀⠀⠀⠀⠈⢻⣿⣿⣿⣿⣿⣿⣿⣿
							⣿⣿⣿⡇⠀⣾⣿⣿⣾⣇⠀⠠⠄⠀⢀⢀⠀⠀⢠⢀⠀⣼⣿⣏⣀⡈⢻⣿⣿⣿
							⣿⣿⣿⣷⡀⠈⠙⠛⠙⠋⠀⠀⠀⠀⠀⠈⠀⠀⠀⠀⠘⠛⠿⠿⠟⠃⢸⣿⣿⣿
							⣿⣿⣿⣿⣿⣶⡄⠂⡒⣂⡨⠃⠀⠀⠈⠈⠀⠀⠈⢐⠒⠀⠠⠤⠤⡞⣿⣿⣿⣿
							⣿⣿⣿⣿⣿⣿⠀⠖⠀⠀⢀⡠⠊⠀⢠⣶⠀⠈⠢⣀⡀⠀⠑⠲⠀⣏⣿⣿⣿⣿
							⣿⣿⣿⣿⣿⣿⣶⣦⣀⡁⠜⢁⢴⠀⢀⣷⣧⡀⠱⡈⠂⠓⠀⢰⣾⣿⣿⣿⣿⣿
							⣿⣿⣿⣿⣿⣿⣿⣿⣿⡁⠈⠕⠘⠀⠘⠿⠿⠇⢠⠿⠀⣶⣾⣿⣿⣿⣿⣿⣿⣿
							⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣄⣁⣀⣆⡐⣶⣶⣧⣴⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿
							⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿`
					bot.SendMessage(strings.TrimSpace(msg))
				}

				// Fjoell Feature Request: Coin Throw
				if strings.HasPrefix(contents, "!coin") {
					rand.Seed(time.Now().UnixNano())
					side := 1 + rand.Int()%2
					result := ""
					if side == 1 {
						result = "heads"
					} else {
						result = "tails"
					}

					bot.SendMessage(fmt.Sprintf("@%s flipped: %s", evt.Sender, result))
				}
			}
		}),
	)
}
