package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	owner         = "j4ck4l-24"    // repo owner
	repo          = "testing_tool" // repo name
	basePath      = "./"           //path to save files locally
	dirToDownload = "test1"        // folder to download from repo
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("GITHUB_TOKEN environment variable not set")
		return
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	getContents(ctx, client, dirToDownload)
}

func getContents(ctx context.Context, client *github.Client, path string) {
	opts := &github.RepositoryContentGetOptions{}
	_, directoryContent, _, err := client.Repositories.GetContents(ctx, owner, repo, path, opts)
	if err != nil {
		fmt.Printf("Error getting contents: %s\n", err)
		return
	}

	for _, content := range directoryContent {
		localPath := filepath.Join(basePath, *content.Path)
		fmt.Println("Processing:", *content.Type, *content.Path)

		if *content.Type == "file" {
			downloadFile(ctx, client, content, localPath)
		} else if *content.Type == "dir" {
			getContents(ctx, client, *content.Path)
		}
	}
}

func downloadFile(ctx context.Context, client *github.Client, content *github.RepositoryContent, localPath string) {
	rc, err := client.Repositories.DownloadContents(ctx, owner, repo, *content.Path, nil)
	if err != nil {
		fmt.Printf("Error downloading file: %s\n", err)
		return
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		fmt.Printf("Error reading content from repository: %s\n", err)
		return
	}

	if err = os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		fmt.Printf("Error creating directories: %s\n", err)
		return
	}

	fmt.Println("Writing file:", localPath)
	err = ioutil.WriteFile(localPath, data, 0644)
	if err != nil {
		fmt.Printf("Error writing file to disk: %s\n", err)
	}
}
