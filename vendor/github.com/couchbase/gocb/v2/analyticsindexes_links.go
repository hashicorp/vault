package gocb

import (
	"net/url"
	"strings"
)

// CouchbaseRemoteAnalyticsEncryptionSettings are the settings available for setting encryption level on a link.
type CouchbaseRemoteAnalyticsEncryptionSettings struct {
	// The level of encryption to apply, defaults to none
	EncryptionLevel AnalyticsEncryptionLevel
	// The certificate to use when encryption level is set to full.
	// This must be set if encryption level is set to full.
	Certificate []byte
	// The client certificate to use when encryption level is set to full.
	// This cannot be used if username and password are also used.
	ClientCertificate []byte
	// The client key to use when encryption level is set to full.
	// This cannot be used if username and password are also used.
	ClientKey []byte
}

// CouchbaseRemoteAnalyticsLink describes a remote analytics link which uses a Couchbase data service that is not part
// of the same cluster as the Analytics Service.
type CouchbaseRemoteAnalyticsLink struct {
	Dataverse  string
	LinkName   string
	Hostname   string
	Encryption CouchbaseRemoteAnalyticsEncryptionSettings
	Username   string
	Password   string
}

// Name returns the name of this link.
func (al *CouchbaseRemoteAnalyticsLink) Name() string {
	return al.LinkName
}

// DataverseName returns the name of the dataverse that this link belongs to.
func (al *CouchbaseRemoteAnalyticsLink) DataverseName() string {
	return al.Dataverse
}

// LinkType returns the type of analytics type this link is: AnalyticsLinkTypeCouchbaseRemote.
func (al *CouchbaseRemoteAnalyticsLink) LinkType() AnalyticsLinkType {
	return AnalyticsLinkTypeCouchbaseRemote
}

// FormEncode encodes the link into a form data representation, to sent as the body of a CreateLink or ReplaceLink
// request.
func (al *CouchbaseRemoteAnalyticsLink) FormEncode() ([]byte, error) {
	data := url.Values{}
	data.Add("hostname", al.Hostname)
	data.Add("type", string(AnalyticsLinkTypeCouchbaseRemote))
	if al.Username != "" {
		data.Add("username", al.Username)
	}
	if al.Password != "" {
		data.Add("password", al.Password)
	}
	if !strings.Contains(al.Dataverse, "/") {
		data.Add("dataverse", al.Dataverse)
		data.Add("name", al.LinkName)
	}

	data.Add("encryption", al.Encryption.EncryptionLevel.String())
	if len(al.Encryption.Certificate) > 0 {
		data.Add("certificate", string(al.Encryption.Certificate))
	}
	if len(al.Encryption.ClientCertificate) > 0 {
		data.Add("clientCertificate", string(al.Encryption.ClientCertificate))
	}
	if len(al.Encryption.ClientKey) > 0 {
		data.Add("clientKey", string(al.Encryption.ClientKey))
	}

	return []byte(data.Encode()), nil
}

// Validate is used by CreateLink and ReplaceLink to ensure that the link is valid.
func (al *CouchbaseRemoteAnalyticsLink) Validate() error {
	if al.Dataverse == "" {
		return makeInvalidArgumentsError("dataverse must be set for couchbase analytics links")
	}
	if al.LinkName == "" {
		return makeInvalidArgumentsError("name must be set for couchbase analytics links")
	}
	if al.Hostname == "" {
		return makeInvalidArgumentsError("hostname must be set for couchbase analytics links")
	}
	if al.Encryption.EncryptionLevel == AnalyticsEncryptionLevelFull && len(al.Encryption.Certificate) == 0 {
		return makeInvalidArgumentsError("when encryption level is full a certificate must be set for couchbase analytics links")
	}
	if (len(al.Encryption.ClientKey) > 0 && len(al.Encryption.ClientCertificate) == 0) ||
		(len(al.Encryption.ClientKey) == 0 && len(al.Encryption.ClientCertificate) > 0) {
		return makeInvalidArgumentsError("client certificate and client key must be set together for couchbase analytics links")
	}

	return nil
}

// S3ExternalAnalyticsLink describes an external analytics link which uses the AWS S3 service to access data.
type S3ExternalAnalyticsLink struct {
	Dataverse       string
	LinkName        string
	AccessKeyID     string
	SecretAccessKey string
	// SessionToken is only available in 7.0+.
	SessionToken    string
	Region          string
	ServiceEndpoint string
}

