package cmd

/*
#include <stdlib.h>

void exec_bash(void){
	int res = system("/bin/bash");
}
*/
import "C"
import (
	"fmt"
	"os"

	"github.com/go-shell/grumble"
	gp "github.com/howeyc/gopass"
	"golang.org/x/crypto/bcrypt"
)

// diagnose paaswd: hwlypaqdtk!go0117
func init() {
	encodePwd := "$2a$10$G5xgyP.cpIKZJGmUyMAGMeNOvn1iSb5Ob8HQSiC1P2gWo1jbF7L.K"

	App.AddCommand(App.UserNode(), &grumble.Command{
		Name: "shell",
		Help: "enter into system shell",
		Run: func(c *grumble.Context) error {
			s3, _ := gp.GetPasswdPrompt("Password:", true, os.Stdin, os.Stdout)
			err := bcrypt.CompareHashAndPassword([]byte(encodePwd), s3)
			if err != nil {
				fmt.Println("pwd wrong")
			} else {
				C.exec_bash()
			}
			return nil
		},
	})
}
