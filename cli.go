package cli

import (
	"errors"

	"github.com/open-cmi/cli/prompt"
	"github.com/open-cmi/cli/view"
)

type CLI struct {
	Prompt              *prompt.Prompt
	DefaultViewName     string
	DefaultPromptPrefix string
}

func New(title string) *CLI {
	p := prompt.New(
		view.Executor,
		view.Completer,
		prompt.OptionPrefix(">"),
		prompt.OptionTitle(title),
	)
	return &CLI{
		Prompt: p,
	}
}

// NewDefaultView new default view
func (c *CLI) NewDefaultView(name string, promptPrefix string) error {
	// SystemView 系统视图
	dv := view.GetView(name)
	if dv != nil {
		return errors.New("view exist")
	}

	if c.DefaultViewName != "" {
		return errors.New("default view has been set")
	}

	sys := view.NewView(name, promptPrefix)
	c.DefaultViewName = name
	c.DefaultPromptPrefix = promptPrefix
	return view.Register(sys, nil)
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
	view.SetCurrentView(c.DefaultViewName, c.DefaultPromptPrefix, nil)
	c.Prompt.Run()
}
