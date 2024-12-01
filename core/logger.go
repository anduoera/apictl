package core

import "github.com/pterm/pterm"

type LogLevel int

const (
	LoggerLevelDebug   = "debug"
	LoggerLevelError   = "error"
	LoggerLevelInfo    = "info"
	LoggerLevelWarn    = "warn"
	LoggerLevelSuccess = "success"
)

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarningLevel
	ErrorLevel
	SuccessLevel
)

var LoggerLevelMap = map[string]LogLevel{
	LoggerLevelInfo:    InfoLevel,
	LoggerLevelError:   ErrorLevel,
	LoggerLevelDebug:   DebugLevel,
	LoggerLevelWarn:    WarningLevel,
	LoggerLevelSuccess: SuccessLevel,
}

type Logger struct {
	Level LogLevel
}

func (c *Context) Logger(level LogLevel, msg ...string) {
	message := ""
	for _, v := range msg {
		message = message + v
	}
	if level >= c.logger.Level {
		switch level {
		case DebugLevel:
			pterm.Debug.Println(message)
		case InfoLevel:
			pterm.Info.Println(message)
		case WarningLevel:
			pterm.Warning.Println(message)
		case ErrorLevel:
			pterm.Error.Println(message)
		case SuccessLevel:
			pterm.Success.Println(message)
		}
	}
}

func (c *Context) SetLoggerLevel(level string) {
	if _, ok := LoggerLevelMap[level]; !ok {
		return
	}
	c.logger.Level = LoggerLevelMap[level]
}

func (c *Context) SetLoggerLevelStr(args []string) {
	level := ""
	for _, v := range args {
		if _, ok := LoggerLevelMap[v]; !ok {
			continue
		}
		level = v
	}

	c.SetLoggerLevel(level)
}
