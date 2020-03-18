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

func gitDeploy(repoPath string) error {
	err := gitCommit(repoPath)
	if isError(err) {
		return err
	}
	err = gitPush(repoPath)
	if isError(err) {
		return err
	}
	return nil
}

func cloneTSUI(rootDir, githubUsername string) error {
	return gitClone(rootDir+"/ts-ui", fmt.Sprintf("git@github.com:%s/ts-ui.git", githubUsername))
}

func cloneTSServerless(rootDir, githubUsername string) error {
	return gitClone(rootDir+"/ts-serverless", fmt.Sprintf("git@github.com:%s/ts-serverless.git", githubUsername))
}

func cloneTelarWeb(rootDir, githubUsername string) error {
	return gitClone(rootDir+"/telar-web", fmt.Sprintf("git@github.com:%s/telar-web.git", githubUsername))
}
