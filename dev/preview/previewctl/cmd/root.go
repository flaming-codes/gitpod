/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	branch = ""
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "previewctl",
		Short: "Your best friend when interacting with Preview Environments :)",
		Long:  `previewctl is your best friend when interacting with Preview Environments :)`,
	}

	cmd.PersistentFlags().StringVar(&branch, "branch", "", "From which branch's preview previewctl should interact with. By default it will use the result of \"git rev-parse --abbrev-ref HEAD\"")

	cmd.AddCommand(
		installContextCmd(),
	)
	return cmd
}
