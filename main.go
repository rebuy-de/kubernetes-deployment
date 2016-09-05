package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	version string
)

func main() {
	app := &App{}

	flag.StringVar(
		&app.KubeConfigPath,
		"kubeconfig", "~/.kube/config",
		"path to the kubernetes configuration")
	flag.StringVar(
		&app.ProjectConfigPath,
		"config", "config/services.yaml",
		"project configuration file")
	flag.StringVar(
		&app.OutputPath,
		"output", "./output",
		"output path of configuration file after shuffling and Kubernetes manifests")
	flag.IntVar(
		&app.SleepInterval,
		"sleep", 0,
		"sleep interval between applying projects")
	flag.BoolVar(
		&app.SkipShuffle,
		"skip-shuffle", false,
		"skip shuffling of project order")

	printVersion := flag.Bool(
		"version", false,
		"prints version and exits")

	flag.Parse()

	if printVersion != nil && *printVersion {
		fmt.Printf("kubernetes-deployment version %s\n", version)
		os.Exit(0)
	}

	err := app.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Execution of the deployment failed:")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
