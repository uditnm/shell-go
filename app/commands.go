package main

const (
	Exit          = "exit"
	Echo          = "echo"
	Type          = "type"
	Pwd           = "pwd"
	Cd            = "cd"
	HomeDirectory = "~"
)

var commands = []string{Exit, Echo, Type, Pwd, Cd}

func Commands() []string {
	out := make([]string, len(commands))
	copy(out, commands)
	return out
}
