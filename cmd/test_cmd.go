package cmd

import (
	"fmt"

	"github.com/go-shell/grumble"
)

// diagnose paaswd: hwlypaqdtk!go0117
func init() {

	App.AddCommand(App.TestNode(), &grumble.Command{
		Name: "test_cmd",
		Help: "test_cmd print",
		Run: func(c *grumble.Context) error {
			fmt.Println("test_cmd")
			return nil
		},
	})
}
