package view

import (
	"errors"
	"os"

	"github.com/open-cmi/prompt-cli/commands"
	"github.com/open-cmi/prompt-cli/parser"
	"github.com/open-cmi/prompt-cli/prompt"
)

// Prompt terminal prompt
var Prompt *prompt.Prompt

// Context view context
type Context struct {
	Name  string
	Value interface{}
}

// ParseFunc parse func
type ParseFunc func(ctx Context, cmddefs []commands.CommandWordDef) error

// View 系统视图
type View struct {
	Name            string
	DefaultViewText string
	CommandLists    []commands.CommandList
	ParserProcs     []ParseFunc
	Ctx             Context
	Parent          *View
	Children        []*View
}

// ViewMap view map
var ViewMap map[string]*View
var curViewName string

// SetPrompt 设置prompt
func SetPrompt(p *prompt.Prompt) {
	Prompt = p
	return
}

// Exit exit
func (v *View) Exit() {
	if v.Parent == nil {
		os.Exit(0)
		return
	}
	SetCurrentView(v.Parent.Name, "", nil)
	return
}

// RegisterCommandParser 视图下的注册
func (v *View) RegisterCommandParser(commandParser *parser.CommandParser, f ParseFunc) (err error) {
	// 进入命令解析
	err = commandParser.Parse()
	if err != nil {
		return
	}

	v.CommandLists = append(v.CommandLists, commandParser.CommandList)
	v.ParserProcs = append(v.ParserProcs, f)
	return
}

// RegisterCommandParser register commands to view
func RegisterCommandParser(name string, commandParser *parser.CommandParser, f ParseFunc) (err error) {
	if name == "all" {
		for _, v := range ViewMap {
			p := *commandParser
			err = v.RegisterCommandParser(&p, f)
			if err != nil {
				break
			}
		}
	} else {
		v := GetView(name)
		if v != nil {
			err = v.RegisterCommandParser(commandParser, f)
		}
	}

	return
}

// GetCurrentView 获取当前view视图
func GetCurrentView() (v *View) {
	return ViewMap[curViewName]
}

// SetCurrentView set view
func SetCurrentView(name string, text string, data interface{}) (err error) {
	curViewName = name
	view := ViewMap[name]
	if view == nil {
		return errors.New("view not exist")
	}
	if text == "" {
		Prompt.SetPrefix(view.DefaultViewText)
	} else {
		Prompt.SetPrefix(text)
	}
	view.Ctx = Context{
		Name:  name,
		Value: data,
	}
	return
}

// GetView get view
func GetView(name string) (v *View) {
	return ViewMap[name]
}

// Register register view
func Register(v *View, parent *View) (err error) {
	if ViewMap[v.Name] == nil {
		ViewMap[v.Name] = v
		v.Parent = parent
		if parent != nil {
			parent.Children = append(parent.Children, v)
		}
	} else {
		err = errors.New("view exist")
	}
	return
}

// NewView new view
func NewView(name string, prompt string) (v *View) {
	v = &View{
		Name:            name,
		DefaultViewText: prompt,
		CommandLists:    []commands.CommandList{},
		ParserProcs:     []ParseFunc{},
	}

	return
}

func init() {
	ViewMap = make(map[string]*View, 8)

}
