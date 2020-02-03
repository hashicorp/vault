//
// Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package storage

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/joyent/triton-go/client"
	tt "github.com/joyent/triton-go/errors"
	"github.com/pkg/errors"
)

type ObjectsClient struct {
	client *client.Client
}

// AbortMpuInput represents parameters to an AbortMpu operation
type AbortMpuInput struct {
	PartsDirectoryPath string
}

func (s *ObjectsClient) AbortMultipartUpload(ctx context.Context, input *AbortMpuInput) error {
	return abortMpu(*s, ctx, input)
}

// CommitMpuInput represents parameters to a CommitMpu operation
type CommitMpuInput struct {
	Id      string
	Headers map[string]string
	Body    CommitMpuBody
}

// CommitMpuBody represents the body of a CommitMpu request
type CommitMpuBody struct {
	Parts []string `json:"parts"`
}

func (s *ObjectsClient) CommitMultipartUpload(ctx context.Context, input *CommitMpuInput) error {
	return commitMpu(*s, ctx, input)
}

// CreateMpuInput represents parameters to a CreateMpu operation.
type CreateMpuInput struct {
	Body            CreateMpuBody
	ContentLength   uint64
	ContentMD5      string
	DurabilityLevel uint64
	ForceInsert     bool //Force the creation of the directory tree
}

// CreateMpuOutput represents the response from a CreateMpu operation
type CreateMpuOutput struct {
	Id             string `json:"id"`
	PartsDirectory string `json:"partsDirectory"`
}

// CreateMpuBody represents the body of a CreateMpu request.
type CreateMpuBody struct {
	ObjectPath string            `json:"objectPath"`
	Headers    map[string]string `json:"headers,omitempty"`
}

func (s *ObjectsClient) CreateMultipartUpload(ctx context.Context, input *CreateMpuInput) (*CreateMpuOutput, error) {
	return createMpu(*s, ctx, input)
}

// GetObjectInput represents parameters to a GetObject operation.
type GetInfoInput struct {
	ObjectPath string
	Headers    map[string]string
}

// GetObjectOutput contains the outputs for a GetObject operation. It is your
// responsibility to ensure that the io.ReadCloser ObjectReader is closed.
type GetInfoOutput struct {
	ContentLength uint64
	ContentType   string
	LastModified  time.Time
	ContentMD5    string
	ETag          string
	Metadata      map[string]string
}

// GetInfo sends a HEAD request to an object in the Manta service. This function
// does not return a response body.
func (s *ObjectsClient) GetInfo(ctx context.Context, input *GetInfoInput) (*GetInfoOutput, error) {
	absPath := absFileInput(s.client.AccountName, input.ObjectPath)

	headers := &http.Header{}
	for key, value := range input.Headers {
		headers.Set(key, value)
	}

	reqInput := client.RequestInput{
		Method:  http.MethodHead,
		Path:    string(absPath),
		Headers: headers,
	}
	_, respHeaders, err := s.client.ExecuteRequestStorage(ctx, reqInput)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get info")
	}

	response := &GetInfoOutput{
		ContentType: respHeaders.Get("Content-Type"),
		ContentMD5:  respHeaders.Get("Content-MD5"),
		ETag:        respHeaders.Get("Etag"),
	}

	lastModified, err := time.Parse(time.RFC1123, respHeaders.Get("Last-Modified"))
	if err == nil {
		response.LastModified = lastModified
	}

	contentLength, err := strconv.ParseUint(respHeaders.Get("Content-Length"), 10, 64)
	if err == nil {
		response.ContentLength = contentLength
	}

	metadata := map[string]string{}
	for key, values := range respHeaders {
		if strings.HasPrefix(key, "m-") {
			metadata[key] = strings.Join(values, ", ")
		}
	}
	response.Metadata = metadata

	return response, nil
}

// IsDir is a convenience wrapper around the GetInfo function which takes an
// ObjectPath and returns a boolean whether or not the object is a directory
// type in Manta. Returns an error if GetInfo failed upstream for some reason.
func (s *ObjectsClient) IsDir(ctx context.Context, objectPath string) (bool, error) {
	info, err := s.GetInfo(ctx, &GetInfoInput{
		ObjectPath: objectPath,
	})
	if err != nil {
		return false, err
	}
	if info != nil {
		return strings.HasSuffix(info.ContentType, "type=directory"), nil
	}
	return false, nil
}

