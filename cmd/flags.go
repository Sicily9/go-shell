package cmd

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/go-shell/grumble"
)

func init() {
	App.AddCommand(App.UserNode(), &grumble.Command{
		Name: "flags",
		Help: "test flags",
		Flags: func(f *grumble.Flags) {
			f.Duration("d", "dur", time.Second, "duration test")
			f.Int("i", "int", 1, "test int")
			f.Int64("l", "int64", 2, "test int64")
			f.Uint("u", "uint", 3, "test uint")
			f.Uint64("j", "uint64", 4, "test uint64")
			f.Float64("f", "float", 5.55, "test float64")
		},
		Run: func(c *grumble.Context) error {
			exec.Command("cat", "/etc/shells")
			// 执行命令，并返回结果
			fmt.Println("duration ", c.Flags.Duration("dur"))
			fmt.Println("int      ", c.Flags.Int("int"))
			fmt.Println("int64    ", c.Flags.Int64("int64"))
			fmt.Println("uint     ", c.Flags.Uint("uint"))
			fmt.Println("uint64   ", c.Flags.Uint64("uint64"))
			fmt.Println("float    ", c.Flags.Float64("float"))
			return nil
		},
	})
}
