package gevent

import "fmt"

// Log is a interface
type Log interface {
	Error(msg string)
}

type consoleLogger struct {
}

func newConsoleLogger() Log {
	return &consoleLogger{}
}

func (c consoleLogger) Error(msg string) {
	fmt.Println(msg)
}
