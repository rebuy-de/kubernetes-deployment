package main

import (
	"os"

	"github.com/rebuy-de/kubernetes-deployment/cmd"
)

func main() {
	os.Exit(cmd.Main(os.Args[1:]...))
}
