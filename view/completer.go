package view

import (
	"sort"

	"github.com/open-cmi/prompt-cli/commands"
	"github.com/open-cmi/prompt-cli/prompt"
)

// EndPrompt end prompt
var EndPrompt prompt.Suggest = prompt.Suggest{
	Text:            "<cr>",
	Description:     "",
	DisableComplete: true,
}

// MergeSuggests merge suggests
func MergeSuggests(destSuggests *[]prompt.Suggest, srcSuggests []prompt.Suggest) {
	for _, srcsug := range srcSuggests {
		bFind := false
		for _, destsug := range *destSuggests {
			if destsug.Text == srcsug.Text && destsug.Description == srcsug.Description {
				bFind = true
			}
		}
		if !bFind {
			*destSuggests = append(*destSuggests, srcsug)
		}
	}

	return
}

func (v *View) getSuggestsFromElem(iterator *commands.CommandIterator,
	elem *commands.CommandWordDefElem) (match bool, suggests []prompt.Suggest, matchString string) {
	command, err := iterator.Next()
	if err != nil {
		return
	}

	var bElemMatch = true
	// elem match，只要有一个不匹配，则为不匹配
	for idx, worddef := range elem.CommandWordDefs {
		if iterator.Last() {
			// 如果是最后一个，则判断与当前 worddef 是否匹配, 最后一个，非完全匹配
			if worddef.MatchCommandDef(command, false) {

				disableComplete := false
				if worddef.Type == "value" {
					disableComplete = true
				}
				suggests = append(suggests, prompt.Suggest{
					Text:            worddef.Name,
					Description:     worddef.Helper,
					DisableComplete: disableComplete,
				})
				if matchString != "" {
					matchString += " "
				}
				matchString += command
			} else {
				bElemMatch = false
			}
			break
		} else {

			if !worddef.MatchCommandDef(command, true) {
				// 如果未匹配，则直接 break
				bElemMatch = false
				break
			}

			if matchString != "" {
				matchString += " "
			}
			matchString += command

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
	}

	return
}

func (v *View) getSuggestsFromWildGroup(iterator *commands.CommandIterator, group *commands.CommandWordGroup) (match bool, suggests []prompt.Suggest) {
	// 如果是全部匹配，那么只记录一个，方便继续查找下一个分组。注册命令行应该保证分组唯一性
	// 如果是部分匹配，则记录全部的，此时，不需要担心下一个分组的匹配情况
	shouldAddEndPrompt := false
	matchMap := make(map[int]bool, 1)
	elemIndex := 0
	for elemIndex < len(group.CommandWordDefElems) {
		if matchMap[elemIndex] {
			elemIndex++
			continue
		}
		elem := group.CommandWordDefElems[elemIndex]
		iterator.Backup()
		m, prompts, matchString := v.getSuggestsFromElem(iterator, &elem)
		if m {
			match = m
			matchMap[elemIndex] = true
			if len(prompts) != 0 {
				// 这里说明是命令行不完整，恢复后，进行下一轮匹配
				MergeSuggests(&suggests, prompts)
				iterator.Restore()
				if matchString != "" {
					shouldAddEndPrompt = false
				}

			} else { //如果至少匹配了一个，则可以输出<cr>
				if matchString != "" {
					shouldAddEndPrompt = true
				}
			}
			// 有匹配的，那就再匹配一遍
			elemIndex = 0
		} else {
			iterator.Restore()
			elemIndex++
		}
	}
	if shouldAddEndPrompt {
		MergeSuggests(&suggests, []prompt.Suggest{EndPrompt})
	}
	return
}

func (v *View) getSuggestsFromNormalGroup(iterator *commands.CommandIterator, group *commands.CommandWordGroup) (match bool, suggests []prompt.Suggest) {
	// 如果是全部匹配，那么只记录一个，方便继续查找下一个分组。注册命令行应该保证分组唯一性
	// 如果是部分匹配，则记录全部的，此时，不需要担心下一个分组的匹配情况

	for _, elem := range group.CommandWordDefElems {
		iterator.Backup()
		m, prompts, _ := v.getSuggestsFromElem(iterator, &elem)
		if m {
			match = m
			if len(prompts) != 0 {
				MergeSuggests(&suggests, prompts)
				iterator.Restore()
			} else {
				// 如果不带提示的，则说明该元素完全匹配，那么不再继续匹配下一个
				break
			}
		} else {
			iterator.Restore()
		}
	}
	return
}

func (v *View) getSuggestsFromGroup(iterator *commands.CommandIterator, group *commands.CommandWordGroup) (match bool, suggests []prompt.Suggest) {
	if group.IsWild != true {
		match, suggests = v.getSuggestsFromNormalGroup(iterator, group)
	} else {
		match, suggests = v.getSuggestsFromWildGroup(iterator, group)
	}
	return
}

// GetSuggests 获取当前视图的 suggest
func (v *View) GetSuggests(input string) (suggests []prompt.Suggest) {
	iterator := commands.NewCommandIterator(input)
	for _, commandlist := range v.CommandLists {
		iterator.Reset()
		addEndPrompt := false
		for groupIndex, group := range commandlist.CommandWordGroups {
			// 匹配 group 有三种情况
			// 1. 不匹配，则下一个 command list
			// 2. iterator 中元素大于等于elem 中的元素数，并且完全匹配某个 elem，此时无提示信息，继续下一个分组的匹配
			// 3. iterator 中元素少，返回提示信息，此命令列表终结

			match, prompts := v.getSuggestsFromGroup(iterator, &group)
			// 未匹配时，require 进入下一个命令行匹配，option 进入下一组匹配
			if !match {
				if group.Type == "option" {
					continue
				}
				break
			}

			if len(prompts) != 0 {
				MergeSuggests(&suggests, prompts)
				break
			}

			if groupIndex == len(commandlist.CommandWordGroups)-1 {
				addEndPrompt = true
			}
			
			// 当当前组为最后一个组或者最后一个require 时，且命令行已经跑完的情况下, 提示<cr>
			lstRequire := true
			for gi := groupIndex + 1; gi < len(commandlist.CommandWordGroups); gi++ {
				if commandlist.CommandWordGroups[gi].Type == "require" {
					lstRequire = false
				}
			}
			
			if lstRequire && !iterator.HasVisualCharactor() {
				addEndPrompt = true
			}
			
			if addEndPrompt {
				MergeSuggests(&suggests, []prompt.Suggest{EndPrompt})
			}
		}
	}

	sort.SliceStable(suggests, func(i, j int) bool {
		if suggests[i].Text == "<cr>" {
			return false
		}
		if suggests[j].Text == "<cr>" {
			return true
		}
		return suggests[i].Text < suggests[j].Text
	})
	return
}

// Completer 提示信息
func Completer(in prompt.Document) []prompt.Suggest {
	command := in.TextBeforeCursor()

	suggests := []prompt.Suggest{}

	if in.LastKeyStroke() != prompt.Tab && in.LastKeyStroke() != prompt.QuestionMark {
		return suggests
	}

	cur := GetCurrentView()
	return cur.GetSuggests(command)
}
