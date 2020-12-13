package cmd

import (
	"fmt"

	"github.com/go-shell/grumble"
)

func init() {
	adminCommand := &grumble.Command{
		Name:     "admin",
		Help:     "admin tools",
		LongHelp: "super administration tools",
	}
	App.AddCommand(App.UserNode(), adminCommand)

	adminCommand.AddCommand(&grumble.Command{
		Name: "root",
		Help: "root the machine",
		Run: func(c *grumble.Context) error {
			fmt.Println(c.Flags.String("directory"))
			return fmt.Errorf("failed")
		},
	})

	adminCommand.AddCommand(&grumble.Command{
		Name: "kill",
		Help: "kill the process",
		Run: func(c *grumble.Context) error {
			return fmt.Errorf("failed")
		},
	})
}
