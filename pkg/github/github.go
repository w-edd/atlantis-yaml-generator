package github

import (
	"context"
	"errors"
	"fmt"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/w-edd/atlantis-yaml-generator/pkg/config"
	"github.com/w-edd/atlantis-yaml-generator/pkg/helpers"
	"golang.org/x/oauth2"
)

type GithubRequest struct {
	AuthToken         string
	Owner             string
	Repo              string
	PullRequestNumber string
}

// NewGitHubClient creates a new GitHub client with the provided auth token.
func newGitHubClient(authToken string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: authToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

// runGHRequest returns a list of changed files in a pull request.
func runGHRequest(authToken, owner, repo, pullReqNum string) ([]string, error) {
	var changedFiles []string
	prNum, err := strconv.Atoi(pullReqNum)
	if err != nil {
		return nil, err
	}
	client := newGitHubClient(authToken)
	files, _, err := client.PullRequests.ListFiles(context.Background(), owner, repo, prNum, nil)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		changedFiles = append(changedFiles, *file.Filename)
	}
	return changedFiles, err
}

// GetChangedFiles gets the parameters to call a ghrequest that returns a list of changed files.
func GetChangedFiles() (ChangedFiles []string, err error) {
	// Parse the token from the git config file
	token, _ := getTokenFromGitCredentialsFile()
	if token == "" {
		token = config.GlobalConfig.Parameters["gh-token"]
	}
	if token == "" {
		err = errors.New("gh-token could not be parsed from .git/config file.\n" +
			"Please use gh-token parameter or GH_TOKEN environment variable to set the token.")
		return ChangedFiles, err
	}
	prChangedFiles, err := runGHRequest(
		token,
		config.GlobalConfig.Parameters["base-repo-owner"],
		config.GlobalConfig.Parameters["base-repo-name"],
		config.GlobalConfig.Parameters["pull-num"])
	if err != nil {
		return []string{}, err
	}
	return prChangedFiles, err
}

func getTokenFromGitCredentialsFile() (string, error) {
	// Get the current user
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	// Construct the path to the .git-credentials file
	credentialsFilePath := filepath.Join(usr.HomeDir, ".git-credentials")

	// ReadFile credentials file
	file, err := helpers.ReadFile(credentialsFilePath)
	if err != nil {
		return "", err
	}

	token, err := extractTokenFromURL(file)
	return token, err
}

func extractTokenFromURL(urlLine string) (string, error) {
	// Split by "x-access-token:" to extract the token.
	parts := strings.Split(urlLine, "x-access-token:")
	if len(parts) == 2 {
		tokenPart := parts[1]

		// Split again by "@" to extract the token.
		tokenParts := strings.Split(tokenPart, "@")
		if len(tokenParts) >= 1 {
			return tokenParts[0], nil
		}
	}

	return "", fmt.Errorf("token not found in url line")
}
