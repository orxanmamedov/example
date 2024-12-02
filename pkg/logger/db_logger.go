package logger

import "go.uber.org/zap"

type PrintableLogger struct {
	*zap.SugaredLogger
}

func NewPrintableLogger() *PrintableLogger {
	return &PrintableLogger{global}
}

func (g *PrintableLogger) Print(v ...interface{}) {
	g.Info(v...)
}

func (g *PrintableLogger) Println(v ...interface{}) {
	g.Info(v...)
}

func (g *PrintableLogger) Printf(format string, v ...interface{}) {
	g.Infof(format, v...)
}
