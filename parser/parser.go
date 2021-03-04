package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/open-cmi/prompt-cli/commands"
)

// TokenGroup token group
type TokenGroup struct {
	Elems  []string
	Type   string // require, option, wild
	IsWild bool   // 是否为 wild
}

// CommandParser command parser
type CommandParser struct {
	Expression   string
	Tokens       []commands.CommandWordDef
	Description  string
	CommandList  commands.CommandList
	commandline  []TokenGroup
	parsinggroup TokenGroup //处理过程中的 group
	parsingtoken string
	registerd    bool
}

func (ce *CommandParser) terminateToken() {
	if ce.parsingtoken != "" {
		ce.parsinggroup.Elems = append(ce.parsinggroup.Elems, ce.parsingtoken)
		ce.parsingtoken = ""
	}

	return
}

func (ce *CommandParser) terminateGroup(nextType string) {
	ce.terminateToken()
	if len(ce.parsinggroup.Elems) != 0 {
		ce.commandline = append(ce.commandline, ce.parsinggroup)
	}

	ce.parsinggroup = TokenGroup{
		Elems: []string{},
		Type:  nextType,
	}

	return
}

func (ce *CommandParser) parsegroup(exp *bytes.Buffer) {

	for cur, err := exp.ReadByte(); err != io.EOF; cur, err = exp.ReadByte() {
		if cur == ' ' {
			continue
		}

		if cur == '{' {
			// 终结本组，进入子组分析
			ce.terminateGroup("require")
			ce.parsegroup(exp)
		} else if cur == '}' {
			ce.terminateGroup("require")
			return
		} else if cur == '[' {
			ce.terminateGroup("option")
			ce.parsegroup(exp)
		} else if cur == ']' {
			ce.terminateGroup("require")
			return
		} else if cur == '|' {
			ce.terminateToken()
		} else if cur == '$' {
			t, err := exp.ReadBytes(' ')

			if err != nil && err != io.EOF {
				fmt.Println("comand express parse failed")
				return
			}
			ce.parsingtoken += (string(t) + " ")
		} else if cur == '*' {
			var wild = &ce.commandline[len(ce.commandline)-1]
			wild.IsWild = true
		}
	}
	ce.terminateGroup("require")
	return
}

// ParseGroup parse how many group
func (ce *CommandParser) ParseGroup() {
	expbuff := bytes.NewBufferString(ce.Expression)
	ce.parsegroup(expbuff)
}

func (ce *CommandParser) translateToCommandWordDefElem(token string) (cmdElem commands.CommandWordDefElem, err error) {
	arr := strings.Split(token, " ")
	for _, t := range arr {
		if strings.Trim(t, " ") != "" {
			idx, err := strconv.Atoi(t)
			if err != nil {
				return cmdElem, errors.New("register command failed")
			}

			cmddef := ce.Tokens[idx]
			cmdElem.CommandWordDefs = append(cmdElem.CommandWordDefs, cmddef)
		}
	}
	return cmdElem, nil
}

// Parse register commands to view
func (ce *CommandParser) Parse() (err error) {
	if ce.registerd {
		return errors.New("command has been register")
	}
	ce.registerd = true
	ce.terminateGroup("require")
	ce.ParseGroup()

	for _, tokengroup := range ce.commandline {
		var commandGroup commands.CommandWordGroup
		for _, token := range tokengroup.Elems {
			cmdElem, err := ce.translateToCommandWordDefElem(token)
			if err != nil {
				return err
			}
			commandGroup.CommandWordDefElems = append(commandGroup.CommandWordDefElems, cmdElem)
			commandGroup.Type = tokengroup.Type
			commandGroup.IsWild = tokengroup.IsWild
		}
		ce.CommandList.CommandWordGroups = append(ce.CommandList.CommandWordGroups, commandGroup)
	}
	return
}
