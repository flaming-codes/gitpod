/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/

package preview

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type Preview struct {
	Branch         string
	KubeconfigPath string
}

func New(branch string) *Preview {
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

	return &Preview{
		Branch:         branch,
		KubeconfigPath: "",
	}
}

func (p *Preview) InstallContext(shouldWait bool) error {
	// previewName := p.GetPreviewName()

	fmt.Println("Install Context isn't implemented")
	return nil
}

func (p *Preview) GetPreviewName() string {
	withoutRefsHead := strings.Replace(p.Branch, "/refs/heads/", "", 1)
	lowerCased := strings.ToLower(withoutRefsHead)

	var re = regexp.MustCompile(`[^-a-z0-9]`)
	sanitizedBranch := re.ReplaceAllString(lowerCased, `$1-$2`)

	if len(sanitizedBranch) > 20 {
		h := sha256.New()
		h.Write([]byte(sanitizedBranch))
		hashedBranch := hex.EncodeToString(h.Sum(nil))

		sanitizedBranch = sanitizedBranch[0:10] + hashedBranch[0:10]
	}

	return sanitizedBranch
}
