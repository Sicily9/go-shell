package cmd

import (
	"github.com/go-shell/grumble"
)

func init() {
	promptCommand := &grumble.Command{
		Name: "prompt",
		Help: "set a custom prompt",
	}
	App.AddCommand(App.UserNode(), promptCommand)

	promptCommand.AddCommand(&grumble.Command{
		Name: "set",
		Help: "set a custom prompt",
		Run: func(c *grumble.Context) error {
			c.App.SetPrompt("CUSTOM PROMPT >> ")
			return nil
		},
	})

	promptCommand.AddCommand(&grumble.Command{
		Name: "reset",
		Help: "reset to default prompt",
		Run: func(c *grumble.Context) error {
			c.App.SetDefaultPrompt()
			return nil
		},
	})
}
