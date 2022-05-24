/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/

package preview_test

import (
	"log"
	"testing"

	"github.com/gitpod-io/gitpod/previewctl/pkg/preview"
)

func TestInstallContext(t *testing.T) {

	p := preview.New()

	err := p.InstallContext("", false)
	if err != nil {
		log.Fatal("Expected to succeed!")
	}
}

func TestGetPreviewName(t *testing.T) {

}
