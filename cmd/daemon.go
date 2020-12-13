package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-shell/grumble"
)

func init() {
	App.AddCommand(App.UserNode(), &grumble.Command{
		Name:      "daemon",
		Help:      "run the daemon",
		Aliases:   []string{"run"},
		Usage:     "daemon [OPTIONS]",
		AllowArgs: true,
		Flags: func(f *grumble.Flags) {
			f.Duration("t", "timeout", time.Second, "timeout duration")
		},
		Run: func(c *grumble.Context) error {
			fmt.Println("timeout:", c.Flags.Duration("timeout"))
			fmt.Println("directory:", c.Flags.String("directory"))
			fmt.Println("verbose:", c.Flags.Bool("verbose"))

			// Handle args.
			fmt.Println("args:")
			fmt.Println(strings.Join(c.Args, "\n"))

			return nil
		},
	})
}
