package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCommand returns a new instance of the root command
func NewRootCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "argocd-terraform-plugin",
		Short: "This is a plugin to replace <placeholders> with terraform outputs",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.AddCommand(NewGenerateCommand())
	command.AddCommand(NewVersionCommand())

	return command
}
