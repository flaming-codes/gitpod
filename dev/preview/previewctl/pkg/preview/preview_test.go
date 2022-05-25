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

	p := preview.New("")

	err := p.InstallContext(false)
	if err != nil {
		log.Fatal("Expected to succeed!")
	}
}

func TestGetPreviewName(t *testing.T) {
	testCases := []struct {
		testName       string
		branch         string
		expectedResult string
	}{
		{
			testName:       "Short branch without special characters",
			branch:         "/refs/heads/testing",
			expectedResult: "testing",
		},
		{
			testName:       "Upper to lower case",
			branch:         "/refs/heads/SCREAMMING",
			expectedResult: "screamming",
		},
		{
			testName:       "Special characters",
			branch:         "/refs/heads/as/test&123.4",
			expectedResult: "as-test-123-4",
		},
		{
			testName:       "Hashed long branch",
			branch:         "/refs/heads/this-is-a-long-branch-that-should-be-replaced-with-a-hash",
			expectedResult: "this-is-a-a868caa3c3",
		},
	}

	for _, tc := range testCases {
		preview := &preview.Preview{
			Branch: tc.branch,
		}

		previewName := preview.GetPreviewName()

		if tc.expectedResult != previewName {
			log.Fatalf("Test '%s' failed. Expected '%s' but got '%s'", tc.testName, tc.expectedResult, previewName)
		}
	}
}
