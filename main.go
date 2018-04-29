package main

import "github.com/opb/seekrits/cmd"

var Version = "unknown"

func main() {
	cmd.Version = Version
	cmd.Execute()
}
