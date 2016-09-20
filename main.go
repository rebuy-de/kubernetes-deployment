package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rebuy-de/kubernetes-deployment/kubernetes"
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
	app.ProjectConfigPath = defaultProjectConfigPath

	fs.StringVar(
		&app.LocalConfigPath,
		"config", defaultProjectConfigPath,
		"project configuration file")
	fs.BoolVar(
		&printVersion,
		"version", false,
		"prints version and exits")

	fs.BoolVar(
		&app.IgnoreDeployFailures,
		"ignore-deploy-failures", false,
		"continue deploying services, if any service fails")
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

	err := fs.Parse(args)
	if err != nil {
		return 2
	}

	if printVersion {
		fmt.Printf("kubernetes-deployment version %s\n", version)
		return 0
	}

	app.KubectlBuilder = func(kubeconfig *string) (kubernetes.API, error) {
		return kubernetes.New(*kubeconfig)
	}

	err = app.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Execution of the deployment failed:")
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}
