package cmdutil

import "os"

type ExitCode struct{ Code int }

func Exit(code int) {
	panic(ExitCode{Code: code})
}

func HandleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(ExitCode); ok == true {
			os.Exit(exit.Code)
		}
		panic(e) // not an Exit, bubble up
	}
}
