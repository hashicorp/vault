package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
)

// A KVClient is used to perform reads and writes against a KV secrets engine in Vault.
//
// The mount path is the location where the target KV secrets engine resides
// in Vault. The version refers to the version of the target KV secrets engine.
//
// Vault development servers tend to have "secret" as the mount path
// and 2 as the version, as these are the default settings when a server
// is started in -dev mode.
//
// Learn more about the KV secrets engine here:
// https://www.vaultproject.io/docs/secrets/kv
type KVClient struct {
	c         *Client
	mountPath string
	version   int
}

type KVSecret struct {
	Data     map[string]interface{}
	Metadata *VersionMetadata
	Raw      *Secret
}

type KVMetadata struct {
	CASRequired        bool                   `mapstructure:"cas_required"`
	CreatedTime        time.Time              `mapstructure:"created_time"`
	CurrentVersion     int                    `mapstructure:"current_version"`
	CustomMetadata     map[string]interface{} `mapstructure:"custom_metadata"`
	DeleteVersionAfter time.Duration          `mapstructure:"delete_version_after"`
	MaxVersions        int                    `mapstructure:"max_versions"`
	OldestVersion      int                    `mapstructure:"oldest_version"`
	UpdatedTime        time.Time              `mapstructure:"updated_time"`
	// Keys are stringified ints, e.g. "3"
	Versions map[string]VersionMetadata `mapstructure:"versions"`
}

type VersionMetadata struct {
	Version      int       `mapstructure:"version"`
	CreatedTime  time.Time `mapstructure:"created_time"`
	DeletionTime time.Time `mapstructure:"deletion_time"`
	Destroyed    bool      `mapstructure:"destroyed"`
	// There is currently no version-specific custom metadata.
	// This field is just a copy of what's in the CustomMetadata field
	// for the full KVMetadata of the secret.
	CustomMetadata map[string]string `mapstructure:"custom_metadata"`
}

func (c *Client) KVv1(mountPath string) *KVClient {
	return &KVClient{c: c, mountPath: mountPath, version: 1}
}

func (c *Client) KVv2(mountPath string) *KVClient {
	return &KVClient{c: c, mountPath: mountPath, version: 2}
}

func (kv *KVClient) Read(ctx context.Context, secretPath string) (*KVSecret, error) {
	pathToRead, err := kv.getFullPath(secretPath)
	if err != nil {
		return nil, fmt.Errorf("error assembling full path to KV secret: %v", err)
	}

	secret, err := kv.c.Logical().ReadWithContext(ctx, pathToRead)
	if err != nil {
		return nil, fmt.Errorf("error encountered while reading secret at %s: %v", pathToRead, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no secret found at %s", pathToRead)
	}

	return extractDataAndVersionMetadata(secret, kv.version)
}

// ReadVersion returns the data and metadata for a specific version of the
// given secret. If that version has been deleted, the Data field on the
// returned secret will be nil, and the Metadata field will contain the deletion time.
func (kv *KVClient) ReadVersion(ctx context.Context, secretPath string, version int) (*KVSecret, error) {
	pathToRead, err := kv.getFullPath(secretPath)
	if err != nil {
		return nil, fmt.Errorf("error assembling full path to KV secret: %v", err)
	}

	queryParams := map[string][]string{"version": {strconv.Itoa(version)}}
	secret, err := kv.c.Logical().ReadWithDataWithContext(ctx, pathToRead, queryParams)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no secret version found at %s", pathToRead)
	}

	return extractDataAndVersionMetadata(secret, kv.version)
}

func (kv *KVClient) ReadMetadata(ctx context.Context, secretPath string) (*KVMetadata, error) {
	pathToRead, err := kv.getFullMetadataPath(secretPath)
	if err != nil {
		return nil, fmt.Errorf("error assembling full path to KV secret's metadata: %v", err)
	}

	secret, err := kv.c.Logical().ReadWithContext(ctx, pathToRead)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no secret metadata found at %s", pathToRead)
	}

	return extractFullMetadata(secret)
}

