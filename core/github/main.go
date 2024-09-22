package github

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// DownloadChallengeFiles downloads CTF challenge files from a GitHub repository
func DownloadChallengeFiles(token, owner, repo, dirToDownload, basePath string) error {
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
    tc := oauth2.NewClient(ctx, ts)
    client := github.NewClient(tc)

    return getContents(ctx, client, owner, repo, dirToDownload, basePath)
}

func getContents(ctx context.Context, client *github.Client, owner, repo, path, basePath string) error {
    opts := &github.RepositoryContentGetOptions{}
    _, directoryContent, _, err := client.Repositories.GetContents(ctx, owner, repo, path, opts)
    if err != nil {
        return fmt.Errorf("error getting contents: %w", err)
    }

    for _, content := range directoryContent {
        localPath := filepath.Join(basePath, *content.Path)
        fmt.Println("Processing:", *content.Type, *content.Path)
        if *content.Type == "file" {
            if err := downloadFile(ctx, client, owner, repo, content, localPath); err != nil {
                return err
            }
        } else if *content.Type == "dir" {
            if err := getContents(ctx, client, owner, repo, *content.Path, basePath); err != nil {
                return err
            }
        }
    }

    return nil
}

func downloadFile(ctx context.Context, client *github.Client, owner, repo string, content *github.RepositoryContent, localPath string) error {
    rc, err := client.Repositories.DownloadContents(ctx, owner, repo, *content.Path, nil)
    if err != nil {
        return fmt.Errorf("error downloading file: %w", err)
    }
    defer rc.Close()

    data, err := io.ReadAll(rc)
    if err != nil {
        return fmt.Errorf("error reading content from repository: %w", err)
    }

    if err = os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
        return fmt.Errorf("error creating directories: %w", err)
    }

    fmt.Println("Writing file:", localPath)
    return os.WriteFile(localPath, data, 0644)
}