package view

import (
	"errors"
	"os"

	"github.com/open-cmi/cli/commands"
	"github.com/open-cmi/cli/parser"
	"github.com/open-cmi/cli/prompt"
)

// gPrompt terminal prompt
var gPrompt *prompt.Prompt

// Context view context
type Context struct {
	Name  string
	Value interface{}
}

// ParseFunc parse func
type ParseFunc func(ctx Context, cmddefs []commands.CommandWordDef) error

type GetViewContextListFunc func() (ctxs []Context, err error)

// View 系统视图
type View struct {
	Name            string
	DefaultViewText string
	CommandLists    []commands.CommandList
	ParserProcs     []ParseFunc
	Ctx             Context
	Parent          *View
	Children        []*View
	GetContextList  GetViewContextListFunc
}

// gViewMap view map
var gViewMap map[string]*View
var gCurViewName string

// SetPrompt 设置prompt
func SetPrompt(p *prompt.Prompt) {
	gPrompt = p
}

// Exit exit
func (v *View) Exit() {
	if v.Parent == nil {
		os.Exit(0)
		return
	}
	SetCurrentView(v.Parent.Name, "", nil)
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
		for _, v := range gViewMap {
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
	return gViewMap[gCurViewName]
}

// SetCurrentView set view
func SetCurrentView(name string, text string, data interface{}) (err error) {

	view := gViewMap[name]
	if view == nil {
		return errors.New("view not exist")
	}
	gCurViewName = name
	if text == "" {
		gPrompt.SetPrefix(view.DefaultViewText)
	} else {
		gPrompt.SetPrefix(text)
	}
	view.Ctx = Context{
		Name:  name,
		Value: data,
	}
	return
}

func SetViewGetContextListFunc(name string, getList GetViewContextListFunc) error {
	view := gViewMap[name]
	if view == nil {
		return errors.New("view not exist")
	}
	view.GetContextList = getList
	return nil
}

// GetView get view
func GetView(name string) (v *View) {
	return gViewMap[name]
}

// Register register view
func Register(v *View, parent *View) (err error) {
	if gViewMap[v.Name] == nil {
		gViewMap[v.Name] = v
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
	gViewMap = make(map[string]*View)
}