func (kv *KVClient) getFullPath(secretPath string) (string, error) {
	var pathToRead string
	switch kv.version {
	case 1:
		pathToRead = fmt.Sprintf("%s/%s", kv.mountPath, secretPath)
	case 2:
		pathToRead = fmt.Sprintf("%s/data/%s", kv.mountPath, secretPath)
	default:
		return "", fmt.Errorf("KV client was initialized with invalid KV secrets engine version")
	}

	return pathToRead, nil
}

func (kv *KVClient) getFullMetadataPath(secretPath string) (string, error) {
	var pathToRead string
	switch kv.version {
	case 1:
		return "", fmt.Errorf("metadata is not supported in v1 of the KV secrets engine")
	case 2:
		pathToRead = fmt.Sprintf("%s/metadata/%s", kv.mountPath, secretPath)
	default:
		return "", fmt.Errorf("KV client was initialized with invalid KV secrets engine version")
	}

	return pathToRead, nil
}

func extractDataAndVersionMetadata(secret *Secret, version int) (*KVSecret, error) {
	var data map[string]interface{}
	var metadata *VersionMetadata
	switch version {
	case 1:
		data = secret.Data
		metadata = nil
	case 2:
		dataInterface, ok := secret.Data["data"]
		if !ok {
			return nil, fmt.Errorf("missing expected 'data' element")
		}

		if dataInterface == nil {
			// this can happen when the secret has been deleted, but the metadata is still available
			data = nil
		} else {
			data, ok = dataInterface.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("unexpected type for 'data' element: %T (%#v)", data, data)
			}
		}

		metadataInterface, ok := secret.Data["metadata"]
		if !ok {
			return nil, fmt.Errorf("missing expected 'metadata' element")
		}

		metadataMap, ok := metadataInterface.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected type for 'metadata' element: %T (%#v)", metadata, metadata)
		}

		// deletion_time usually comes in as an empty string which can't be
		// processed as time.RFC3339, so we reset it to a convertible value
		if metadataMap["deletion_time"] == "" {
			metadataMap["deletion_time"] = time.Time{}
		}

		d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			DecodeHook: mapstructure.StringToTimeHookFunc(time.RFC3339),
			Result:     &metadata,
		})
		if err != nil {
			return nil, fmt.Errorf("error setting up decoder for API response: %v", err)
		}

		err = d.Decode(metadataMap)
		if err != nil {
			return nil, fmt.Errorf("error decoding metadata from API response into VersionMetadata: %v", err)
		}
	default:
		return nil, fmt.Errorf("cannot parse secret without specifying valid KV secrets engine version")
	}

	return &KVSecret{
		Data:     data,
		Metadata: metadata,
		Raw:      secret,
	}, nil
}

func extractFullMetadata(secret *Secret) (*KVMetadata, error) {
	var metadata *KVMetadata

	// deletion_time usually comes in as an empty string which can't be
	// processed as time.RFC3339, so we reset it to a convertible value
	if versions, ok := secret.Data["versions"]; ok {
		versionsMap := versions.(map[string]interface{})
		if len(versionsMap) > 0 {
			for version, metadata := range versionsMap {
				metadataMap := metadata.(map[string]interface{})
				if metadataMap["deletion_time"] == "" {
					metadataMap["deletion_time"] = time.Time{}
				}
				versionsMap[version] = metadataMap // save the updated copy of the metadata map
			}
		}
		secret.Data["versions"] = versionsMap // save the updated copy of the versions map
	}

	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.StringToTimeDurationHookFunc(),
		),
		Result: &metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("error setting up decoder for API response: %v", err)
	}

	err = d.Decode(secret.Data)
	if err != nil {
		return nil, fmt.Errorf("error decoding metadata from API response into KVMetadata: %v", err)
	}

	return metadata, nil
}
