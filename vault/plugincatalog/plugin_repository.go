package plugincatalog

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

// TODO create an interface for the different types of plugin repositories

func downloadPluginBinary(pluginDirectory string, plugin pluginutil.SetPluginInput) error {
	filePath := filepath.Join(pluginDirectory, plugin.Name+"_"+plugin.Version)

	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	if err := os.Chmod(filePath, 0o775); err != nil {
		return fmt.Errorf("failed to set file permissions: %v", err)
	}

	// TODO construct the URL from the plugin metadata
	resp, err := http.Get("https://releases.hashicorp.com/vault-plugin-database-elasticsearch/0.17.0/vault-plugin-database-elasticsearch_0.17.0_linux_amd64.zip")
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: received status code %d", resp.StatusCode)
	}

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	fmt.Println("File downloaded successfully:", filePath)
	return nil
}

func (c *PluginCatalog) downloadPluginImage(ociImage string) error {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error creating Docker client: %w", err)
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctxWithTimeout, ociImage, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	defer reader.Close()

	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return fmt.Errorf("error reading pull output: %w", err)
	}

	return nil
}
