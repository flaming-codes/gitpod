/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/

package preview

import (
	"fmt"
	"log"
	"os/exec"
)

type Preview struct {
	Name           string
	KubeconfigPath string
}

func New() *Preview {
	return &Preview{
		Name:           "",
		KubeconfigPath: "",
	}
}

func (p *Preview) InstallContext(branch string, shouldWait bool) error {
	p.Name = p.GetPreviewName(branch)

	// TODO:

	fmt.Println("Install Context isn't implemented")
	return nil
}

func (p *Preview) GetPreviewName(branch string) string {
	if branch == "" {
		out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
		if err != nil {
			log.Fatalf("Could not retrieve branch name, err: %v", err)
		}
		branch = string(out)
	} else {
		_, err := exec.Command("git", "rev-parse", "--verify", branch).Output()
		if err != nil {
			log.Fatalf("Branch '%s' does not exist", branch)
		}
	}

	//TODO: Sanitize branch name and return

	return branch
}
