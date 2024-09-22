package docker

import (
    "context"
    "fmt"
    "io"
    "os"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    "github.com/docker/docker/pkg/archive"
)

// CreateContainerImage creates a Docker image from a Dockerfile
func CreateContainerImage(dockerfilePath string) error {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return fmt.Errorf("failed to create Docker client: %w", err)
    }
    defer cli.Close()

    ctx := context.Background()

    tar, err := archive.TarWithOptions(dockerfilePath, &archive.TarOptions{})
    if err != nil {
        return fmt.Errorf("failed to create tar archive: %w", err)
    }

    opts := types.ImageBuildOptions{
        Dockerfile: "Dockerfile",
        Tags:       []string{"stellarpods-image:latest"},
    }

    resp, err := cli.ImageBuild(ctx, tar, opts)
    if err != nil {
        return fmt.Errorf("failed to build Docker image: %w", err)
    }
    defer resp.Body.Close()

    _, err = io.Copy(os.Stdout, resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read build response: %w", err)
    }

    fmt.Println("Docker image built successfully!")
    return nil
}