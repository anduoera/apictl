package core

import (
	"context"
	"github.com/pterm/pterm"
	"sync"
)

var (
	ctx  *Context
	once sync.Once
)

type Context struct {
	context context.Context
	logger  Logger
}

func GetInstance() *Context {
	once.Do(func() {
		ctx = &Context{
			context: context.Background(),
			logger: Logger{
				Level: ErrorLevel,
			},
		}
	})
	return ctx
}

func init() {
	pterm.EnableDebugMessages()
}
