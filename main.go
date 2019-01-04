package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/rebuy-de/kubernetes-deployment/cmd/kubernetes-deployment"
	"github.com/rebuy-de/rebuy-go-sdk/cmdutil"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	defer cmdutil.HandleExit()
	if err := cmd.NewRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
