/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2018 Roland Singer [roland.singer@deserbit.com]
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package grumble

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/desertbit/closer/v3"
	shlex "github.com/desertbit/go-shlex"
	"github.com/fatih/color"
	"github.com/go-shell/readline"
	gp "github.com/howeyc/gopass"
)

// 表示命令行节点类型
const (
	AuthNODE int32 = 0
	UserNODE int32 = 1
	TestNODE int32 = 2
)

// AuthNode return enum AuthNODE
func (a *App) AuthNode() int32 {
	return AuthNODE
}

// UserNode return enum UserNODE
func (a *App) UserNode() int32 {
	return UserNODE
}

// TestNode return enum TestNODE
func (a *App) TestNode() int32 {
	return TestNODE
}

// App is the entrypoint.
type App struct {
	closer.Closer

	rl            *readline.Instance
	config        *Config
	commandsMap   map[int32]*Commands // 每个视图有固定的命令
	isShell       bool
	currentPrompt string
	currentView   int32

	flags   Flags
	flagMap FlagMap

	initHook  func(a *App, flags FlagMap) error
	shellHook func(a *App) error

	printHelp        func(a *App, shell bool)
	printCommandHelp func(a *App, cmd *Command, shell bool)
	interruptHandler func(a *App, count int)
	printASCIILogo   func(a *App)
}

// New creates a new app.
// Panics if the config is invalid.
func New(c *Config) (a *App) {
	// Prepare the config.
	c.SetDefaults()
	err := c.Validate()
	if err != nil {
		panic(err)
	}

	// APP.
	a = &App{
		Closer:           closer.New(),
		config:           c,
		currentPrompt:    c.prompt(),
		commandsMap:      make(map[int32]*Commands),
		currentView:      AuthNODE,
		flagMap:          make(FlagMap),
		printHelp:        defaultPrintHelp,
		printCommandHelp: defaultPrintCommandHelp,
		interruptHandler: defaultInterruptHandler,
	}

	a.commandsMap[AuthNODE] = &Commands{}
	a.commandsMap[UserNODE] = &Commands{}
	a.commandsMap[TestNODE] = &Commands{}

	// Register the builtin flags.
	a.flags.Bool("h", "help", false, "display help")
	a.flags.BoolL("nocolor", false, "disable color output")

	// Register the user flags if present.
	if c.Flags != nil {
		c.Flags(&a.flags)
	}

	return
}

// SetPrompt a new prompt.
func (a *App) SetPrompt(p string) {
	if !a.config.NoColor {
		p = a.config.PromptColor.Sprint(p)
	}
	a.currentPrompt = p
}

// SetDefaultPrompt resets the current prompt to the default prompt as
// configured in the config.
func (a *App) SetDefaultPrompt() {
	a.currentPrompt = a.config.prompt()
}

// IsShell indicates, if this is a shell session.
func (a *App) IsShell() bool {
	return a.isShell
}

// SetCurrentView set the app's current view
func (a *App) SetCurrentView(val int32) {
	a.currentView = val
}

// CurrentView returns the app's current view
func (a *App) CurrentView() int32 {
	return a.currentView
}

// Config returns the app's config value.
func (a *App) Config() *Config {
	return a.config
}

// SetAutoCompleter set the app's AutoCompleter.
func (a *App) SetAutoCompleter(commands *Commands) {
	completer := &completer{
		commands: commands,
	}
	a.rl.SetAutoCompleter(completer)
}

// CurrentCommands returns the app's commands.
func (a *App) CurrentCommands() *Commands {
	var ret *Commands
	commands, ok := a.commandsMap[a.currentView]
	if ok {
		ret = commands
	} else {
		ret = nil
	}
	return ret
}

// PrintError prints the given error.
func (a *App) PrintError(err error) {
	if a.config.NoColor {
		a.Printf("error: %v\n", err)
	} else {
		a.config.ErrorColor.Fprint(a, "error: ")
		a.Printf("%v\n", err)
	}
}

// Printf formats according to a format specifier and writes to terminal output.
// Printf writes to standard output if terminal output is not yet active.
func (a *App) Printf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(a, format, args...)
}

// Println writes to terminal output followed by a newline.
// Println writes to standard output if terminal output is not yet active.
func (a *App) Println(args ...interface{}) (int, error) {
	return fmt.Fprintln(a, args...)
}

// OnInit sets the function which will be executed before the first command
// is executed. App flags can be handled here.
func (a *App) OnInit(f func(a *App, flags FlagMap) error) {
	a.initHook = f
}

