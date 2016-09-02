package git

func SparseCheckout(target string, repo string, branch string, dir string) error {
	var err error

	git, err := New(target)
	if err != nil {
		return err
	}

	err = git.Init()
	if err != nil {
		return err
	}

	err = git.RemoteAdd(repo)
	if err != nil {
		return err
	}

	err = git.Config("core.sparseCheckout", "true")
	if err != nil {
		return err
	}

	err = git.SetCheckoutPath(dir)
	if err != nil {
		return err
	}

	err = git.PullShallow(branch)
	if err != nil {
		return err
	}

	return nil
}
