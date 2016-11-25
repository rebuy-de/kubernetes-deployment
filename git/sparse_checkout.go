package git

func SparseCheckout(target string, repo string, branch string, dir string) (*Git, error) {
	var err error

	git, err := New(target)
	if err != nil {
		return nil, err
	}

	err = git.Init()
	if err != nil {
		return nil, err
	}

	err = git.RemoteAdd(repo)
	if err != nil {
		return nil, err
	}

	err = git.Config("core.sparseCheckout", "true")
	if err != nil {
		return nil, err
	}

	err = git.SetCheckoutPath(dir)
	if err != nil {
		return nil, err
	}

	err = git.PullShallow(branch)
	if err != nil {
		return nil, err
	}

	return git, nil
}