// OnShell sets the function which will be executed before the shell starts.
func (a *App) OnShell(f func(a *App) error) {
	a.shellHook = f
}

// SetInterruptHandler sets the interrupt handler function.
func (a *App) SetInterruptHandler(f func(a *App, count int)) {
	a.interruptHandler = f
}

// SetPrintHelp sets the print help function.
func (a *App) SetPrintHelp(f func(a *App, shell bool)) {
	a.printHelp = f
}

// SetPrintCommandHelp sets the print help function for a single command.
func (a *App) SetPrintCommandHelp(f func(a *App, c *Command, shell bool)) {
	a.printCommandHelp = f
}

// SetPrintASCIILogo sets the function to print the ASCII logo.
func (a *App) SetPrintASCIILogo(f func(a *App)) {
	a.printASCIILogo = func(a *App) {
		if !a.config.NoColor {
			a.config.ASCIILogoColor.Set()
			defer color.Unset()
		}
		f(a)
	}
}

// Write to the underlying output, using readline if available.
func (a *App) Write(p []byte) (int, error) {
	return a.Stdout().Write(p)
}

// Stdout returns a writer to Stdout, using readline if available.
// Note that calling before Run() will return a different instance.
func (a *App) Stdout() io.Writer {
	if a.rl != nil {
		return a.rl.Stdout()
	}
	return os.Stdout
}

// Stderr returns a writer to Stderr, using readline if available.
// Note that calling before Run() will return a different instance.
func (a *App) Stderr() io.Writer {
	if a.rl != nil {
		return a.rl.Stderr()
	}
	return os.Stderr
}

// AddCommandAll adds a new command.
// Panics on error.
func (a *App) AddCommandAll(cmd *Command) {
	a.addCommandAll(cmd, true)
}

// AddCommand adds a new command.
// Panics on error.
func (a *App) AddCommand(view int32, cmd *Command) {
	a.addCommand(view, cmd, true)
}

// addCommand adds a new command.
// If addHelpFlag is true, a help flag is automatically
// added to the command which displays its usage on use.
// Panics on error.
func (a *App) addCommand(view int32, cmd *Command, addHelpFlag bool) {
	err := cmd.validate()
	if err != nil {
		panic(err)
	}
	cmd.registerFlags(addHelpFlag)

	commands, ok := a.commandsMap[view]
	if !ok {
		fmt.Println("unknow view , only aupport authNode userNode testNode")
		return
	}
	commands.Add(cmd)
}

func (a *App) addCommandAll(cmd *Command, addHelpFlag bool) {
	err := cmd.validate()
	if err != nil {
		panic(err)
	}
	cmd.registerFlags(addHelpFlag)

	for _, v := range a.commandsMap {
		v.Add(cmd)
	}
}

