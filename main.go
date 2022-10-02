package main

import (
	"os"

	"github.com/KazanExpress/argocd-terraform-plugin/cmd"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
