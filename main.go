package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rebuy-de/kubernetes-deployment/kubernetes"
)

var (
	version = "unknown"

	defaultKubeConfigPath    = "config/kubeconfig.yaml"
	defaultProjectConfigPath = "config/services.yaml"
	defaultOutputPath        = "./output"
)

func main() {
	os.Exit(Main(os.Args[1:]...))
}

func Main(args ...string) int {
	log.SetOutput(os.Stdout)

	app := &App{}
	printVersion := false

	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	kubeConfigPath := fs.String(
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
	fs.DurationVar(
		&app.SleepInterval,
		"sleep", time.Second,
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
	fs.DurationVar(
		&app.RetrySleep,
		"retry-sleep", 250*time.Millisecond,
		"sleep interval between applying projects")
	fs.IntVar(
		&app.RetryCount,
		"retry-count", 3,
		"sleep interval between applying projects")
	fs.BoolVar(
		&app.IgnoreDeployFailures,
		"ignore-deploy-failures", false,
		"continue deploying services, if any service fails")

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

	if _, err := os.Stat(*kubeConfigPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "The kubeconfig '%s' does not exist.\n", *kubeConfigPath)
		return 1
	}

	app.Kubectl, err = kubernetes.New(*kubeConfigPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	err = app.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Execution of the deployment failed:")
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}
