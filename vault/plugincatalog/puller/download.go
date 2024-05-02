// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package puller

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
)

const (
	releasesAPI = "https://api.releases.hashicorp.com"
)

type DownloadPluginInput struct {
	Directory string
	Command   string
	Version   string
	SHA256Sum []byte // Optional. Verified if set.
}

func (in DownloadPluginInput) targetFile() string {
	return filepath.Join(in.Directory, fmt.Sprintf("%s-%s-%s", in.Command, in.Version, hex.EncodeToString(in.SHA256Sum)))
}

// EnsurePluginDownloaded downloads the plugin if it doesn't exist at the target
// location and returns the SHA256 sum of the plugin binary.
func EnsurePluginDownloaded(ctx context.Context, sources []pluginSource, in DownloadPluginInput) (ver string, sha256Sum []byte, err error) {
	// TODO: Per-plugin locking.

	// TODO: Allow the releases API to be omitted.
	sources = append([]pluginSource{newHTTPPluginSource(releasesAPI)}, sources...)

	if in.Version == "" {
		for _, p := range sources {
			metadata, err := p.listMetadata(ctx, in.Command)
			if err != nil {
				if errors.Is(err, errNotFound) {
					// TODO: Info/debug log.
					continue
				}

				return "", nil, fmt.Errorf("failed to list versions for plugin %s: %w", in.Command, err)
			}

			if len(metadata) == 0 {
				continue
			}

			// Pick the latest available version.
			var versions []*version.Version
			for _, m := range metadata {
				v, err := version.NewVersion(m.Version)
				if err != nil {
					return "", nil, err
				}
				versions = append(versions, v)
			}
			sort.Sort(version.Collection(versions))
			in.Version = "v" + versions[len(versions)-1].String()
			break
		}

		if in.Version == "" {
			return "", nil, fmt.Errorf("plugin %s not found from any available sources", in.Command)
		}
	}

	if exists, err := checkExisting(in); exists || err != nil {
		return in.Version, in.SHA256Sum, err
	}

	for _, p := range sources {
		metadata, err := p.getMetadata(ctx, in.Command, strings.TrimPrefix(in.Version, "v"))
		if err != nil {
			if errors.Is(err, errNotFound) {
				// TODO: Info/debug log.
				continue
			}

			return "", nil, err
		}

		sha256Sum, err = downloadPlugin(ctx, p, in, metadata)
		if err != nil {
			return "", nil, fmt.Errorf("failed to download plugin %s: %w", in.Command, err)
		}

		return in.Version, sha256Sum, nil
	}

	return "", nil, fmt.Errorf("plugin %s not found from any available sources", in.Command)
}

