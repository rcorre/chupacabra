package main

import (
	"flag"
	"os"
	"path/filepath"
)

type Config struct {
	KubeConfigPath string
}

func NewConfig() *Config {
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE") // windows
	}

	var kubeconfig *string

	if home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	return &Config{
		KubeConfigPath: *kubeconfig,
	}
}
