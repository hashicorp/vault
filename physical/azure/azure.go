// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package azure

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/permitpool"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/physical"
)

const (
	// MaxBlobSize at this time
	MaxBlobSize = 1024 * 1024 * 4
	// MaxListResults is the current default value, setting explicitly
	MaxListResults = 5000
)

// AzureBackend is a physical backend that stores data
// within an Azure blob container.
type AzureBackend struct {
	container  *azblob.ContainerURL
	logger     log.Logger
	permitPool *permitpool.Pool
}

// Verify AzureBackend satisfies the correct interfaces
var _ physical.Backend = (*AzureBackend)(nil)

// NewAzureBackend constructs an Azure backend using a pre-existing
// bucket. Credentials can be provided to the backend, sourced
// from the environment, via HCL or by using managed identities.
func NewAzureBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	name := os.Getenv("AZURE_BLOB_CONTAINER")
	useMSI := false

	if name == "" {
		name = conf["container"]
		if name == "" {
			return nil, fmt.Errorf("'container' must be set")
		}
	}

	if err := validateContainerName(name); err != nil {
		return nil, fmt.Errorf("invalid container name %s: %w", name, err)
	}

	accountName := os.Getenv("AZURE_ACCOUNT_NAME")
	if accountName == "" {
		accountName = conf["accountName"]
		if accountName == "" {
			return nil, fmt.Errorf("'accountName' must be set")
		}
	}
	if err := validateAccountName(accountName); err != nil {
		return nil, fmt.Errorf("invalid account name %s: %w", accountName, err)
	}

	accountKey := os.Getenv("AZURE_ACCOUNT_KEY")
	if accountKey == "" {
		accountKey = conf["accountKey"]
		if accountKey == "" {
			logger.Info("accountKey not set, using managed identity auth")
			useMSI = true
		}
	}

	environmentName := os.Getenv("AZURE_ENVIRONMENT")
	if environmentName == "" {
		environmentName = conf["environment"]
		if environmentName == "" {
			environmentName = "AzurePublicCloud"
		}
	}

	environmentURL := os.Getenv("AZURE_ARM_ENDPOINT")
	if environmentURL == "" {
		environmentURL = conf["arm_endpoint"]
	}

	var environment azure.Environment
	var URL *url.URL
	var err error

	testHost := conf["testHost"]
	switch {
	case testHost != "":
		URL = &url.URL{Scheme: "http", Host: testHost, Path: fmt.Sprintf("/%s/%s", accountName, name)}
	default:
		if environmentURL != "" {
			environment, err = azure.EnvironmentFromURL(environmentURL)
			if err != nil {
				return nil, fmt.Errorf("failed to look up Azure environment descriptor for URL %q: %w", environmentURL, err)
			}
		} else {
			environment, err = azure.EnvironmentFromName(environmentName)
			if err != nil {
				return nil, fmt.Errorf("failed to look up Azure environment descriptor for name %q: %w", environmentName, err)
			}
		}
		URL, err = url.Parse(
			fmt.Sprintf("https://%s.blob.%s/%s", accountName, environment.StorageEndpointSuffix, name))
		if err != nil {
			return nil, fmt.Errorf("failed to create Azure client: %w", err)
		}
	}

	var credential azblob.Credential
	if useMSI {
		authToken, err := getAuthTokenFromIMDS(environment.ResourceIdentifiers.Storage)
		if err != nil {
			return nil, fmt.Errorf("failed to obtain auth token from IMDS %q: %w", environmentName, err)
		}

		credential = azblob.NewTokenCredential(authToken.OAuthToken(), func(c azblob.TokenCredential) time.Duration {
			err = authToken.Refresh()
			if err != nil {
				logger.Error("couldn't refresh token credential", "error", err)
				return 0
			}

			expIn, err := authToken.Token().ExpiresIn.Int64()
			if err != nil {
				logger.Error("couldn't retrieve jwt claim for 'expiresIn' from refreshed token", "error", err)
				return 0
			}

			logger.Debug("token refreshed, new token expires in", "access_token_expiry", expIn)
			c.SetToken(authToken.OAuthToken())

			// tokens are valid for 23h59m (86399s) by default, refresh after ~21h
			return time.Duration(int(float64(expIn)*0.9)) * time.Second
		})
	} else {
		credential, err = azblob.NewSharedKeyCredential(accountName, accountKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create Azure client: %w", err)
		}
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	containerURL := azblob.NewContainerURL(*URL, p)
	_, err = containerURL.GetProperties(ctx, azblob.LeaseAccessConditions{})
	if err != nil {
		var e azblob.StorageError
		if errors.As(err, &e) {
			switch e.ServiceCode() {
			case azblob.ServiceCodeContainerNotFound:
				_, err := containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
				if err != nil {
					return nil, fmt.Errorf("failed to create %q container: %w", name, err)
				}
			default:
				return nil, fmt.Errorf("failed to get properties for container %q: %w", name, err)
			}
		}
	}

	maxParStr, ok := conf["max_parallel"]
	var maxParInt int
	if ok {
		maxParInt, err = strconv.Atoi(maxParStr)
		if err != nil {
			return nil, fmt.Errorf("failed parsing max_parallel parameter: %w", err)
		}
		if logger.IsDebug() {
			logger.Debug("max_parallel set", "max_parallel", maxParInt)
		}
	}

	a := &AzureBackend{
		container:  &containerURL,
		logger:     logger,
		permitPool: permitpool.New(maxParInt),
	}
	return a, nil
}

