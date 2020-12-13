module go-shell

go 1.15

require github.com/go-shell/cmd v0.0.0

require github.com/go-shell/grumble v0.0.0

require github.com/go-shell/readline v0.0.0 // indirect

require (
	github.com/chzyer/logex v1.1.10 // indirect
	github.com/chzyer/test v0.0.0-20180213035817-a1ea475d72b1 // indirect
	github.com/desertbit/closer/v3 v3.1.2 // indirect
	github.com/desertbit/columnize v2.1.0+incompatible // indirect
	github.com/desertbit/go-shlex v0.1.1 // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/howeyc/gopass v0.0.0-20190910152052-7cb4b85ec19c // indirect
	golang.org/x/crypto v0.0.0-20201208171446-5f87f3452ae9 // indirect
	gopkg.in/AlecAivazis/survey.v1 v1.8.8 // indirect
)

replace github.com/go-shell/cmd => ./cmd

replace github.com/go-shell/grumble => ./grumble

replace github.com/go-shell/readline => ./readline