// GetObjectInput represents parameters to a GetObject operation.
type GetObjectInput struct {
	ObjectPath string
	Headers    map[string]string
}

// GetObjectOutput contains the outputs for a GetObject operation. It is your
// responsibility to ensure that the io.ReadCloser ObjectReader is closed.
type GetObjectOutput struct {
	ContentLength uint64
	ContentType   string
	LastModified  time.Time
	ContentMD5    string
	ETag          string
	Metadata      map[string]string
	ObjectReader  io.ReadCloser
}

// Get retrieves an object from the Manta service. If error is nil (i.e. the
// call returns successfully), it is your responsibility to close the
// io.ReadCloser named ObjectReader in the operation output.
func (s *ObjectsClient) Get(ctx context.Context, input *GetObjectInput) (*GetObjectOutput, error) {
	absPath := absFileInput(s.client.AccountName, input.ObjectPath)

	headers := &http.Header{}
	for key, value := range input.Headers {
		headers.Set(key, value)
	}

	reqInput := client.RequestInput{
		Method:  http.MethodGet,
		Path:    string(absPath),
		Headers: headers,
	}
	respBody, respHeaders, err := s.client.ExecuteRequestStorage(ctx, reqInput)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get object")
	}

	response := &GetObjectOutput{
		ContentType:  respHeaders.Get("Content-Type"),
		ContentMD5:   respHeaders.Get("Content-MD5"),
		ETag:         respHeaders.Get("Etag"),
		ObjectReader: respBody,
	}

	lastModified, err := time.Parse(time.RFC1123, respHeaders.Get("Last-Modified"))
	if err == nil {
		response.LastModified = lastModified
	}

	contentLength, err := strconv.ParseUint(respHeaders.Get("Content-Length"), 10, 64)
	if err == nil {
		response.ContentLength = contentLength
	}

	metadata := map[string]string{}
	for key, values := range respHeaders {
		if strings.HasPrefix(key, "m-") {
			metadata[key] = strings.Join(values, ", ")
		}
	}
	response.Metadata = metadata

	return response, nil
}

// DeleteObjectInput represents parameters to a DeleteObject operation.
type DeleteObjectInput struct {
	ObjectPath string
	Headers    map[string]string
}

// DeleteObject deletes an object.
func (s *ObjectsClient) Delete(ctx context.Context, input *DeleteObjectInput) error {
	absPath := absFileInput(s.client.AccountName, input.ObjectPath)

	headers := &http.Header{}
	for key, value := range input.Headers {
		headers.Set(key, value)
	}

	reqInput := client.RequestInput{
		Method:  http.MethodDelete,
		Path:    string(absPath),
		Headers: headers,
	}
	respBody, _, err := s.client.ExecuteRequestStorage(ctx, reqInput)
	if respBody != nil {
		defer respBody.Close()
	}
	if err != nil {
		return errors.Wrap(err, "unable to delete object")
	}

	return nil
}

// GetMpuInput represents parameters to a GetMpu operation
type GetMpuInput struct {
	PartsDirectoryPath string
}

type GetMpuHeaders struct {
	ContentLength int64  `json:"content-length"`
	ContentMd5    string `json:"content-md5"`
}

type GetMpuOutput struct {
	Id             string        `json:"id"`
	State          string        `json:"state"`
	PartsDirectory string        `json:"partsDirectory"`
	TargetObject   string        `json:"targetObject"`
	Headers        GetMpuHeaders `json:"headers"`
	NumCopies      int64         `json:"numCopies"`
	CreationTimeMs int64         `json:"creationTimeMs"`
}

func (s *ObjectsClient) GetMultipartUpload(ctx context.Context, input *GetMpuInput) (*GetMpuOutput, error) {
	return getMpu(*s, ctx, input)
}

type ListMpuPartsInput struct {
	Id string
}

type ListMpuPart struct {
	ETag       string
	PartNumber int
	Size       int64
}

type ListMpuPartsOutput struct {
	Parts []ListMpuPart
}

func (s *ObjectsClient) ListMultipartUploadParts(ctx context.Context, input *ListMpuPartsInput) (*ListMpuPartsOutput, error) {
	return listMpuParts(*s, ctx, input)
}

// PutObjectMetadataInput represents parameters to a PutObjectMetadata operation.
type PutObjectMetadataInput struct {
	ObjectPath  string
	ContentType string
	Metadata    map[string]string
}

