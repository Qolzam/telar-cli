package cmd

import (
	"fmt"
	"os"
	"os/user"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

func gitClone(projectDirectory string, repoURL string) error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	sshAuth, err := ssh.NewPublicKeysFromFile("git", currentUser.HomeDir+"/.ssh/id_rsa", "")
	if err != nil {
		return err
	}
	_, err = git.PlainClone(projectDirectory, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
		Auth:     sshAuth,
	})

	if err != nil {
		return err
	}
	return nil
}

func gitCommit(path string) error {
	repo, err := git.PlainOpen(path)
	if isError(err) {
		return err
	}
	w, err := repo.Worktree()
	if isError(err) {
		return err
	}

	_, err = w.Add(".")
	if isError(err) {
		return err
	}

	fmt.Println(path)
	status, err := w.Status()

	fmt.Println(status)

	commit, err := w.Commit("[TELAR] Initialize Telar Social.", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Telar",
			Email: "support@red-gold.tech",
			When:  time.Now(),
		},
	})
	if isError(err) {
		return err
	}

	obj, err := repo.CommitObject(commit)
	if isError(err) {
		return err
	}

	fmt.Println(obj)
	return nil
}

func gitPush(path string) error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	repo, err := git.PlainOpen(path)
	if isError(err) {
		return err
	}

	sshAuth, err := ssh.NewPublicKeysFromFile("git", currentUser.HomeDir+"/.ssh/id_rsa", "")
	if err != nil {
		return err
	}
	err = repo.Push(&git.PushOptions{
		Auth:     sshAuth,
		Progress: os.Stdout,
	})
	if isError(err) {
		return err
	}
	return nil
}

func cloneTSUI(rootDir string) error {
	return gitClone(rootDir+"/ts-ui", "git@github.com:red-gold/ts-ui.git")
}

func cloneTSServerless(rootDir string) error {
	return gitClone(rootDir+"/ts-serverless", "git@github.com:red-gold/ts-serverless.git")
}

func cloneTelarWeb(rootDir string) error {
	return gitClone(rootDir+"/telar-web", "git@github.com:red-gold/telar-web.git")
}
