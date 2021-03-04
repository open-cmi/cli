package main

import (
	prompt "github.com/open-cmi/prompt-cli/prompt"
	view "github.com/open-cmi/prompt-cli/view"
)

var defaultPrefix = ">"
var useSock bool = false

func main() {

	// view context
	p := prompt.New(
		view.Executor,
		view.Completer,
		prompt.OptionPrefix(defaultPrefix),
		prompt.OptionTitle("prompt-cli"),
	)
	p.Run()
}
