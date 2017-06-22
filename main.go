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

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
