package main

const (
	Exit             = "exit"
	Echo             = "echo"
	Type             = "type"
	Pwd              = "pwd"
	Cd               = "cd"
	HomeDirectory    = "~"
	Redirect         = ">"
	ErrorRedirect    = "2>"
	StandardRedirect = "1>"
)

var commands = []string{Exit, Echo, Type, Pwd, Cd}

func Commands() []string {
	out := make([]string, len(commands))
	copy(out, commands)
	return out
}
