package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/avast/retry-go/v4"
	"github.com/google/go-github/v32/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func getLatestRelease(ctx context.Context, client *github.Client, repoUrl string) string {
	parts := strings.Split(repoUrl, "/")

	// The owner and repository name are the last two parts
	owner := parts[len(parts)-2]
	repo := parts[len(parts)-1]

	// Use path.Clean() to remove any trailing slash
	owner = path.Clean(owner)
	repo = path.Clean(repo)

	var latestRelease *github.RepositoryRelease
	var response *github.Response
	var err error

	err = retry.Do(
		func() error {
			latestRelease, response, err = client.Repositories.GetLatestRelease(ctx, owner, repo)
			if err != nil {
				return err
			}

			return nil
		},

		retry.Attempts(5),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("%s #%d: %s\n", repoUrl, n, err)
		}),
	)

	if err != nil {
		fmt.Println(err)
	}

	if response.StatusCode == 401 {
		log.Fatal("401 Bad credentials")
	}

	if latestRelease.GetName() == "" {
		return fmt.Sprintf(`[%s repo]
There's no latest release.`, repo)
	}

	return fmt.Sprintf(`[%s repo]
Current release: %s 
Url: %s 
Last published: %s`, repo, latestRelease.GetName(),
		latestRelease.GetHTMLURL(),
		NewTimeFormat(latestRelease.GetPublishedAt()))
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	var githubToken string
	var provider string
	// Load the GitHub personal access token and provider from an environment variable
	githubToken = os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
	provider = os.Getenv("PROVIDER")

	if provider == "" {
		log.Fatal("Provider is missing.")
	}

	if githubToken == "" {
		log.Fatal("Github Token is missing.")
	}

	provider, err := checkProvider(provider)

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)

	tc := oauth2.NewClient(ctx, ts)

	// Create a new GitHub client
	client := github.NewClient(tc)

	file, err := os.Open("list-repo-to-watch.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		repoUrl := scanner.Text()
		message := getLatestRelease(ctx, client, repoUrl)
		sendNotification(provider, message)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
