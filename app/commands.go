package main

const (
	Exit                 = "exit"
	Echo                 = "echo"
	Type                 = "type"
	Pwd                  = "pwd"
	Cd                   = "cd"
	HomeDirectory        = "~"
	Redirect             = ">"
	ErrorRedirect        = "2>"
	StandardRedirect     = "1>"
	OutputAppend         = ">>"
	StandardOutputAppend = "1>>"
	ErrorAppend          = "2>>"
)

var commands = []string{Exit, Echo, Type, Pwd, Cd}
var redirectOutputCommands = []string{Redirect, ErrorRedirect, StandardRedirect, OutputAppend, StandardOutputAppend, ErrorAppend}

func Commands() []string {
	out := make([]string, len(commands))
	copy(out, commands)
	return out
}

func RedirectOutputCommands() []string {
	out := make([]string, len(redirectOutputCommands))
	copy(out, redirectOutputCommands)
	return out
}
