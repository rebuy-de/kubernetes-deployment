package main

import (
	"fmt"
	"os"

	"github.com/rebuy-de/kubernetes-deployment/cmd"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
