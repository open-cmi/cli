package main

import (
	"github.com/open-cmi/cli"
)

func main() {

	// view context
	c := cli.New("cli")
	c.AddDefaultView("sys", ">")
	c.AppendView("sys", "service", ">")
	c.Run()
}
