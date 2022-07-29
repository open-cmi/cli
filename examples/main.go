package main

import (
	"github.com/open-cmi/cli"
)

func main() {

	// view context
	c := cli.New("cli")
	c.NewDefaultView("sys", ">")
	c.AppendView("sys", "service", ">")
	c.Run()
}