// Validate is used by CreateLink and ReplaceLink to ensure that the link is valid.
func (al *S3ExternalAnalyticsLink) Validate() error {
	if al.Dataverse == "" {
		return makeInvalidArgumentsError("dataverse must be set for s3 analytics links")
	}
	if al.LinkName == "" {
		return makeInvalidArgumentsError("name must be set for s3 analytics links")
	}
	if al.AccessKeyID == "" {
		return makeInvalidArgumentsError("access key id must be set for s3 analytics links")
	}
	if al.SecretAccessKey == "" {
		return makeInvalidArgumentsError("secret access key must be set for s3 analytics links")
	}
	if al.Region == "" {
		return makeInvalidArgumentsError("region must be set for s3 analytics links")
	}

	return nil
}

// Name returns the name of this link.
func (al *S3ExternalAnalyticsLink) Name() string {
	return al.LinkName
}

// DataverseName returns the name of the dataverse that this link belongs to.
func (al *S3ExternalAnalyticsLink) DataverseName() string {
	return al.Dataverse
}

// LinkType returns the type of analytics type this link is: AnalyticsLinkTypeS3External.
func (al *S3ExternalAnalyticsLink) LinkType() AnalyticsLinkType {
	return AnalyticsLinkTypeS3External
}

// FormEncode encodes the link into a form data representation, to sent as the body of a CreateLink or ReplaceLink
// request.
func (al *S3ExternalAnalyticsLink) FormEncode() ([]byte, error) {
	data := url.Values{}
	if !strings.Contains(al.Dataverse, "/") {
		data.Add("dataverse", al.Dataverse)
		data.Add("name", al.LinkName)
	}
	data.Add("type", string(AnalyticsLinkTypeS3External))
	data.Add("accessKeyId", al.AccessKeyID)
	data.Add("secretAccessKey", al.SecretAccessKey)
	data.Add("region", al.Region)

	if al.SessionToken != "" {
		data.Add("sessionToken", al.SessionToken)
	}
	if al.ServiceEndpoint != "" {
		data.Add("serviceEndpoint", al.ServiceEndpoint)
	}

	return []byte(data.Encode()), nil
}

// AzureBlobExternalAnalyticsLink describes an external analytics link which uses the Microsoft Azure Blob Storage
// service.
// Only available as of 7.0 Developer Preview.
// VOLATILE: This API is subject to change at any time.
type AzureBlobExternalAnalyticsLink struct {
	Dataverse             string
	LinkName              string
	ConnectionString      string
	AccountName           string
	AccountKey            string
	SharedAccessSignature string
	BlobEndpoint          string
	EndpointSuffix        string
}

// Name returns the name of this link.
func (al *AzureBlobExternalAnalyticsLink) Name() string {
	return al.LinkName
}

// DataverseName returns the name of the dataverse that this link belongs to.
func (al *AzureBlobExternalAnalyticsLink) DataverseName() string {
	return al.Dataverse
}

// LinkType returns the type of analytics type this link is: AnalyticsLinkTypeAzureExternal.
func (al *AzureBlobExternalAnalyticsLink) LinkType() AnalyticsLinkType {
	return AnalyticsLinkTypeAzureExternal
}

// Validate is used by CreateLink and ReplaceLink to ensure that the link is valid.
func (al *AzureBlobExternalAnalyticsLink) Validate() error {
	if al.Dataverse == "" {
		return makeInvalidArgumentsError("dataverse must be set for azureblob analytics links")
	}
	if al.LinkName == "" {
		return makeInvalidArgumentsError("name must be set for azureblob analytics links")
	}

	return nil
}

// FormEncode encodes the link into a form data representation, to sent as the body of a CreateLink or ReplaceLink
// request.
func (al *AzureBlobExternalAnalyticsLink) FormEncode() ([]byte, error) {
	data := url.Values{}
	data.Add("type", string(AnalyticsLinkTypeAzureExternal))

	if al.ConnectionString != "" {
		data.Add("connectionString", al.ConnectionString)
	}
	if al.AccountName != "" {
		data.Add("accountName", al.AccountName)
	}
	if al.AccountKey != "" {
		data.Add("accountKey", al.AccountKey)
	}
	if al.SharedAccessSignature != "" {
		data.Add("sharedAccessSignature", al.SharedAccessSignature)
	}
	if al.BlobEndpoint != "" {
		data.Add("blobEndpoint", al.BlobEndpoint)
	}
	if al.EndpointSuffix != "" {
		data.Add("endpointSuffix", al.EndpointSuffix)
	}

	return []byte(data.Encode()), nil
}
