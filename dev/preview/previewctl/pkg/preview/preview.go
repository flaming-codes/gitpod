/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/

package preview

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"
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
	if shouldWait {
		installTicker := time.NewTicker(30 * time.Second)

		for {
			select {
			case <-installTicker.C:
				if err := installContext(p.Branch); err == nil {
					// No error means successful context installation
					return nil
				}
			}
		}
	}

	return installContext(p.Branch)
}

func installContext(branch string) error {
	return exec.Command("bash", "/workspace/gitpod/dev/preview/install-k3s-kubeconfig.sh", "-b", branch).Run()
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
