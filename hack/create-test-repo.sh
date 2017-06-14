#!/bin/bash

set -exu -o pipefail

temp=$( mktemp -d -t kubernetes-deployment.XXXXXX )

(
	cd ${temp}

	git init
	git remote add origin git@github.com:svenwltr/repo-test.git

	echo "bluber" > master.md
	git add master.md
	git commit -m "initial commit"
	git push

	git checkout -b merged-branch
	echo "foobar" > merged-branch.md
	git add merged-branch.md
	git commit -m "merged branch"
	git push
	hub pull-request -m "merged branch"

	git checkout master
	git merge --no-ff merged-branch
	git push

	git checkout -b pending-branch
	echo "bimbaz" > pending-branch.md
	git add pending-branch.md
	git commit -m "pending branch"
	git push
	hub pull-request -m "pending branch"

	git checkout master
	git checkout -b deleted-branch
	echo "bish bash bosh" > deleted-branch.md
	git add deleted-branch.md
	git commit -m "deleted branch"
	git push
	hub pull-request -m "deleted branch"

	git checkout master
	git merge --no-ff deleted-branch
	git push
	git branch -d deleted-branch
	git push origin --delete deleted-branch
)


rm -rf ${temp}
