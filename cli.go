package cli

import (
	"errors"

	"github.com/open-cmi/cli/prompt"
	"github.com/open-cmi/cli/view"
)

type CLI struct {
	Prompt              *prompt.Prompt
	DefaultView         string
	DefaultPromptPrefix string
}

func New(title string) *CLI {
	p := prompt.New(
		view.Executor,
		view.Completer,
		prompt.OptionPrefix(">"),
		prompt.OptionTitle("xsnos-cli"),
	)
	return &CLI{
		Prompt: p,
	}
}

func (c *CLI) AddDefaultView(name string, promptPrefix string) error {
	// UserView 用户视图
	dv := view.GetView(name)
	if dv != nil {
		return errors.New("view exist")
	}
	if c.DefaultView != "" {
		return errors.New("default view has been added")
	}

	user := view.NewView(name, promptPrefix)
	c.DefaultView = name
	c.DefaultPromptPrefix = promptPrefix
	return view.Register(user, nil)
}

func (c *CLI) AppendView(parent string, name string, promptPrefix string) error {
	dv := view.GetView(parent)
	if dv == nil {
		return errors.New("parent view does not exist")
	}

	// add to parent view
	sys := view.NewView(name, promptPrefix)
	return view.Register(sys, dv)
}

func (c *CLI) Run() {
	view.SetPrompt(c.Prompt)
	view.SetCurrentView(c.DefaultView, c.DefaultPromptPrefix, nil)
	c.Prompt.Run()
}