func downloadPlugin(ctx context.Context, p pluginSource, in DownloadPluginInput, metadata metadata) (sha256Sum []byte, err error) {
	// Get the SHA256SUMS file.
	zipSumsReader, err := p.getContentReader(ctx, metadata.URLSHASums)
	if err != nil {
		return nil, err
	}
	defer zipSumsReader.Close()
	zipSums, err := io.ReadAll(zipSumsReader)
	if err != nil {
		return nil, err
	}

	// Verify the SHA256SUMS file with any good signature.
	var verifyErrs error
	for _, sigURL := range metadata.URLSHASumsSignatures {
		sigReader, err := p.getContentReader(ctx, sigURL)
		if err != nil {
			return nil, err
		}
		defer sigReader.Close()
		sig, err := io.ReadAll(sigReader)
		if err != nil {
			return nil, err
		}

		if err := verifySignature(zipSums, sig); err != nil {
			verifyErrs = errors.Join(verifyErrs, fmt.Errorf("error verifying signature from %s: %w", sigURL, err))
		} else {
			break
		}
	}
	if verifyErrs != nil {
		return nil, verifyErrs
	}

	// Get the zip file.
	var build *build
	for _, b := range metadata.Builds {
		if b.Arch == runtime.GOARCH && b.OS == runtime.GOOS {
			build = &b
			break
		}
	}
	if build == nil {
		return nil, fmt.Errorf("no %s build available for %s/%s", in.Command, runtime.GOOS, runtime.GOARCH)
	}

	zipReader, err := p.getContentReader(ctx, build.URL)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()
	tempZipFile, zipFileSum, err := getAsFileAndHash(zipReader)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempZipFile)

	// Verify our zip file matches the expected SHA256 sum.
	found := false
	zipName := filepath.Base(build.URL)
	for _, zipSumLine := range strings.Split(string(zipSums), "\n") {
		expectedZipSum, expectedZipName, valid := strings.Cut(zipSumLine, " ")
		if !valid {
			continue
		}
		if strings.TrimSpace(expectedZipName) == zipName {
			if expectedZipSum == hex.EncodeToString(zipFileSum) {
				found = true
				break
			} else {
				// TODO: User error instead of server error?
				return nil, fmt.Errorf("expected SHA256 sum %s, got %s", expectedZipSum, zipFileSum)
			}
		}
	}
	if !found {
		return nil, fmt.Errorf("missing entry for %s in SHA256SUMS", zipName)
	}

	// Verify before we unzip so we never write a bad plugin to the plugin
	// directory and don't have to clean up an additional temp file if it fails.
	pluginSum, err := sha256SumFromZip(tempZipFile, in.Command)
	if err != nil {
		return nil, err
	}

	// We don't have the SHA256 sum available yet if we're downloading during
	// registration.
	if in.SHA256Sum != nil {
		if !slices.Equal(pluginSum, in.SHA256Sum) {
			return nil, fmt.Errorf("expected SHA256 sum %s, got %s", hex.EncodeToString(in.SHA256Sum), hex.EncodeToString(pluginSum))
		}
	} else {
		in.SHA256Sum = pluginSum
	}

	// Unzip to the target in the plugin directory.
	zipFileReader, err := zip.OpenReader(tempZipFile)
	if err != nil {
		return nil, err
	}
	defer zipFileReader.Close()
	pluginReader, err := zipFileReader.Open(in.Command)
	if err != nil {
		return nil, err
	}
	defer pluginReader.Close()

	f, err := os.OpenFile(in.targetFile(), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o700)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(f, pluginReader); err != nil {
		return nil, err
	}

	return pluginSum, nil
}

func checkExisting(in DownloadPluginInput) (bool, error) {
	if in.SHA256Sum == nil {
		return false, nil
	}

	_, err := os.Stat(in.targetFile())
	switch {
	case err == nil:
		// The file exists, double check the SHA256 sum is correct. This ensures
		// that if anyone modifies the file after initial download, Vault will
		// self-heal if the plugin is still available for download.
		hasher := sha256.New()
		f, err := os.Open(in.targetFile())
		if err != nil {
			return false, err
		}
		defer f.Close()
		if _, err := io.Copy(hasher, f); err != nil {
			return false, err
		}
		fileSum := hasher.Sum(nil)
		if !slices.Equal(fileSum, in.SHA256Sum) {
			return false, nil
		}
		return true, nil
	case os.IsNotExist(err):
		return false, nil
	default:
		return false, fmt.Errorf("error checking if managed plugin file exists: %w", err)
	}
}

func sha256SumFromZip(zipFileName, pluginFileName string) ([]byte, error) {
	zipReader, err := zip.OpenReader(zipFileName)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()
	pluginFile, err := zipReader.Open(pluginFileName)
	if err != nil {
		return nil, err
	}
	defer pluginFile.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, pluginFile); err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}

func getAsFileAndHash(reader io.Reader) (fileName string, shaSum []byte, retErr error) {
	f, err := os.CreateTemp(os.TempDir(), "vault-plugin-temp")
	if err != nil {
		return "", nil, err
	}
	defer func() {
		retErr = errors.Join(retErr, f.Close())
		if retErr != nil {
			retErr = errors.Join(retErr, os.Remove(f.Name()))
		}
	}()

	hasher := sha256.New()
	zipReader := io.TeeReader(reader, hasher)
	if _, err := io.Copy(f, zipReader); err != nil {
		return "", nil, err
	}

	return f.Name(), hasher.Sum(nil), nil
}
