/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/

package cmd

import (
	"log"

	"github.com/gitpod-io/gitpod/previewctl/pkg/preview"
	"github.com/spf13/cobra"
)

var (
	shouldWait = false
	branch     = ""
)

func installContextCmd() *cobra.Command {
	p := preview.New()

	cmd := &cobra.Command{
		Use:   "install-context",
		Short: "Installs the kubectl context of a preview environment.",
		PreRunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := p.InstallContext(branch, shouldWait)
			if err != nil {
				log.Fatalf("Couldn't install context for the '%s' preview", p.Name)
			}
		},
	}

	cmd.Flags().StringVar(&branch, "branch", "", "From which branch previewctl should install the context. By default it will use the result of \"git rev-parse --abbrev-ref HEAD\"")
	cmd.Flags().BoolVar(&shouldWait, "wait", false, "If wait is enabled, previewctl will keep trying to install the kube-context every 30 seconds.")
	return cmd
}
