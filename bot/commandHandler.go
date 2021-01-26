package bot

import (
	"strings"
)

// Command represents a known Command from Config that the Bot can respond to
type Command struct {
	Prefix          string
	MessageTemplate HandlerTemplate
}

// Return the interpolated Message for the given command
func (c *Command) ParsedMessage(evt Event) string {
	return c.MessageTemplate.Parse(evt)
}

func (bot *Bot) HandleCommands(knownCommands []Command) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsChatEvent() {
				contents := evt.Message

				for _, command := range knownCommands {
					if strings.HasPrefix(contents, command.Prefix) {
						bot.SendMessage(command.ParsedMessage(evt))
					}
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
			}
		}),
	)
}
