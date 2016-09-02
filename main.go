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
	kubeConfigPath := flag.String(
		"kubeconfig", "~/.kube/config",
		"path to the kubernetes configuration")
	configPath := flag.String(
		"config", "config/services.yaml",
		"project configuration file")
	sleepInterval := flag.Int(
		"sleep", 0,
		"sleep interval between applying projects")
	skipShuffle := flag.Bool(
		"skip-shuffle", false,
		"skip shuffling of project order")
	outputConfigPath := flag.String(
		"output-config", "config/services-last.yaml",
		"output path of configuration file after shuffling")
	printVersion := flag.Bool(
		"version", false,
		"prints version and exits")

	flag.Parse()

	if printVersion != nil && *printVersion {
		fmt.Printf("kubernetes-deployment version %s\n", version)
		os.Exit(0)
	}

	_ = configPath
	_ = outputConfigPath
	_ = sleepInterval
	_ = skipShuffle
	_ = kubeConfigPath

}
