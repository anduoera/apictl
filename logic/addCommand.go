package logic

import "autoApi/core"

type addCommand struct {
	ctx  *core.Context
	args []string
}

func init() {
	AddCommandTag(AddCommandUse, newAddCommand)
}

func newAddCommand(ctx *core.Context, args []string) Command {
	return &addCommand{
		ctx:  ctx,
		args: args,
	}
}

func (a *addCommand) Run() {
	if a.CheckProjectStruct() {
		a.ctx.Logger("project stuct fail")
	}
}

func (a *addCommand) CheckProjectStruct() bool {

	return false
}
