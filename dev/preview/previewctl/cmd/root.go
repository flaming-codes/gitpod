/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "previewctl",
		Short: "Your best friend when interacting with Preview Environments :)",
		Long:  `previewctl is your best friend when interacting with Preview Environments :)`,
	}
	cmd.AddCommand(
		installContextCmd(),
	)
	return cmd
}