// validation rules for containers are defined here:
// https://learn.microsoft.com/en-us/rest/api/storageservices/Naming-and-Referencing-Containers--Blobs--and-Metadata#container-names
var containerNameRegex = regexp.MustCompile("^[a-z0-9]+(-[a-z0-9]+)*$")

func validateContainerName(name string) error {
	if len(name) < 3 || len(name) > 63 {
		return errors.New("name must be between 3 and 63 characters long")
	}

	if !containerNameRegex.MatchString(name) {
		return errors.New("name is invalid")
	}
	return nil
}

// validation rules are defined here:
// https://learn.microsoft.com/en-us/azure/azure-resource-manager/troubleshooting/error-storage-account-name?tabs=bicep#cause
var accountNameRegex = regexp.MustCompile("^[a-z0-9]+$")

func validateAccountName(name string) error {
	if len(name) < 3 || len(name) > 24 {
		return errors.New("name must be between 3 and 24 characters long")
	}
	if !accountNameRegex.MatchString(name) {
		return errors.New("name is invalid")
	}
	return nil
}

// Put is used to insert or update an entry
func (a *AzureBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"azure", "put"}, time.Now())

	if len(entry.Value) >= MaxBlobSize {
		return fmt.Errorf("value is bigger than the current supported limit of 4MBytes")
	}

	if err := a.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer a.permitPool.Release()

	blobURL := a.container.NewBlockBlobURL(entry.Key)
	_, err := azblob.UploadBufferToBlockBlob(ctx, entry.Value, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize: MaxBlobSize,
	})

	return err
}

// Get is used to fetch an entry
func (a *AzureBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"azure", "get"}, time.Now())

	if err := a.permitPool.Acquire(ctx); err != nil {
		return nil, err
	}
	defer a.permitPool.Release()

	blobURL := a.container.NewBlockBlobURL(key)
	clientOptions := azblob.ClientProvidedKeyOptions{}

	res, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, clientOptions)
	if err != nil {
		var e azblob.StorageError
		if errors.As(err, &e) {
			switch e.ServiceCode() {
			case azblob.ServiceCodeBlobNotFound:
				return nil, nil
			default:
				return nil, fmt.Errorf("failed to download blob %q: %w", key, err)
			}
		}
		return nil, err
	}

	reader := res.Body(azblob.RetryReaderOptions{})
	defer reader.Close()

	data, err := io.ReadAll(reader)

	ent := &physical.Entry{
		Key:   key,
		Value: data,
	}

	return ent, err
}

// Delete is used to permanently delete an entry
func (a *AzureBackend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"azure", "delete"}, time.Now())

	if err := a.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer a.permitPool.Release()

	blobURL := a.container.NewBlockBlobURL(key)
	_, err := blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
	if err != nil {
		var e azblob.StorageError
		if errors.As(err, &e) {
			switch e.ServiceCode() {
			case azblob.ServiceCodeBlobNotFound:
				return nil
			default:
				return fmt.Errorf("failed to delete blob %q: %w", key, err)
			}
		}
	}

	return err
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (a *AzureBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"azure", "list"}, time.Now())

	if err := a.permitPool.Acquire(ctx); err != nil {
		return nil, err
	}
	defer a.permitPool.Release()

	var keys []string
	for marker := (azblob.Marker{}); marker.NotDone(); {
		listBlob, err := a.container.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{
			Prefix:     prefix,
			MaxResults: MaxListResults,
		})
		if err != nil {
			return nil, err
		}

		for _, blobInfo := range listBlob.Segment.BlobItems {
			key := strings.TrimPrefix(blobInfo.Name, prefix)
			if i := strings.Index(key, "/"); i == -1 {
				// file
				keys = append(keys, key)
			} else {
				// subdirectory
				keys = strutil.AppendIfMissing(keys, key[:i+1])
			}
		}

		marker = listBlob.NextMarker
	}

	sort.Strings(keys)
	return keys, nil
}

// getAuthTokenFromIMDS uses the Azure Instance Metadata Service to retrieve a short-lived credential using OAuth
// more info on this https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview
func getAuthTokenFromIMDS(resource string) (*adal.ServicePrincipalToken, error) {
	msiEndpoint, err := adal.GetMSIEndpoint()
	if err != nil {
		return nil, err
	}

	spt, err := adal.NewServicePrincipalTokenFromMSI(msiEndpoint, resource)
	if err != nil {
		return nil, err
	}

	if err := spt.Refresh(); err != nil {
		return nil, err
	}

	token := spt.Token()
	if token.IsZero() {
		return nil, err
	}

	return spt, nil
}
