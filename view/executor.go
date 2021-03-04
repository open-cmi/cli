package view

import (
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/prompt-cli/commands"
)

// MatchElem match elem
func (v *View) MatchElem(iterator *commands.CommandIterator, elem *commands.CommandWordDefElem) (match bool, defs []commands.CommandWordDef) {
	command, err := iterator.Next()
	if err != nil {
		return
	}
	var bElemMatch = true
	// elem match，只要有一个不匹配，则为不匹配
	for idx := 0; idx < len(elem.CommandWordDefs); idx++ {
		worddef := &elem.CommandWordDefs[idx]
		if iterator.Last() {
			// 如果是最后一个，则判断与当前 worddef 是否匹配, 最后一个，非完全匹配
			if !worddef.MatchCommandDef(command, true) {
				bElemMatch = false
			}
			break
		} else {
			if !worddef.MatchCommandDef(command, true) {
				// 如果未匹配，则直接 break
				bElemMatch = false
				break
			}
			// 取下一个字符
			if idx != len(elem.CommandWordDefs)-1 {
				command, err = iterator.Next()
				if err != nil {
					break
				}
			}
		}
	}

	if bElemMatch {
		match = true
		for _, worddef := range elem.CommandWordDefs {
			if worddef.Identity != 0 {
				defs = append(defs, worddef)
			}
		}
	}

	return
}

// MatchWildGroup match
func (v *View) MatchWildGroup(iterator *commands.CommandIterator, group *commands.CommandWordGroup) (match bool, defs []commands.CommandWordDef) {
	// 如果是全部匹配，那么只记录一个，方便继续查找下一个分组。注册命令行应该保证分组唯一性
	// 如果是部分匹配，则记录全部的，此时，不需要担心下一个分组的匹配情况
	matchMap := make(map[int]bool, 1)
	elemIndex := 0
	for elemIndex < len(group.CommandWordDefElems) {
		if matchMap[elemIndex] {
			elemIndex++
			continue
		}
		elem := group.CommandWordDefElems[elemIndex]
		iterator.Backup()
		m, cmddefs := v.MatchElem(iterator, &elem)
		if m {
			match = m
			defs = append(defs, cmddefs...)
			matchMap[elemIndex] = true
			// 有匹配的，那就再匹配一遍
			elemIndex = 0
		} else {
			iterator.Restore()
			elemIndex++
		}
	}
	return
}

// MatchNormalGroup match
func (v *View) MatchNormalGroup(iterator *commands.CommandIterator, group *commands.CommandWordGroup) (match bool, defs []commands.CommandWordDef) {
	// 每个分组，只能匹配一个 elem
	elemIndex := 0
	for {
		elem := group.CommandWordDefElems[elemIndex]
		iterator.Backup()
		m, cmddefs := v.MatchElem(iterator, &elem)
		if m {
			match = m
			defs = append(defs, cmddefs...)
			break
		} else {
			iterator.Restore()
		}
		elemIndex++
		if elemIndex == len(group.CommandWordDefElems) {
			break
		}
	}

	return
}

// MatchGroup match group
func (v *View) MatchGroup(iterator *commands.CommandIterator, group *commands.CommandWordGroup) (match bool, defs []commands.CommandWordDef) {
	if group.IsWild == true {
		match, defs = v.MatchWildGroup(iterator, group)
	} else {
		match, defs = v.MatchNormalGroup(iterator, group)
	}
	return
}

// ExecCommnadLine 执行命令
func (v *View) ExecCommnadLine(command string) (err error) {
	err = errors.New("parse command failed")
	var commanddefs []commands.CommandWordDef
	iterator := commands.NewCommandIterator(command)
	for lstIndex, commandlist := range v.CommandLists {
		iterator.Reset()
		cmdMatch := true
		for _, group := range commandlist.CommandWordGroups {
			// 匹配 group 有三种情况
			// 1. 不匹配，则下一个 command list
			// 2. iterator 中元素大于等于elem 中的元素数，并且完全匹配某个 elem，此时无提示信息，继续下一个分组的匹配
			// 3. iterator 中元素少，返回提示信息，此命令列表终结

			match, defs := v.MatchGroup(iterator, &group)
			if group.Type != "option" && !match {
				cmdMatch = false
				// 任意一个不匹配，则不匹配
				break
			}
			commanddefs = append(commanddefs, defs...)
		}
		if cmdMatch && v.ParserProcs[lstIndex] != nil {
			curView := GetCurrentView()
			err = v.ParserProcs[lstIndex](curView.Ctx, commanddefs)
			break
		}
	}
	return
}

// Executor 输入命令后，执行
func Executor(commandLine string) {
	cmdline := strings.Trim(commandLine, " ")
	if cmdline == "" {
		return
	}
	cur := GetCurrentView()
	if cmdline == "quit" || cmdline == "exit" {
		cur.Exit()
		return
	}

	err := cur.ExecCommnadLine(cmdline)
	if err != nil {
		fmt.Println(err.Error())
	}

	return
}
