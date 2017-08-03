package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/rebuy-de/kubernetes-deployment/cmd"
	"github.com/rebuy-de/kubernetes-deployment/pkg/cmdutil"
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
