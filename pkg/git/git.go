package git

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/rebuy-de/kubernetes-deployment/pkg/util"
)

type Git struct {
	Directory  string
	RemoteName string
	GitPath    string
}

func New(directory string) (*Git, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}

	git := &Git{
		Directory:  directory,
		RemoteName: "origin",
		GitPath:    gitPath,
	}

	err = git.Exec("version")
	if err != nil {
		return nil, err
	}

	return git, nil

}

func (g *Git) CommitID() (string, error) {
	cmd := exec.Command(g.GitPath, "rev-parse", "--short", "HEAD")
	cmd.Dir = g.Directory

	log.Infof("$ git %s", strings.Join(cmd.Args, " "))

	raw, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	id := strings.ToLower(strings.TrimSpace(string(raw)))

	matched, err := regexp.MatchString("[0-9a-f]{7}", id)
	if err != nil {
		return "", err
	}

	if !matched {
		return "", fmt.Errorf("Invalid return value from git: %s", id)
	}

	return id, nil
}

func (g *Git) Exec(args ...string) error {
	cmd := exec.Command(g.GitPath, args...)
	cmd.Dir = g.Directory

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	go util.PipeToLogrus(log.WithFields(log.Fields{
		"args":      cmd.Args,
		"directory": g.Directory,
		"stream":    "stdout",
	}), stdout)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	go util.PipeToLogrus(log.WithFields(log.Fields{
		"args":      cmd.Args,
		"directory": g.Directory,
		"stream":    "stderr",
	}), stderr)

	log.Infof("$ git %s", strings.Join(args, " "))
	return cmd.Run()
}

func (g *Git) Init() error {
	return g.Exec("init")
}

func (g *Git) RemoteAdd(repo string) error {
	return g.Exec("remote", "add", "-f", "origin", repo)
}

func (g *Git) Config(key, value string) error {
	return g.Exec("config", key, value)
}

func (g *Git) PullShallow(branch string) error {
	return g.Exec("pull", "--depth=1", "origin", branch)
}

func (g *Git) SetCheckoutPath(dir string) error {
	infoFile := path.Join(g.Directory, ".git", "info", "sparse-checkout")
	log.Infof("Writing '%s' to %s", dir, infoFile)
	return ioutil.WriteFile(infoFile, []byte(dir), 0644)
}
