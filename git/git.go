package git

import (
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"strings"

	"github.com/rebuy-de/kubernetes-deployment/util"
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

func (g *Git) Exec(args ...string) error {
	cmd := exec.Command(g.GitPath, args...)
	cmd.Dir = g.Directory

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	go util.PipeToLog(" ", stdout)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	go util.PipeToLog("!", stderr)

	log.Printf("$ git %s", strings.Join(args, " "))
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
	log.Printf("Writing '%s' to %s", dir, infoFile)
	return ioutil.WriteFile(infoFile, []byte(dir), 0644)
}
