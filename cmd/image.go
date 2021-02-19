package cmd

import (
	"os"
	"strings"

	"github.com/openfaas/faas-cli/schema"
	stack "github.com/openfaas/faas-cli/stack"
	"github.com/openfaas/openfaas-cloud/sdk"
)

func formatImageShaTag(registry string, function *stack.Function, sha string, owner string, repo string) string {
	imageName := function.Image

	repoIndex := strings.LastIndex(imageName, "/")
	if repoIndex > -1 {
		imageName = imageName[repoIndex+1:]
	}

	sha = sdk.FormatShortSHA(sha)

	imageName = schema.BuildImageName(schema.BranchAndSHAFormat, imageName, sha, buildBranch())

	var imageRef string
	sharedRepo := strings.HasSuffix(registry, "/")
	if sharedRepo {
		imageRef = registry[:len(registry)-1] + "/" + owner + "-" + repo + "-" + imageName
	} else {
		imageRef = registry + "/" + owner + "/" + repo + "-" + imageName
	}

	return imageRef
}

func buildBranch() string {
	branch := os.Getenv("build_branch")
	if branch == "" {
		return "master"
	}
	return branch
}
