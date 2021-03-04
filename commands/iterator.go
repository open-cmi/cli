package commands

import (
	"io"
	"strings"

	"github.com/open-cmi/prompt-cli/container/stack"
)

// CommandIterator command iterator
type CommandIterator struct {
	Expression      string
	CommandKeywords []string
	backups         *stack.Stack
	len             int
	curpos          int
}

// NewCommandIterator 创建辛的迭代器
func NewCommandIterator(expression string) (iter *CommandIterator) {
	keywords := []string{}

	commands := strings.Split(expression, " ")
	for idx, cmd := range commands {
		cmd = strings.Trim(cmd, " ")
		if idx != len(commands)-1 {
			if cmd == "" {
				continue
			}
		}
		keywords = append(keywords, cmd)
	}

	curpos := -1
	backups := stack.NewStack()

	return &CommandIterator{
		Expression:      expression,
		CommandKeywords: keywords,
		len:             len(keywords),
		curpos:          curpos,
		backups:         backups,
	}
}

// Backup 备份
func (ci *CommandIterator) Backup() {
	ci.backups.Push(ci.curpos)
}

// Restore 恢复
func (ci *CommandIterator) Restore() {
	ci.curpos = ci.backups.Pop().(int)
}

// First whether it is first element
func (ci *CommandIterator) First() bool {
	if ci.curpos == 0 {
		return true
	}
	return false
}

// Last whether it is last element
func (ci *CommandIterator) Last() bool {
	if ci.curpos == ci.len-1 {
		return true
	}
	return false
}

// Next 命令获取下一个
func (ci *CommandIterator) Next() (command string, err error) {

	if ci.curpos == ci.len-1 {
		return "", io.EOF
	}
	ci.curpos++
	command = ci.CommandKeywords[ci.curpos]
	return command, nil
}

// Reset reset from 0
func (ci *CommandIterator) Reset() {
	ci.curpos = -1
	for !ci.backups.Empty() {
		ci.backups.Pop()
	}
}
