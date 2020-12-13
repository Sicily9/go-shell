package cmd

import (
	"fmt"
	"os"

	"github.com/go-shell/grumble"
	gp "github.com/howeyc/gopass"
	"golang.org/x/crypto/bcrypt"
)

// test passwd: venus.test
func init() {
	encodePwd := "$2a$10$vphNZMGoKmOG/bFW8vGM9OOhKLuNT4GuwEvONexb7.4DRfmrdc75u"

	App.AddCommand(App.UserNode(), &grumble.Command{
		Name: "test",
		Help: "enter into test view",
		Run: func(c *grumble.Context) error {
			s3, _ := gp.GetPasswdPrompt("Password:", true, os.Stdin, os.Stdout)
			err := bcrypt.CompareHashAndPassword([]byte(encodePwd), s3)
			if err != nil {
				fmt.Println("pwd wrong")
			} else {
				App.SetPrompt("test > ")
				App.SetCurrentView(App.TestNode())
				App.SetAutoCompleter(App.CurrentCommands())
			}
			return nil
		},
	})
}