// PutObjectMetadata allows you to overwrite the HTTP headers for an already
// existing object, without changing the data. Note this is an idempotent "replace"
// operation, so you must specify the complete set of HTTP headers you want
// stored on each request.
//
// You cannot change "critical" headers:
// 	- Content-Length
//	- Content-MD5
//	- Durability-Level
func (s *ObjectsClient) PutMetadata(ctx context.Context, input *PutObjectMetadataInput) error {
	absPath := absFileInput(s.client.AccountName, input.ObjectPath)
	query := &url.Values{}
	query.Set("metadata", "true")

	headers := &http.Header{}
	headers.Set("Content-Type", input.ContentType)
	for key, value := range input.Metadata {
		headers.Set(key, value)
	}

	reqInput := client.RequestInput{
		Method:  http.MethodPut,
		Path:    string(absPath),
		Query:   query,
		Headers: headers,
	}
	respBody, _, err := s.client.ExecuteRequestStorage(ctx, reqInput)
	if respBody != nil {
		defer respBody.Close()
	}
	if err != nil {
		return errors.Wrap(err, "unable to put metadata")
	}

	return nil
}

// PutObjectInput represents parameters to a PutObject operation.
type PutObjectInput struct {
	ObjectPath       string
	DurabilityLevel  uint64
	ContentType      string
	ContentMD5       string
	IfMatch          string
	IfModifiedSince  *time.Time
	ContentLength    uint64
	MaxContentLength uint64
	ObjectReader     io.Reader
	Headers          map[string]string
	ForceInsert      bool //Force the creation of the directory tree
}

func (s *ObjectsClient) Put(ctx context.Context, input *PutObjectInput) error {
	absPath := absFileInput(s.client.AccountName, input.ObjectPath)
	if input.ForceInsert {
		absDirName := _AbsCleanPath(path.Dir(string(absPath)))
		exists, err := checkDirectoryTreeExists(*s, ctx, absDirName)
		if err != nil {
			return err
		}
		if !exists {
			err := createDirectory(*s, ctx, absDirName)
			if err != nil {
				return err
			}
			return putObject(*s, ctx, input, absPath)
		}
	}

	return putObject(*s, ctx, input, absPath)
}

// UploadPartInput represents parameters to a UploadPart operation.
type UploadPartInput struct {
	Id           string
	PartNum      uint64
	ContentMD5   string
	Headers      map[string]string
	ObjectReader io.Reader
}

// UploadPartOutput represents the response from a
type UploadPartOutput struct {
	Part string `json:"part"`
}

func (s *ObjectsClient) UploadPart(ctx context.Context, input *UploadPartInput) (*UploadPartOutput, error) {
	return uploadPart(*s, ctx, input)
}

// _AbsCleanPath is an internal type that means the input has been
// path.Clean()'ed and is an absolute path.
type _AbsCleanPath string

func absFileInput(accountName, objPath string) _AbsCleanPath {
	cleanInput := path.Clean(objPath)
	if strings.HasPrefix(cleanInput, path.Join("/", accountName, "/")) {
		return _AbsCleanPath(cleanInput)
	}

	cleanAbs := path.Clean(path.Join("/", accountName, objPath))
	return _AbsCleanPath(cleanAbs)
}

func putObject(c ObjectsClient, ctx context.Context, input *PutObjectInput, absPath _AbsCleanPath) error {
	if input.MaxContentLength != 0 && input.ContentLength != 0 {
		return errors.New("ContentLength and MaxContentLength may not both be set to non-zero values.")
	}

	headers := &http.Header{}
	for key, value := range input.Headers {
		headers.Set(key, value)
	}
	if input.DurabilityLevel != 0 {
		headers.Set("Durability-Level", strconv.FormatUint(input.DurabilityLevel, 10))
	}
	if input.ContentType != "" {
		headers.Set("Content-Type", input.ContentType)
	}
	if input.ContentMD5 != "" {
		headers.Set("Content-MD$", input.ContentMD5)
	}
	if input.IfMatch != "" {
		headers.Set("If-Match", input.IfMatch)
	}
	if input.IfModifiedSince != nil {
		headers.Set("If-Modified-Since", input.IfModifiedSince.Format(time.RFC1123))
	}
	if input.ContentLength != 0 {
		headers.Set("Content-Length", strconv.FormatUint(input.ContentLength, 10))
	}
	if input.MaxContentLength != 0 {
		headers.Set("Max-Content-Length", strconv.FormatUint(input.MaxContentLength, 10))
	}

	reqInput := client.RequestNoEncodeInput{
		Method:  http.MethodPut,
		Path:    string(absPath),
		Headers: headers,
		Body:    input.ObjectReader,
	}
	respBody, _, err := c.client.ExecuteRequestNoEncode(ctx, reqInput)
	if respBody != nil {
		defer respBody.Close()
	}
	if err != nil {
		return errors.Wrap(err, "unable to put object")
	}

	return nil
}

