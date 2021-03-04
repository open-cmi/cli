package commands

import (
	"errors"
	"strconv"
	"strings"
)

// ValidFunc func
type ValidFunc func(valuetype string, input string, strict bool) bool

// DefaultValidFunc default valic func when custom valid func is nil
func DefaultValidFunc(valuetype string, input string, strict bool) bool {
	if input == "" && !strict {
		return true
	}

	if valuetype == "string" || valuetype == "bool" {
		return true
	}

	if valuetype == "int" || valuetype == "integer" {
		_, err := strconv.Atoi(input)
		if err == nil {
			return true
		}
	}
	return false
}

// SourceParserFunc func
type SourceParserFunc func(valuetype string, input string) (value interface{}, err error)

// DefaultSourceParserFunc parse input and assign value
func DefaultSourceParserFunc(valuetype string, input string) (value interface{}, err error) {
	if valuetype == "string" {
		return input, nil
	}

	if valuetype == "int" || valuetype == "integer" {
		value, err := strconv.Atoi(input)
		return value, err
	}
	return "", errors.New("value type error")
}

// CommandWordDef define command keyword
type CommandWordDef struct {
	Identity     int64            // 命令字的ID
	Name         string           // 命令行助记符
	Helper       string           // 命令行帮助提示
	Type         string           // 如果是 Name,就用助记符比较，如果是 Value，则用 Value 比较
	ValueType    string           // 值类型, bool, integer, string
	Value        interface{}      // 携带的值，初始值为默认值
	SourceParser SourceParserFunc // 源解析函数，通过输入
	ValidFunc    ValidFunc        // 值校验函数
}

// MatchCommandDefByValue match by value
func (cw *CommandWordDef) MatchCommandDefByValue(command string, strict bool) bool {
	var bMatch = false
	// 不根据 name 比对，根据值比对，只要值合法，就可以返回该命令字

	ValidFunc := cw.ValidFunc
	if ValidFunc == nil {
		ValidFunc = DefaultValidFunc
	}
	if ValidFunc(cw.ValueType, command, strict) {
		bMatch = true
		if strict {
			SourceParserFunc := cw.SourceParser
			if SourceParserFunc == nil {
				SourceParserFunc = DefaultSourceParserFunc
			}

			value, err := SourceParserFunc(cw.ValueType, command)
			if err == nil {
				cw.Value = value
			}
		}
	}
	return bMatch
}

// MatchCommandDefByName match by name
func (cw *CommandWordDef) MatchCommandDefByName(command string, strict bool) bool {
	var bMatch = false
	if strict {
		if command != "" && strings.HasPrefix(cw.Name, command) {
			bMatch = true
		}
	} else {
		if strings.HasPrefix(cw.Name, command) {
			bMatch = true
		}
	}
	return bMatch
}

// MatchCommandDef 判断是否匹配
func (cw *CommandWordDef) MatchCommandDef(command string, strict bool) bool {
	if cw.Type == "name" {
		return cw.MatchCommandDefByName(command, strict)
	} else if cw.Type == "value" {
		return cw.MatchCommandDefByValue(command, strict)
	}
	return false
}

// CommandWordDefElem 命令元组
type CommandWordDefElem struct {
	CommandWordDefs []CommandWordDef
}

// CommandWordGroup 元组
type CommandWordGroup struct {
	Type                string
	IsWild              bool
	CommandWordDefElems []CommandWordDefElem
}

// CommandList 命令参数
type CommandList struct {
	CommandWordGroups []CommandWordGroup
}

// CommandWord 命令参数
type CommandWord struct {
	CommandWordDef
	RequireCommands []CommandWord //必须的命令字
}
