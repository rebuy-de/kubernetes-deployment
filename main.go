package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/rebuy-de/kubernetes-deployment/cmd"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func handleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(cmd.Exit); ok == true {
			os.Exit(exit.Code)
		}
		panic(e) // not an Exit, bubble up
	}
}

func main() {
	defer handleExit()
	if err := cmd.NewRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