// RunCommand runs a single command.
func (a *App) RunCommand(args []string) error {
	// Parse the arguments string and obtain the command path to the root.
	cmds, fg, args, err := a.CurrentCommands().parse(args, a.flagMap, false)
	if err != nil {
		return err
	} else if len(cmds) == 0 {
		return fmt.Errorf("incorrect input, try 'help'")
	}

	// The last command is the final command.
	cmd := cmds[len(cmds)-1]

	// Check if arguments are allowed.
	if !cmd.AllowArgs && len(args) > 0 {
		return fmt.Errorf("command '%s' requires no arguments, try 'help'", cmd.Name)
	}

	// Print the command help if the command run function is nil or if the help flag is set.
	if fg.Bool("help") || cmd.Run == nil {
		a.printCommandHelp(a, cmd, a.isShell)
		return nil
	}

	// Create the context and pass the rest args.
	ctx := newContext(a, cmd, fg, args)

	// Run the command.
	err = cmd.Run(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Run the application and parse the command line arguments.
// This method blocks.
func (a *App) Run() (err error) {
	defer a.Close()

	// Sort all commands by their name.
	a.CurrentCommands().SortRecursive()

	// Remove the program name from the args.
	args := os.Args
	if len(args) > 0 {
		args = args[1:]
	}

	// Parse the app command line flags.
	args, err = a.flags.parse(args, a.flagMap)
	if err != nil {
		return err
	}

	// Check if nocolor was set.
	a.config.NoColor = a.flagMap.Bool("nocolor")

	// Determine if this is a shell session.
	a.isShell = len(args) == 0

	// Add general builtin commands.
	a.addCommandAll(&Command{
		Name:      "help",
		Help:      "use 'help [command]' for command help",
		AllowArgs: true,
		Run: func(c *Context) error {
			if len(c.Args) == 0 {
				a.printHelp(a, a.isShell)
				return nil
			}
			cmd, _, err := a.CurrentCommands().FindCommand(c.Args)
			if err != nil {
				return err
			} else if cmd == nil {
				a.PrintError(fmt.Errorf("command not found"))
				return nil
			}
			a.printCommandHelp(a, cmd, a.isShell)
			return nil
		},
	}, false)

	// Check if help should be displayed.
	if a.flagMap.Bool("help") {
		a.printHelp(a, false)
		return nil
	}

	// Run the init hook.
	if a.initHook != nil {
		err = a.initHook(a, a.flagMap)
		if err != nil {
			return err
		}
	}

	// Check if a command chould be executed in non-interactive mode.
	if !a.isShell {
		return a.RunCommand(args)
	}

	// Add shell builtin commands.
	a.AddCommandAll(&Command{
		Name: "exit",
		Help: "exit the shell",
		Run: func(c *Context) error {
			var view int32
			view = a.CurrentView()
			switch view {
			case UserNODE:
				a.SetCurrentView(AuthNODE)
			case TestNODE:
				a.SetPrompt("venus > ")
				a.SetCurrentView(UserNODE)
				a.SetAutoCompleter(a.CurrentCommands())
			default:
				c.Stop()
			}
			fmt.Print("\033c")
			return nil
		},
	})
	a.AddCommandAll(&Command{
		Name: "clear",
		Help: "clear the screen",
		Run: func(c *Context) error {
			readline.ClearScreen(a.rl)
			return nil
		},
	})

	// Create the readline instance.
	a.rl, err = readline.NewEx(&readline.Config{
		Prompt:                 a.currentPrompt,
		HistorySearchFold:      true, // enable case-insensitive history searching
		DisableAutoSaveHistory: true,
		HistoryFile:            a.config.HistoryFile,
		HistoryLimit:           a.config.HistoryLimit,
		AutoComplete:           newCompleter(a.CurrentCommands()),
	})
	if err != nil {
		return err
	}
	a.OnClose(a.rl.Close)

	// Run the shell hook.
	if a.shellHook != nil {
		err = a.shellHook(a)
		if err != nil {
			return err
		}
	}

	// Run the shell.
	return a.runShell()
}

func (a *App) runShell() error {
	var interruptCount int
	var lines []string
	multiActive := false
Loop:
	for !a.IsClosing() {
		//Auth
		if a.CurrentView() == AuthNODE {
			fmt.Print("\033c")
			var name string
			fmt.Print("Username:")
			fmt.Scanf("%v", &name)
			s3, _ := gp.GetPasswdPrompt("Password:", true, os.Stdin, os.Stdout)
			passWd := string(s3)
			if passWd == "venus70" && name == "adm" {
				a.SetPrompt("venus > ")
				a.SetCurrentView(UserNODE)
				a.SetAutoCompleter(a.CurrentCommands())
				readline.ClearScreen(a.rl)
				if a.printASCIILogo != nil {
					a.printASCIILogo(a)
				}
			} else {
				a.SetCurrentView(AuthNODE)
				a.SetAutoCompleter(a.CurrentCommands())
				fmt.Print("\033c")
				continue Loop
			}

		}
		// Set the prompt.
		if multiActive {
			a.rl.SetPrompt(a.config.multiPrompt())
		} else {
			a.rl.SetPrompt(a.currentPrompt)
		}
		multiActive = false
		// Readline.
		line, err := a.rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				interruptCount++
				a.interruptHandler(a, interruptCount)
				continue Loop
			} else if err == io.EOF {
				return nil
			} else {
				return err
			}
		}

		// Reset the interrupt count.
		interruptCount = 0

		// Handle multiline input.
		if strings.HasSuffix(line, "\\") {
			multiActive = true
			line = strings.TrimSpace(line[:len(line)-1]) // Add without suffix and trim spaces.
			lines = append(lines, line)
			continue Loop
		}
		lines = append(lines, strings.TrimSpace(line))

		line = strings.Join(lines, " ")
		line = strings.TrimSpace(line)
		lines = lines[:0]

		// Skip if the line is empty.
		if len(line) == 0 {
			continue Loop
		}

		// Save command history.
		err = a.rl.SaveHistory(line)
		if err != nil {
			a.PrintError(err)
			continue Loop
		}

		// Split the line to args.
		args, err := shlex.Split(line, true)
		if err != nil {
			a.PrintError(fmt.Errorf("invalid args: %v", err))
			continue Loop
		}

		// Execute the command.
		err = a.RunCommand(args)
		if err != nil {
			a.PrintError(err)
			continue Loop
		}
	}

	return nil
}
