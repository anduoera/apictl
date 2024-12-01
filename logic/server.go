package logic

import (
	"apictl/core"
	"fmt"
)

const (
	AddCommandUse = "add"
)

type NewCommandType = func(ctx *core.Context, args []string) Command

var commandsMap = make(map[string]NewCommandType)

func GetCommandTag(command string, ctx *core.Context, args []string) Command {
	if _, ok := commandsMap[command]; !ok {
		ctx.Logger(core.ErrorLevel, fmt.Sprintf("no such %s", command))
	}
	fun := commandsMap[command]
	return fun(ctx, args)
}

func AddCommandTag(command string, newFun func(ctx *core.Context, args []string) Command) {
	commandsMap[command] = newFun
}

type Command interface {
	Run()
}
