package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	version = "unknown"
	defaultProjectConfigPath = "config/services.yaml"
)

func main() {
	os.Exit(Main(os.Args[1:]...))
}

func Main(args ...string) int {
	log.SetOutput(os.Stdout)

	app := &App{}
	printVersion := false

	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	fs.StringVar(
		&app.ProjectConfigPath,
		"config", defaultProjectConfigPath,
		"project configuration file")
	fs.BoolVar(
		&printVersion,
		"version", false,
		"prints version and exits")

	err := fs.Parse(args)
	if err != nil {
		return 2
	}

	if printVersion {
		fmt.Printf("kubernetes-deployment version %s\n", version)
		return 0
	}

	err = app.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Execution of the deployment failed:")
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}
