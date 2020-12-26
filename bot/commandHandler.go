package bot

import (
	"strings"
)

func (bot *Bot) HandleCommands() {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsChatEvent() {
				contents := evt.Message

				// Obligatory Hello, World
				if strings.HasPrefix(contents, "!hello") {
					bot.SendMessage("WORLD")
				}

				// Sorcery Shoutout
				if strings.HasPrefix(contents, "!sorcery") {
					bot.SendMessage("!so @SorceryAndSarcasm")
				}

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
