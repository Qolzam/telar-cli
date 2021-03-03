package cmd

import (
	"fmt"
	"os/exec"

	"github.com/Qolzam/telar-cli/pkg/log"
)

func gitClone(projectDirectory string, repoURL string) error {

	cmd := exec.Command("git", "clone", repoURL, projectDirectory)
	err := cmd.Run()

	if err != nil {
		return err
	}
	return nil
}

func gitStatus(path string) error {
	cmd := exec.Command("git", "status")
	cmd.Dir = path
	err := cmd.Run()

	if err != nil {
		return err
	}
	return nil
}

func gitAdd(path string) error {
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = path
	err := cmd.Run()

	if err != nil {
		return err
	}
	return nil
}

func gitShortSHA(path string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Dir = path
	output, err := cmd.Output()

	if err != nil {
		return "", err
	}
	return string(output), nil
}

func gitCommit(path string) error {

	err := gitAdd(path)
	if err != nil {
		return err
	}

	err = gitStatus(path)
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "commit", "-am", "chore: initialize telar social.")
	cmd.Dir = path
	err = cmd.Run()

	if err != nil {
		return err
	}
	return nil
}

func gitPush(path string) error {
	cmd := exec.Command("git", "push")
	cmd.Dir = path
	err := cmd.Run()

	if err != nil {
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
	path := rootDir + "/ts-ui"
	err, exist := directoryFileExist(path)
	if isError(err) {
		return err
	}
	if exist == true {
		log.Info("Telar social user interface repository exist. %s", path)
		return nil
	}
	return gitClone(rootDir+"/ts-ui", fmt.Sprintf("git@github.com:%s/ts-ui.git", githubUsername))
}

func cloneTSServerless(rootDir, githubUsername string) error {
	path := rootDir + "/ts-serverless"
	err, exist := directoryFileExist(path)
	if isError(err) {
		return err
	}
	if exist == true {
		log.Info("Telar social serverless repository exist. %s", path)
		return nil
	}
	return gitClone(path, fmt.Sprintf("git@github.com:%s/ts-serverless.git", githubUsername))
}

func cloneTelarWeb(rootDir, githubUsername string) error {
	path := rootDir + "/telar-web"
	err, exist := directoryFileExist(path)
	if isError(err) {
		return err
	}
	if exist == true {
		log.Info("Telar web repository exist. %s ", path)
		return nil
	}
	return gitClone(path, fmt.Sprintf("git@github.com:%s/telar-web.git", githubUsername))
}
