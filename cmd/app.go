package cmd

import (
	"github.com/fatih/color"
	"github.com/go-shell/grumble"
)

var App = grumble.New(&grumble.Config{
	Name:                  "venus",
	Description:           "An awesome venus bar",
	HistoryFile:           "/tmp/venus_bar.hist",
	Prompt:                "venus > ",
	PromptColor:           color.New(color.FgGreen, color.Bold),
	HelpHeadlineColor:     color.New(color.FgGreen),
	HelpHeadlineUnderline: true,
	HelpSubCommands:       true,

	Flags: func(f *grumble.Flags) {
		f.String("d", "directory", "DEFAULT", "set an alternative root directory path")
		f.Bool("v", "verbose", false, "enable verbose mode")
	},
})

func init() {
	App.SetPrintASCIILogo(func(a *grumble.App) {
		a.Println(" _          ______________________       __________")
		a.Println("\\ \\        /    _______|   _____  |     |  ________|")
		a.Println(" \\ \\      / /| |_______   |     | |     | |________ ")
		a.Println("  \\ \\    / / |  _______|  |     | |     |  ______  |   ")
		a.Println("   \\ \\__/ /  | |_______   |     | |_____| |______| | ")
		a.Println("    \\____/   |_________|__|     |__________________|  ")
		a.Println("")
	})
	//App.SetPrintASCIILogo(func(a *grumble.App) {
	//		a.Println("                   _   _     ")
	//		a.Println(" ___ ___ _ _ _____| |_| |___ ")
	//		a.Println("| . |  _| | |     | . | | -_|")
	//		a.Println("|_  |_| |___|_|_|_|___|_|___|")
	//		a.Println("|___|                        ")
	//		a.Println()
	//	})
}
