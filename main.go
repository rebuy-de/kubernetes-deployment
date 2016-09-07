package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	version = "unknown"

	defaultKubeConfigPath    = "~/.kube/config"
	defaultProjectConfigPath = "config/services.yaml"
	defaultOutputPath        = "./output"
)

func main() {
	os.Exit(Main(os.Args...))
}

func Main(args ...string) int {
	log.SetOutput(os.Stdout)

	app := &App{}

	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	fs.StringVar(
		&app.KubeConfigPath,
		"kubeconfig", defaultKubeConfigPath,
		"path to the kubernetes configuration")
	fs.StringVar(
		&app.ProjectConfigPath,
		"config", defaultProjectConfigPath,
		"project configuration file")
	fs.StringVar(
		&app.OutputPath,
		"output", defaultOutputPath,
		"output path of configuration file after shuffling and Kubernetes manifests")
	fs.IntVar(
		&app.SleepInterval,
		"sleep", 0,
		"sleep interval between applying projects")
	fs.BoolVar(
		&app.SkipShuffle,
		"skip-shuffle", false,
		"skip shuffling of project order")
	fs.BoolVar(
		&app.SkipFetch,
		"skip-fetch", false,
		"skip fetching files via git; requires valid files in the output directory")
	fs.BoolVar(
		&app.SkipDeploy,
		"skip-deploy", false,
		"skip applying the manifests to kubectl")

	printVersion := fs.Bool(
		"version", false,
		"prints version and exits")

	err := fs.Parse(args)
	if err != nil {
		return 2
	}

	if printVersion != nil && *printVersion {
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