func createDirectory(c ObjectsClient, ctx context.Context, absPath _AbsCleanPath) error {
	dirClient := &DirectoryClient{
		client: c.client,
	}

	// An abspath starts w/ a leading "/" which gets added to the slice as an
	// empty string. Start all array math at 1.
	parts := strings.Split(string(absPath), "/")
	if len(parts) < 2 {
		return errors.New("no path components to create directory")
	}

	folderPath := parts[1]
	// Don't attempt to create a manta account as a directory
	for i := 2; i < len(parts); i++ {
		part := parts[i]
		folderPath = path.Clean(path.Join("/", folderPath, part))
		err := dirClient.Put(ctx, &PutDirectoryInput{
			DirectoryName: folderPath,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func abortMpu(c ObjectsClient, ctx context.Context, input *AbortMpuInput) error {
	reqInput := client.RequestInput{
		Method:  http.MethodPost,
		Path:    input.PartsDirectoryPath + "/abort",
		Headers: &http.Header{},
		Body:    nil,
	}
	respBody, _, err := c.client.ExecuteRequestStorage(ctx, reqInput)
	if err != nil {
		return errors.Wrap(err, "unable to abort mpu")
	}

	if respBody != nil {
		defer respBody.Close()
	}

	return nil
}

func commitMpu(c ObjectsClient, ctx context.Context, input *CommitMpuInput) error {
	headers := &http.Header{}
	for key, value := range input.Headers {
		headers.Set(key, value)
	}

	// The mpu directory prefix length is derived from the final character
	// in the mpu identifier which we'll call P. The mpu prefix itself is
	// the first P characters of the mpu identifier. In order to derive the
	// correct directory structure we need to parse this information from
	// the mpu identifier
	id := input.Id
	idLength := len(id)
	prefixLen, err := strconv.Atoi(id[idLength-1 : idLength])
	if err != nil {
		return errors.Wrap(err, "unable to commit mpu due to invalid mpu prefix length")
	}
	prefix := id[:prefixLen]
	partPath := "/" + c.client.AccountName + "/uploads/" + prefix + "/" + input.Id + "/commit"

	reqInput := client.RequestInput{
		Method:  http.MethodPost,
		Path:    partPath,
		Headers: headers,
		Body:    input.Body,
	}
	respBody, _, err := c.client.ExecuteRequestStorage(ctx, reqInput)
	if err != nil {
		return errors.Wrap(err, "unable to commit mpu")
	}

	if respBody != nil {
		defer respBody.Close()
	}

	return nil
}

func createMpu(c ObjectsClient, ctx context.Context, input *CreateMpuInput) (*CreateMpuOutput, error) {
	absPath := absFileInput(c.client.AccountName, input.Body.ObjectPath)

	// Because some clients will be treating Manta like S3, they will
	// include slashes in object names which we'll need to convert to
	// directories
	if input.ForceInsert {
		absDirName := _AbsCleanPath(path.Dir(string(absPath)))
		exists, _ := checkDirectoryTreeExists(c, ctx, absDirName)
		if !exists {
			err := createDirectory(c, ctx, absDirName)
			if err != nil {
				return nil, errors.Wrap(err, "unable to create directory for create mpu operation")
			}
		}
	}
	headers := &http.Header{}
	for key, value := range input.Body.Headers {
		headers.Set(key, value)
	}
	if input.DurabilityLevel != 0 {
		headers.Set("Durability-Level", strconv.FormatUint(input.DurabilityLevel, 10))
	}
	if input.ContentLength != 0 {
		headers.Set("Content-Length", strconv.FormatUint(input.ContentLength, 10))
	}
	if input.ContentMD5 != "" {
		headers.Set("Content-MD5", input.ContentMD5)
	}

	input.Body.ObjectPath = string(absPath)
	reqInput := client.RequestInput{
		Method:  http.MethodPost,
		Path:    "/" + c.client.AccountName + "/uploads",
		Headers: headers,
		Body:    input.Body,
	}
	respBody, _, err := c.client.ExecuteRequestStorage(ctx, reqInput)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create mpu")
	}
	if respBody != nil {
		defer respBody.Close()
	}

	response := &CreateMpuOutput{}
	decoder := json.NewDecoder(respBody)
	if err = decoder.Decode(&response); err != nil {
		return nil, errors.Wrap(err, "unable to decode create mpu response")
	}

	return response, nil
}

func getMpu(c ObjectsClient, ctx context.Context, input *GetMpuInput) (*GetMpuOutput, error) {
	headers := &http.Header{}

	reqInput := client.RequestInput{
		Method:  http.MethodGet,
		Path:    input.PartsDirectoryPath + "/state",
		Headers: headers,
	}
	respBody, _, err := c.client.ExecuteRequestStorage(ctx, reqInput)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get mpu")
	}

	response := &GetMpuOutput{}
	decoder := json.NewDecoder(respBody)
	if err = decoder.Decode(&response); err != nil {
		return nil, errors.Wrap(err, "unable to decode get mpu response")
	}

	return response, nil
}

func listMpuParts(c ObjectsClient, ctx context.Context, input *ListMpuPartsInput) (*ListMpuPartsOutput, error) {
	id := input.Id
	idLength := len(id)
	prefixLen, err := strconv.Atoi(id[idLength-1 : idLength])
	if err != nil {
		return nil, errors.Wrap(err, "unable to upload part")
	}
	prefix := id[:prefixLen]
	partPath := "/" + c.client.AccountName + "/uploads/" + prefix + "/" + input.Id + "/"
	listDirInput := ListDirectoryInput{
		DirectoryName: partPath,
	}

	dirClient := &DirectoryClient{
		client: c.client,
	}

	listDirOutput, err := dirClient.List(ctx, &listDirInput)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list mpu parts")
	}

	var parts []ListMpuPart
	for num, part := range listDirOutput.Entries {
		parts = append(parts, ListMpuPart{
			ETag:       part.ETag,
			PartNumber: num,
			Size:       int64(part.Size),
		})
	}

	listMpuPartsOutput := &ListMpuPartsOutput{
		Parts: parts,
	}

	return listMpuPartsOutput, nil
}

func uploadPart(c ObjectsClient, ctx context.Context, input *UploadPartInput) (*UploadPartOutput, error) {
	headers := &http.Header{}
	for key, value := range input.Headers {
		headers.Set(key, value)
	}

	if input.ContentMD5 != "" {
		headers.Set("Content-MD5", input.ContentMD5)
	}

	// The mpu directory prefix length is derived from the final character
	// in the mpu identifier which we'll call P. The mpu prefix itself is
	// the first P characters of the mpu identifier. In order to derive the
	// correct directory structure we need to parse this information from
	// the mpu identifier
	id := input.Id
	idLength := len(id)
	partNum := strconv.FormatUint(input.PartNum, 10)
	prefixLen, err := strconv.Atoi(id[idLength-1 : idLength])
	if err != nil {
		return nil, errors.Wrap(err, "unable to upload part due to invalid mpu prefix length")
	}
	prefix := id[:prefixLen]
	partPath := "/" + c.client.AccountName + "/uploads/" + prefix + "/" + input.Id + "/" + partNum

	reqInput := client.RequestNoEncodeInput{
		Method:  http.MethodPut,
		Path:    partPath,
		Headers: headers,
		Body:    input.ObjectReader,
	}
	respBody, respHeader, err := c.client.ExecuteRequestNoEncode(ctx, reqInput)
	if respBody != nil {
		defer respBody.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to upload part")
	}

	uploadPartOutput := &UploadPartOutput{
		Part: respHeader.Get("Etag"),
	}
	return uploadPartOutput, nil
}

func checkDirectoryTreeExists(c ObjectsClient, ctx context.Context, absPath _AbsCleanPath) (bool, error) {
	exists, err := c.IsDir(ctx, string(absPath))
	if err != nil {
		if tt.IsResourceNotFoundError(err) || tt.IsStatusNotFoundCode(err) {
			return false, nil
		}
		return false, err
	}
	if exists {
		return true, nil
	}

	return false, nil
}
