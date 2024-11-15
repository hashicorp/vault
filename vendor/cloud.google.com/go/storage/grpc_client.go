// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net/url"
	"os"

	"cloud.google.com/go/iam/apiv1/iampb"
	"cloud.google.com/go/internal/trace"
	gapic "cloud.google.com/go/storage/internal/apiv2"
	"cloud.google.com/go/storage/internal/apiv2/storagepb"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/mem"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
)

const (
	// defaultConnPoolSize is the default number of channels
	// to initialize in the GAPIC gRPC connection pool. A larger
	// connection pool may be necessary for jobs that require
	// high throughput and/or leverage many concurrent streams
	// if not running via DirectPath.
	//
	// This is only used for the gRPC client.
	defaultConnPoolSize = 1

	// maxPerMessageWriteSize is the maximum amount of content that can be sent
	// per WriteObjectRequest message. A buffer reaching this amount will
	// precipitate a flush of the buffer. It is only used by the gRPC Writer
	// implementation.
	maxPerMessageWriteSize int = int(storagepb.ServiceConstants_MAX_WRITE_CHUNK_BYTES)

	// globalProjectAlias is the project ID alias used for global buckets.
	//
	// This is only used for the gRPC API.
	globalProjectAlias = "_"

	// msgEntityNotSupported indicates ACL entites using project ID are not currently supported.
	//
	// This is only used for the gRPC API.
	msgEntityNotSupported = "The gRPC API currently does not support ACL entities using project ID, use project numbers instead"
)

// defaultGRPCOptions returns a set of the default client options
// for gRPC client initialization.
func defaultGRPCOptions() []option.ClientOption {
	defaults := []option.ClientOption{
		option.WithGRPCConnectionPool(defaultConnPoolSize),
	}

	// Set emulator options for gRPC if an emulator was specified. Note that in a
	// hybrid client, STORAGE_EMULATOR_HOST will set the host to use for HTTP and
	// STORAGE_EMULATOR_HOST_GRPC will set the host to use for gRPC (when using a
	// local emulator, HTTP and gRPC must use different ports, so this is
	// necessary).
	//
	// TODO: When the newHybridClient is not longer used, remove
	// STORAGE_EMULATOR_HOST_GRPC and use STORAGE_EMULATOR_HOST for both the
	// HTTP and gRPC based clients.
	if host := os.Getenv("STORAGE_EMULATOR_HOST_GRPC"); host != "" {
		// Strip the scheme from the emulator host. WithEndpoint does not take a
		// scheme for gRPC.
		host = stripScheme(host)

		defaults = append(defaults,
			option.WithEndpoint(host),
			option.WithGRPCDialOption(grpc.WithInsecure()),
			option.WithoutAuthentication(),
			WithDisabledClientMetrics(),
		)
	} else {
		// Only enable DirectPath when the emulator is not being targeted.
		defaults = append(defaults, internaloption.EnableDirectPath(true), internaloption.EnableDirectPathXds())
	}

	return defaults
}

// grpcStorageClient is the gRPC API implementation of the transport-agnostic
// storageClient interface.
type grpcStorageClient struct {
	raw      *gapic.Client
	settings *settings
}

// newGRPCStorageClient initializes a new storageClient that uses the gRPC
// Storage API.
func newGRPCStorageClient(ctx context.Context, opts ...storageOption) (storageClient, error) {
	s := initSettings(opts...)
	s.clientOption = append(defaultGRPCOptions(), s.clientOption...)
	// Disable all gax-level retries in favor of retry logic in the veneer client.
	s.gax = append(s.gax, gax.WithRetry(nil))

	config := newStorageConfig(s.clientOption...)
	if config.readAPIWasSet {
		return nil, errors.New("storage: GRPC is incompatible with any option that specifies an API for reads")
	}

	if !config.disableClientMetrics {
		// Do not fail client creation if enabling metrics fails.
		if metricsContext, err := enableClientMetrics(ctx, s, config); err == nil {
			s.metricsContext = metricsContext
			s.clientOption = append(s.clientOption, metricsContext.clientOpts...)
		} else {
			log.Printf("Failed to enable client metrics: %v", err)
		}
	}
	g, err := gapic.NewClient(ctx, s.clientOption...)
	if err != nil {
		return nil, err
	}

	return &grpcStorageClient{
		raw:      g,
		settings: s,
	}, nil
}

func (c *grpcStorageClient) Close() error {
	if c.settings.metricsContext != nil {
		c.settings.metricsContext.close()
	}
	return c.raw.Close()
}

// Top-level methods.

// GetServiceAccount is not supported in the gRPC client.
func (c *grpcStorageClient) GetServiceAccount(ctx context.Context, project string, opts ...storageOption) (string, error) {
	return "", errMethodNotSupported
}

func (c *grpcStorageClient) CreateBucket(ctx context.Context, project, bucket string, attrs *BucketAttrs, enableObjectRetention *bool, opts ...storageOption) (*BucketAttrs, error) {
	if enableObjectRetention != nil {
		// TO-DO: implement ObjectRetention once available - see b/308194853
		return nil, status.Errorf(codes.Unimplemented, "storage: object retention is not supported in gRPC")
	}

	s := callSettings(c.settings, opts...)
	b := attrs.toProtoBucket()
	b.Project = toProjectResource(project)
	// If there is lifecycle information but no location, explicitly set
	// the location. This is a GCS quirk/bug.
	if b.GetLocation() == "" && b.GetLifecycle() != nil {
		b.Location = "US"
	}

	req := &storagepb.CreateBucketRequest{
		Parent:   fmt.Sprintf("projects/%s", globalProjectAlias),
		Bucket:   b,
		BucketId: bucket,
	}
	if attrs != nil {
		req.PredefinedAcl = attrs.PredefinedACL
		req.PredefinedDefaultObjectAcl = attrs.PredefinedDefaultObjectACL
	}

	var battrs *BucketAttrs
	err := run(ctx, func(ctx context.Context) error {
		res, err := c.raw.CreateBucket(ctx, req, s.gax...)

		battrs = newBucketFromProto(res)

		return err
	}, s.retry, s.idempotent)

	return battrs, err
}

func (c *grpcStorageClient) ListBuckets(ctx context.Context, project string, opts ...storageOption) *BucketIterator {
	s := callSettings(c.settings, opts...)
	it := &BucketIterator{
		ctx:       ctx,
		projectID: project,
	}

	var gitr *gapic.BucketIterator
	fetch := func(pageSize int, pageToken string) (token string, err error) {

		var buckets []*storagepb.Bucket
		var next string
		err = run(it.ctx, func(ctx context.Context) error {
			// Initialize GAPIC-based iterator when pageToken is empty, which
			// indicates that this fetch call is attempting to get the first page.
			//
			// Note: Initializing the GAPIC-based iterator lazily is necessary to
			// capture the BucketIterator.Prefix set by the user *after* the
			// BucketIterator is returned to them from the veneer.
			if pageToken == "" {
				req := &storagepb.ListBucketsRequest{
					Parent: toProjectResource(it.projectID),
					Prefix: it.Prefix,
				}
				gitr = c.raw.ListBuckets(ctx, req, s.gax...)
			}
			buckets, next, err = gitr.InternalFetch(pageSize, pageToken)
			return err
		}, s.retry, s.idempotent)
		if err != nil {
			return "", err
		}

		for _, bkt := range buckets {
			b := newBucketFromProto(bkt)
			it.buckets = append(it.buckets, b)
		}

		return next, nil
	}
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(
		fetch,
		func() int { return len(it.buckets) },
		func() interface{} { b := it.buckets; it.buckets = nil; return b })

	return it
}

// Bucket methods.

func (c *grpcStorageClient) DeleteBucket(ctx context.Context, bucket string, conds *BucketConditions, opts ...storageOption) error {
	s := callSettings(c.settings, opts...)
	req := &storagepb.DeleteBucketRequest{
		Name: bucketResourceName(globalProjectAlias, bucket),
	}
	if err := applyBucketCondsProto("grpcStorageClient.DeleteBucket", conds, req); err != nil {
		return err
	}
	if s.userProject != "" {
		ctx = setUserProjectMetadata(ctx, s.userProject)
	}

	return run(ctx, func(ctx context.Context) error {
		return c.raw.DeleteBucket(ctx, req, s.gax...)
	}, s.retry, s.idempotent)
}

func (c *grpcStorageClient) GetBucket(ctx context.Context, bucket string, conds *BucketConditions, opts ...storageOption) (*BucketAttrs, error) {
	s := callSettings(c.settings, opts...)
	req := &storagepb.GetBucketRequest{
		Name:     bucketResourceName(globalProjectAlias, bucket),
		ReadMask: &fieldmaskpb.FieldMask{Paths: []string{"*"}},
	}
	if err := applyBucketCondsProto("grpcStorageClient.GetBucket", conds, req); err != nil {
		return nil, err
	}
	if s.userProject != "" {
		ctx = setUserProjectMetadata(ctx, s.userProject)
	}

	var battrs *BucketAttrs
	err := run(ctx, func(ctx context.Context) error {
		res, err := c.raw.GetBucket(ctx, req, s.gax...)

		battrs = newBucketFromProto(res)

		return err
	}, s.retry, s.idempotent)

	if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
		return nil, ErrBucketNotExist
	}

	return battrs, err
}
func (c *grpcStorageClient) UpdateBucket(ctx context.Context, bucket string, uattrs *BucketAttrsToUpdate, conds *BucketConditions, opts ...storageOption) (*BucketAttrs, error) {
	s := callSettings(c.settings, opts...)
	b := uattrs.toProtoBucket()
	b.Name = bucketResourceName(globalProjectAlias, bucket)
	req := &storagepb.UpdateBucketRequest{
		Bucket:                     b,
		PredefinedAcl:              uattrs.PredefinedACL,
		PredefinedDefaultObjectAcl: uattrs.PredefinedDefaultObjectACL,
	}
	if err := applyBucketCondsProto("grpcStorageClient.UpdateBucket", conds, req); err != nil {
		return nil, err
	}
	if s.userProject != "" {
		ctx = setUserProjectMetadata(ctx, s.userProject)
	}

	var paths []string
	fieldMask := &fieldmaskpb.FieldMask{
		Paths: paths,
	}
	if uattrs.CORS != nil {
		fieldMask.Paths = append(fieldMask.Paths, "cors")
	}
	if uattrs.DefaultEventBasedHold != nil {
		fieldMask.Paths = append(fieldMask.Paths, "default_event_based_hold")
	}
	if uattrs.RetentionPolicy != nil {
		fieldMask.Paths = append(fieldMask.Paths, "retention_policy")
	}
	if uattrs.VersioningEnabled != nil {
		fieldMask.Paths = append(fieldMask.Paths, "versioning")
	}
	if uattrs.RequesterPays != nil {
		fieldMask.Paths = append(fieldMask.Paths, "billing")
	}
	if uattrs.BucketPolicyOnly != nil || uattrs.UniformBucketLevelAccess != nil || uattrs.PublicAccessPrevention != PublicAccessPreventionUnknown {
		fieldMask.Paths = append(fieldMask.Paths, "iam_config")
	}
	if uattrs.Encryption != nil {
		fieldMask.Paths = append(fieldMask.Paths, "encryption")
	}
	if uattrs.Lifecycle != nil {
		fieldMask.Paths = append(fieldMask.Paths, "lifecycle")
	}
	if uattrs.Logging != nil {
		fieldMask.Paths = append(fieldMask.Paths, "logging")
	}
	if uattrs.Website != nil {
		fieldMask.Paths = append(fieldMask.Paths, "website")
	}
	if uattrs.PredefinedACL != "" {
		// In cases where PredefinedACL is set, Acl is cleared.
		fieldMask.Paths = append(fieldMask.Paths, "acl")
	}
	if uattrs.PredefinedDefaultObjectACL != "" {
		// In cases where PredefinedDefaultObjectACL is set, DefaultObjectAcl is cleared.
		fieldMask.Paths = append(fieldMask.Paths, "default_object_acl")
	}
	// Note: This API currently does not support entites using project ID.
	// Use project numbers in ACL entities. Pending b/233617896.
	if uattrs.acl != nil {
		// In cases where acl is set by UpdateBucketACL method.
		fieldMask.Paths = append(fieldMask.Paths, "acl")
	}
	if uattrs.defaultObjectACL != nil {
		// In cases where defaultObjectACL is set by UpdateBucketACL method.
		fieldMask.Paths = append(fieldMask.Paths, "default_object_acl")
	}
	if uattrs.StorageClass != "" {
		fieldMask.Paths = append(fieldMask.Paths, "storage_class")
	}
	if uattrs.RPO != RPOUnknown {
		fieldMask.Paths = append(fieldMask.Paths, "rpo")
	}
	if uattrs.Autoclass != nil {
		fieldMask.Paths = append(fieldMask.Paths, "autoclass")
	}
	if uattrs.SoftDeletePolicy != nil {
		fieldMask.Paths = append(fieldMask.Paths, "soft_delete_policy")
	}

	for label := range uattrs.setLabels {
		fieldMask.Paths = append(fieldMask.Paths, fmt.Sprintf("labels.%s", label))
	}

	// Delete a label by not including it in Bucket.Labels but adding the key to the update mask.
	for label := range uattrs.deleteLabels {
		fieldMask.Paths = append(fieldMask.Paths, fmt.Sprintf("labels.%s", label))
	}

	req.UpdateMask = fieldMask

	if len(fieldMask.Paths) < 1 {
		// Nothing to update. Send a get request for current attrs instead. This
		// maintains consistency with JSON bucket updates.
		opts = append(opts, idempotent(true))
		return c.GetBucket(ctx, bucket, conds, opts...)
	}

	var battrs *BucketAttrs
	err := run(ctx, func(ctx context.Context) error {
		res, err := c.raw.UpdateBucket(ctx, req, s.gax...)
		battrs = newBucketFromProto(res)
		return err
	}, s.retry, s.idempotent)

	return battrs, err
}
func (c *grpcStorageClient) LockBucketRetentionPolicy(ctx context.Context, bucket string, conds *BucketConditions, opts ...storageOption) error {
	s := callSettings(c.settings, opts...)
	req := &storagepb.LockBucketRetentionPolicyRequest{
		Bucket: bucketResourceName(globalProjectAlias, bucket),
	}
	if err := applyBucketCondsProto("grpcStorageClient.LockBucketRetentionPolicy", conds, req); err != nil {
		return err
	}

	return run(ctx, func(ctx context.Context) error {
		_, err := c.raw.LockBucketRetentionPolicy(ctx, req, s.gax...)
		return err
	}, s.retry, s.idempotent)

}
func (c *grpcStorageClient) ListObjects(ctx context.Context, bucket string, q *Query, opts ...storageOption) *ObjectIterator {
	s := callSettings(c.settings, opts...)
	it := &ObjectIterator{
		ctx: ctx,
	}
	if q != nil {
		it.query = *q
	}
	req := &storagepb.ListObjectsRequest{
		Parent:                   bucketResourceName(globalProjectAlias, bucket),
		Prefix:                   it.query.Prefix,
		Delimiter:                it.query.Delimiter,
		Versions:                 it.query.Versions,
		LexicographicStart:       it.query.StartOffset,
		LexicographicEnd:         it.query.EndOffset,
		IncludeTrailingDelimiter: it.query.IncludeTrailingDelimiter,
		MatchGlob:                it.query.MatchGlob,
		ReadMask:                 q.toFieldMask(), // a nil Query still results in a "*" FieldMask
		SoftDeleted:              it.query.SoftDeleted,
		IncludeFoldersAsPrefixes: it.query.IncludeFoldersAsPrefixes,
	}
	if s.userProject != "" {
		ctx = setUserProjectMetadata(ctx, s.userProject)
	}
	fetch := func(pageSize int, pageToken string) (token string, err error) {
		var objects []*storagepb.Object
		var gitr *gapic.ObjectIterator
		err = run(it.ctx, func(ctx context.Context) error {
			gitr = c.raw.ListObjects(ctx, req, s.gax...)
			it.ctx = ctx
			objects, token, err = gitr.InternalFetch(pageSize, pageToken)
			return err
		}, s.retry, s.idempotent)
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
				err = ErrBucketNotExist
			}
			return "", err
		}

		for _, obj := range objects {
			b := newObjectFromProto(obj)
			it.items = append(it.items, b)
		}

		// Response is always non-nil after a successful request.
		res := gitr.Response.(*storagepb.ListObjectsResponse)
		for _, prefix := range res.GetPrefixes() {
			it.items = append(it.items, &ObjectAttrs{Prefix: prefix})
		}

		return token, nil
	}
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(
		fetch,
		func() int { return len(it.items) },
		func() interface{} { b := it.items; it.items = nil; return b })

	return it
}

// Object metadata methods.

func (c *grpcStorageClient) DeleteObject(ctx context.Context, bucket, object string, gen int64, conds *Conditions, opts ...storageOption) error {
	s := callSettings(c.settings, opts...)
	req := &storagepb.DeleteObjectRequest{
		Bucket: bucketResourceName(globalProjectAlias, bucket),
		Object: object,
	}
	if err := applyCondsProto("grpcStorageClient.DeleteObject", gen, conds, req); err != nil {
		return err
	}
	if s.userProject != "" {
		ctx = setUserProjectMetadata(ctx, s.userProject)
	}
	err := run(ctx, func(ctx context.Context) error {
		return c.raw.DeleteObject(ctx, req, s.gax...)
	}, s.retry, s.idempotent)
	if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
		return ErrObjectNotExist
	}
	return err
}

func (c *grpcStorageClient) GetObject(ctx context.Context, params *getObjectParams, opts ...storageOption) (*ObjectAttrs, error) {
	s := callSettings(c.settings, opts...)
	req := &storagepb.GetObjectRequest{
		Bucket: bucketResourceName(globalProjectAlias, params.bucket),
		Object: params.object,
		// ProjectionFull by default.
		ReadMask: &fieldmaskpb.FieldMask{Paths: []string{"*"}},
	}
	if err := applyCondsProto("grpcStorageClient.GetObject", params.gen, params.conds, req); err != nil {
		return nil, err
	}
	if s.userProject != "" {
		ctx = setUserProjectMetadata(ctx, s.userProject)
	}
	if params.encryptionKey != nil {
		req.CommonObjectRequestParams = toProtoCommonObjectRequestParams(params.encryptionKey)
	}
	if params.softDeleted {
		req.SoftDeleted = &params.softDeleted
	}

	var attrs *ObjectAttrs
	err := run(ctx, func(ctx context.Context) error {
		res, err := c.raw.GetObject(ctx, req, s.gax...)
		attrs = newObjectFromProto(res)

		return err
	}, s.retry, s.idempotent)

	if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
		return nil, ErrObjectNotExist
	}

	return attrs, err
}

func (c *grpcStorageClient) UpdateObject(ctx context.Context, params *updateObjectParams, opts ...storageOption) (*ObjectAttrs, error) {
	uattrs := params.uattrs
	if params.overrideRetention != nil || uattrs.Retention != nil {
		// TO-DO: implement ObjectRetention once available - see b/308194853
		return nil, status.Errorf(codes.Unimplemented, "storage: object retention is not supported in gRPC")
	}
	s := callSettings(c.settings, opts...)
	o := uattrs.toProtoObject(bucketResourceName(globalProjectAlias, params.bucket), params.object)
	// For Update, generation is passed via the object message rather than a field on the request.
	if params.gen >= 0 {
		o.Generation = params.gen
	}
	req := &storagepb.UpdateObjectRequest{
		Object:        o,
		PredefinedAcl: uattrs.PredefinedACL,
	}
	if err := applyCondsProto("grpcStorageClient.UpdateObject", defaultGen, params.conds, req); err != nil {
		return nil, err
	}
	if s.userProject != "" {
		ctx = setUserProjectMetadata(ctx, s.userProject)
	}
	if params.encryptionKey != nil {
		req.CommonObjectRequestParams = toProtoCommonObjectRequestParams(params.encryptionKey)
	}

	fieldMask := &fieldmaskpb.FieldMask{Paths: nil}
	if uattrs.EventBasedHold != nil {
		fieldMask.Paths = append(fieldMask.Paths, "event_based_hold")
	}
	if uattrs.TemporaryHold != nil {
		fieldMask.Paths = append(fieldMask.Paths, "temporary_hold")
	}
	if uattrs.ContentType != nil {
		fieldMask.Paths = append(fieldMask.Paths, "content_type")
	}
	if uattrs.ContentLanguage != nil {
		fieldMask.Paths = append(fieldMask.Paths, "content_language")
	}
	if uattrs.ContentEncoding != nil {
		fieldMask.Paths = append(fieldMask.Paths, "content_encoding")
	}
	if uattrs.ContentDisposition != nil {
		fieldMask.Paths = append(fieldMask.Paths, "content_disposition")
	}
	if uattrs.CacheControl != nil {
		fieldMask.Paths = append(fieldMask.Paths, "cache_control")
	}
	if !uattrs.CustomTime.IsZero() {
		fieldMask.Paths = append(fieldMask.Paths, "custom_time")
	}
	// Note: This API currently does not support entites using project ID.
	// Use project numbers in ACL entities. Pending b/233617896.
	if uattrs.ACL != nil || len(uattrs.PredefinedACL) > 0 {
		fieldMask.Paths = append(fieldMask.Paths, "acl")
	}

	if uattrs.Metadata != nil {
		// We don't support deleting a specific metadata key; metadata is deleted
		// as a whole if provided an empty map, so we do not use dot notation here
		if len(uattrs.Metadata) == 0 {
			fieldMask.Paths = append(fieldMask.Paths, "metadata")
		} else {
			// We can, however, use dot notation for adding keys
			for key := range uattrs.Metadata {
				fieldMask.Paths = append(fieldMask.Paths, fmt.Sprintf("metadata.%s", key))
			}
		}
	}

	req.UpdateMask = fieldMask

	if len(fieldMask.Paths) < 1 {
		// Nothing to update. To maintain consistency with JSON, we must still
		// update the object because metageneration and other fields are
		// updated even on an empty update.
		// gRPC will fail if the fieldmask is empty, so instead we add an
		// output-only field to the update mask. Output-only fields are (and must
		// be - see AIP 161) ignored, but allow us to send an empty update because
		// any mask that is valid for read (as this one is) must be valid for write.
		fieldMask.Paths = append(fieldMask.Paths, "create_time")
	}

	var attrs *ObjectAttrs
	err := run(ctx, func(ctx context.Context) error {
		res, err := c.raw.UpdateObject(ctx, req, s.gax...)
		attrs = newObjectFromProto(res)
		return err
	}, s.retry, s.idempotent)
	if e, ok := status.FromError(err); ok && e.Code() == codes.NotFound {
		return nil, ErrObjectNotExist
	}

	return attrs, err
}

func (c *grpcStorageClient) RestoreObject(ctx context.Context, params *restoreObjectParams, opts ...storageOption) (*ObjectAttrs, error) {
	s := callSettings(c.settings, opts...)
	req := &storagepb.RestoreObjectRequest{
		Bucket:        bucketResourceName(globalProjectAlias, params.bucket),
		Object:        params.object,
		CopySourceAcl: &params.copySourceACL,
	}
	if err := applyCondsProto("grpcStorageClient.RestoreObject", params.gen, params.conds, req); err != nil {
		return nil, err
	}
	if s.userProject != "" {
		ctx = setUserProjectMetadata(ctx, s.userProject)
	}

	var attrs *ObjectAttrs
	err := run(ctx, func(ctx context.Context) error {
		res, err := c.raw.RestoreObject(ctx, req, s.gax...)
		attrs = newObjectFromProto(res)
		return err
	}, s.retry, s.idempotent)
	if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
		return nil, ErrObjectNotExist
	}
	return attrs, err
}

// Default Object ACL methods.

func (c *grpcStorageClient) DeleteDefaultObjectACL(ctx context.Context, bucket string, entity ACLEntity, opts ...storageOption) error {
	// There is no separate API for PATCH in gRPC.
	// Make a GET call first to retrieve BucketAttrs.
	attrs, err := c.GetBucket(ctx, bucket, nil, opts...)
	if err != nil {
		return err
	}
	// Delete the entity and copy other remaining ACL entities.
	// Note: This API currently does not support entites using project ID.
	// Use project numbers in ACL entities. Pending b/233617896.
	// Return error if entity is not found or a project ID is used.
	invalidEntity := true
	var acl []ACLRule
	for _, a := range attrs.DefaultObjectACL {
		if a.Entity != entity {
			acl = append(acl, a)
		}
		if a.Entity == entity {
			invalidEntity = false
		}
	}
	if invalidEntity {
		return fmt.Errorf("storage: entity %v was not found on bucket %v, got %v. %v", entity, bucket, attrs.DefaultObjectACL, msgEntityNotSupported)
	}
	uattrs := &BucketAttrsToUpdate{defaultObjectACL: acl}
	// Call UpdateBucket with a MetagenerationMatch precondition set.
	if _, err = c.UpdateBucket(ctx, bucket, uattrs, &BucketConditions{MetagenerationMatch: attrs.MetaGeneration}, opts...); err != nil {
		return err
	}
	return nil
}

func (c *grpcStorageClient) ListDefaultObjectACLs(ctx context.Context, bucket string, opts ...storageOption) ([]ACLRule, error) {
	attrs, err := c.GetBucket(ctx, bucket, nil, opts...)
	if err != nil {
		return nil, err
	}
	return attrs.DefaultObjectACL, nil
}

func (c *grpcStorageClient) UpdateDefaultObjectACL(ctx context.Context, bucket string, entity ACLEntity, role ACLRole, opts ...storageOption) error {
	// There is no separate API for PATCH in gRPC.
	// Make a GET call first to retrieve BucketAttrs.
	attrs, err := c.GetBucket(ctx, bucket, nil, opts...)
	if err != nil {
		return err
	}
	// Note: This API currently does not support entites using project ID.
	// Use project numbers in ACL entities. Pending b/233617896.
	var acl []ACLRule
	aclRule := ACLRule{Entity: entity, Role: role}
	acl = append(attrs.DefaultObjectACL, aclRule)
	uattrs := &BucketAttrsToUpdate{defaultObjectACL: acl}
	// Call UpdateBucket with a MetagenerationMatch precondition set.
	if _, err = c.UpdateBucket(ctx, bucket, uattrs, &BucketConditions{MetagenerationMatch: attrs.MetaGeneration}, opts...); err != nil {
		return err
	}
	return nil
}

// Bucket ACL methods.

func (c *grpcStorageClient) DeleteBucketACL(ctx context.Context, bucket string, entity ACLEntity, opts ...storageOption) error {
	// There is no separate API for PATCH in gRPC.
	// Make a GET call first to retrieve BucketAttrs.
	attrs, err := c.GetBucket(ctx, bucket, nil, opts...)
	if err != nil {
		return err
	}
	// Delete the entity and copy other remaining ACL entities.
	// Note: This API currently does not support entites using project ID.
	// Use project numbers in ACL entities. Pending b/233617896.
	// Return error if entity is not found or a project ID is used.
	invalidEntity := true
	var acl []ACLRule
	for _, a := range attrs.ACL {
		if a.Entity != entity {
			acl = append(acl, a)
		}
		if a.Entity == entity {
			invalidEntity = false
		}
	}
	if invalidEntity {
		return fmt.Errorf("storage: entity %v was not found on bucket %v, got %v. %v", entity, bucket, attrs.ACL, msgEntityNotSupported)
	}
	uattrs := &BucketAttrsToUpdate{acl: acl}
	// Call UpdateBucket with a MetagenerationMatch precondition set.
	if _, err = c.UpdateBucket(ctx, bucket, uattrs, &BucketConditions{MetagenerationMatch: attrs.MetaGeneration}, opts...); err != nil {
		return err
	}
	return nil
}

func (c *grpcStorageClient) ListBucketACLs(ctx context.Context, bucket string, opts ...storageOption) ([]ACLRule, error) {
	attrs, err := c.GetBucket(ctx, bucket, nil, opts...)
	if err != nil {
		return nil, err
	}
	return attrs.ACL, nil
}

func (c *grpcStorageClient) UpdateBucketACL(ctx context.Context, bucket string, entity ACLEntity, role ACLRole, opts ...storageOption) error {
	// There is no separate API for PATCH in gRPC.
	// Make a GET call first to retrieve BucketAttrs.
	attrs, err := c.GetBucket(ctx, bucket, nil, opts...)
	if err != nil {
		return err
	}
	// Note: This API currently does not support entites using project ID.
	// Use project numbers in ACL entities. Pending b/233617896.
	var acl []ACLRule
	aclRule := ACLRule{Entity: entity, Role: role}
	acl = append(attrs.ACL, aclRule)
	uattrs := &BucketAttrsToUpdate{acl: acl}
	// Call UpdateBucket with a MetagenerationMatch precondition set.
	if _, err = c.UpdateBucket(ctx, bucket, uattrs, &BucketConditions{MetagenerationMatch: attrs.MetaGeneration}, opts...); err != nil {
		return err
	}
	return nil
}

// Object ACL methods.

func (c *grpcStorageClient) DeleteObjectACL(ctx context.Context, bucket, object string, entity ACLEntity, opts ...storageOption) error {
	// There is no separate API for PATCH in gRPC.
	// Make a GET call first to retrieve ObjectAttrs.
	attrs, err := c.GetObject(ctx, &getObjectParams{bucket, object, defaultGen, nil, nil, false}, opts...)
	if err != nil {
		return err
	}
	// Delete the entity and copy other remaining ACL entities.
	// Note: This API currently does not support entites using project ID.
	// Use project numbers in ACL entities. Pending b/233617896.
	// Return error if entity is not found or a project ID is used.
	invalidEntity := true
	var acl []ACLRule
	for _, a := range attrs.ACL {
		if a.Entity != entity {
			acl = append(acl, a)
		}
		if a.Entity == entity {
			invalidEntity = false
		}
	}
	if invalidEntity {
		return fmt.Errorf("storage: entity %v was not found on bucket %v, got %v. %v", entity, bucket, attrs.ACL, msgEntityNotSupported)
	}
	uattrs := &ObjectAttrsToUpdate{ACL: acl}
	// Call UpdateObject with the specified metageneration.
	params := &updateObjectParams{bucket: bucket, object: object, uattrs: uattrs, gen: defaultGen, conds: &Conditions{MetagenerationMatch: attrs.Metageneration}}
	if _, err = c.UpdateObject(ctx, params, opts...); err != nil {
		return err
	}
	return nil
}

// ListObjectACLs retrieves object ACL entries. By default, it operates on the latest generation of this object.
// Selecting a specific generation of this object is not currently supported by the client.
func (c *grpcStorageClient) ListObjectACLs(ctx context.Context, bucket, object string, opts ...storageOption) ([]ACLRule, error) {
	o, err := c.GetObject(ctx, &getObjectParams{bucket, object, defaultGen, nil, nil, false}, opts...)
	if err != nil {
		return nil, err
	}
	return o.ACL, nil
}

func (c *grpcStorageClient) UpdateObjectACL(ctx context.Context, bucket, object string, entity ACLEntity, role ACLRole, opts ...storageOption) error {
	// There is no separate API for PATCH in gRPC.
	// Make a GET call first to retrieve ObjectAttrs.
	attrs, err := c.GetObject(ctx, &getObjectParams{bucket, object, defaultGen, nil, nil, false}, opts...)
	if err != nil {
		return err
	}
	// Note: This API currently does not support entites using project ID.
	// Use project numbers in ACL entities. Pending b/233617896.
	var acl []ACLRule
	aclRule := ACLRule{Entity: entity, Role: role}
	acl = append(attrs.ACL, aclRule)
	uattrs := &ObjectAttrsToUpdate{ACL: acl}
	// Call UpdateObject with the specified metageneration.
	params := &updateObjectParams{bucket: bucket, object: object, uattrs: uattrs, gen: defaultGen, conds: &Conditions{MetagenerationMatch: attrs.Metageneration}}
	if _, err = c.UpdateObject(ctx, params, opts...); err != nil {
		return err
	}
	return nil
}

// Media operations.

func (c *grpcStorageClient) ComposeObject(ctx context.Context, req *composeObjectRequest, opts ...storageOption) (*ObjectAttrs, error) {
	s := callSettings(c.settings, opts...)
	if s.userProject != "" {
		ctx = setUserProjectMetadata(ctx, s.userProject)
	}

	dstObjPb := req.dstObject.attrs.toProtoObject(req.dstBucket)
	dstObjPb.Name = req.dstObject.name

	if req.sendCRC32C {
		dstObjPb.Checksums.Crc32C = &req.dstObject.attrs.CRC32C
	}

	srcs := []*storagepb.ComposeObjectRequest_SourceObject{}
	for _, src := range req.srcs {
		srcObjPb := &storagepb.ComposeObjectRequest_SourceObject{Name: src.name, ObjectPreconditions: &storagepb.ComposeObjectRequest_SourceObject_ObjectPreconditions{}}
		if src.gen >= 0 {
			srcObjPb.Generation = src.gen
		}
		if err := applyCondsProto("ComposeObject source", defaultGen, src.conds, srcObjPb.ObjectPreconditions); err != nil {
			return nil, err
		}
		srcs = append(srcs, srcObjPb)
	}

	rawReq := &storagepb.ComposeObjectRequest{
		Destination:   dstObjPb,
		SourceObjects: srcs,
	}
	if err := applyCondsProto("ComposeObject destination", defaultGen, req.dstObject.conds, rawReq); err != nil {
		return nil, err
	}
	if req.predefinedACL != "" {
		rawReq.DestinationPredefinedAcl = req.predefinedACL
	}
	if req.dstObject.encryptionKey != nil {
		rawReq.CommonObjectRequestParams = toProtoCommonObjectRequestParams(req.dstObject.encryptionKey)
	}

	var obj *storagepb.Object
	var err error
	if err := run(ctx, func(ctx context.Context) error {
		obj, err = c.raw.ComposeObject(ctx, rawReq, s.gax...)
		return err
	}, s.retry, s.idempotent); err != nil {
		return nil, err
	}

	return newObjectFromProto(obj), nil
}
func (c *grpcStorageClient) RewriteObject(ctx context.Context, req *rewriteObjectRequest, opts ...storageOption) (*rewriteObjectResponse, error) {
	s := callSettings(c.settings, opts...)
	obj := req.dstObject.attrs.toProtoObject("")
	call := &storagepb.RewriteObjectRequest{
		SourceBucket:              bucketResourceName(globalProjectAlias, req.srcObject.bucket),
		SourceObject:              req.srcObject.name,
		RewriteToken:              req.token,
		DestinationBucket:         bucketResourceName(globalProjectAlias, req.dstObject.bucket),
		DestinationName:           req.dstObject.name,
		Destination:               obj,
		DestinationKmsKey:         req.dstObject.keyName,
		DestinationPredefinedAcl:  req.predefinedACL,
		CommonObjectRequestParams: toProtoCommonObjectRequestParams(req.dstObject.encryptionKey),
	}

	// The userProject, whether source or destination project, is decided by the code calling the interface.
	if s.userProject != "" {
		ctx = setUserProjectMetadata(ctx, s.userProject)
	}
	if err := applyCondsProto("Copy destination", defaultGen, req.dstObject.conds, call); err != nil {
		return nil, err
	}
	if err := applySourceCondsProto(req.srcObject.gen, req.srcObject.conds, call); err != nil {
		return nil, err
	}

	if len(req.dstObject.encryptionKey) > 0 {
		call.CommonObjectRequestParams = toProtoCommonObjectRequestParams(req.dstObject.encryptionKey)
	}
	if len(req.srcObject.encryptionKey) > 0 {
		srcParams := toProtoCommonObjectRequestParams(req.srcObject.encryptionKey)
		call.CopySourceEncryptionAlgorithm = srcParams.GetEncryptionAlgorithm()
		call.CopySourceEncryptionKeyBytes = srcParams.GetEncryptionKeyBytes()
		call.CopySourceEncryptionKeySha256Bytes = srcParams.GetEncryptionKeySha256Bytes()
	}

	call.MaxBytesRewrittenPerCall = req.maxBytesRewrittenPerCall

	var res *storagepb.RewriteResponse
	var err error

	retryCall := func(ctx context.Context) error { res, err = c.raw.RewriteObject(ctx, call, s.gax...); return err }

	if err := run(ctx, retryCall, s.retry, s.idempotent); err != nil {
		return nil, err
	}

	r := &rewriteObjectResponse{
		done:     res.GetDone(),
		written:  res.GetTotalBytesRewritten(),
		size:     res.GetObjectSize(),
		token:    res.GetRewriteToken(),
		resource: newObjectFromProto(res.GetResource()),
	}

	return r, nil
}

// Custom codec to be used for unmarshaling ReadObjectResponse messages.
// This is used to avoid a copy of object data in proto.Unmarshal.
type bytesCodecV2 struct {
}

var _ encoding.CodecV2 = bytesCodecV2{}

// Marshal is used to encode messages to send for bytesCodecV2. Since we are only
// using this to send ReadObjectRequest messages we don't need to recycle buffers
// here.
func (bytesCodecV2) Marshal(v any) (mem.BufferSlice, error) {
	vv, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message", v)
	}
	var data mem.BufferSlice
	buf, err := proto.Marshal(vv)
	if err != nil {
		return nil, err
	}
	data = append(data, mem.SliceBuffer(buf))
	return data, nil
}

// Unmarshal is used for data received for ReadObjectResponse. We want to preserve
// the mem.BufferSlice in most cases rather than copying and calling proto.Unmarshal.
func (bytesCodecV2) Unmarshal(data mem.BufferSlice, v any) error {
	switch v := v.(type) {
	case *mem.BufferSlice:
		*v = data
		// Pick up a reference to the data so that it is not freed while decoding.
		data.Ref()
		return nil
	case proto.Message:
		buf := data.MaterializeToBuffer(mem.DefaultBufferPool())
		return proto.Unmarshal(buf.ReadOnlyData(), v)
	default:
		return fmt.Errorf("cannot unmarshal type %T, want proto.Message or mem.BufferSlice", v)
	}
}

func (bytesCodecV2) Name() string {
	return ""
}

func (c *grpcStorageClient) NewRangeReader(ctx context.Context, params *newRangeReaderParams, opts ...storageOption) (r *Reader, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/storage.grpcStorageClient.NewRangeReader")
	defer func() { trace.EndSpan(ctx, err) }()

	s := callSettings(c.settings, opts...)

	s.gax = append(s.gax, gax.WithGRPCOptions(
		grpc.ForceCodecV2(bytesCodecV2{}),
	))

	if s.userProject != "" {
		ctx = setUserProjectMetadata(ctx, s.userProject)
	}

	b := bucketResourceName(globalProjectAlias, params.bucket)
	req := &storagepb.ReadObjectRequest{
		Bucket:                    b,
		Object:                    params.object,
		CommonObjectRequestParams: toProtoCommonObjectRequestParams(params.encryptionKey),
	}
	// The default is a negative value, which means latest.
	if params.gen >= 0 {
		req.Generation = params.gen
	}

	// Define a function that initiates a Read with offset and length, assuming
	// we have already read seen bytes.
	reopen := func(seen int64) (*readStreamResponse, context.CancelFunc, error) {
		// If the context has already expired, return immediately without making
		// we call.
		if err := ctx.Err(); err != nil {
			return nil, nil, err
		}

		cc, cancel := context.WithCancel(ctx)

		req.ReadOffset = params.offset + seen

		// Only set a ReadLimit if length is greater than zero, because <= 0 means
		// to read it all.
		if params.length > 0 {
			req.ReadLimit = params.length - seen
		}

		if err := applyCondsProto("gRPCReader.reopen", params.gen, params.conds, req); err != nil {
			cancel()
			return nil, nil, err
		}

		var stream storagepb.Storage_ReadObjectClient
		var err error
		var decoder *readResponseDecoder

		err = run(cc, func(ctx context.Context) error {
			stream, err = c.raw.ReadObject(ctx, req, s.gax...)
			if err != nil {
				return err
			}

			// Receive the message into databuf as a wire-encoded message so we can
			// use a custom decoder to avoid an extra copy at the protobuf layer.
			databufs := mem.BufferSlice{}
			err := stream.RecvMsg(&databufs)
			// These types of errors show up on the Recv call, rather than the
			// initialization of the stream via ReadObject above.
			if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
				return ErrObjectNotExist
			}
			if err != nil {
				return err
			}
			// Use a custom decoder that uses protobuf unmarshalling for all
			// fields except the object data. Object data is handled separately
			// to avoid a copy.
			decoder = &readResponseDecoder{
				databufs: databufs,
			}
			err = decoder.readFullObjectResponse()
			return err
		}, s.retry, s.idempotent)
		if err != nil {
			// Close the stream context we just created to ensure we don't leak
			// resources.
			cancel()
			// Free any buffers.
			if decoder != nil && decoder.databufs != nil {
				decoder.databufs.Free()
			}
			return nil, nil, err
		}

		return &readStreamResponse{stream, decoder}, cancel, nil
	}

	res, cancel, err := reopen(0)
	if err != nil {
		return nil, err
	}

	// The first message was Recv'd on stream open, use it to populate the
	// object metadata.
	msg := res.decoder.msg
	obj := msg.GetMetadata()
	// This is the size of the entire object, even if only a range was requested.
	size := obj.GetSize()

	// Only support checksums when reading an entire object, not a range.
	var (
		wantCRC  uint32
		checkCRC bool
	)
	if checksums := msg.GetObjectChecksums(); checksums != nil && checksums.Crc32C != nil {
		if params.offset == 0 && params.length < 0 {
			checkCRC = true
		}
		wantCRC = checksums.GetCrc32C()
	}

	r = &Reader{
		Attrs: ReaderObjectAttrs{
			Size:            size,
			ContentType:     obj.GetContentType(),
			ContentEncoding: obj.GetContentEncoding(),
			CacheControl:    obj.GetCacheControl(),
			LastModified:    obj.GetUpdateTime().AsTime(),
			Metageneration:  obj.GetMetageneration(),
			Generation:      obj.GetGeneration(),
			CRC32C:          wantCRC,
		},
		reader: &gRPCReader{
			stream: res.stream,
			reopen: reopen,
			cancel: cancel,
			size:   size,
			// Preserve the decoder to read out object data when Read/WriteTo is called.
			currMsg:   res.decoder,
			settings:  s,
			zeroRange: params.length == 0,
			wantCRC:   wantCRC,
			checkCRC:  checkCRC,
		},
		checkCRC: checkCRC,
	}

	cr := msg.GetContentRange()
	if cr != nil {
		r.Attrs.StartOffset = cr.GetStart()
		r.remain = cr.GetEnd() - cr.GetStart()
	} else {
		r.remain = size
	}

	// For a zero-length request, explicitly close the stream and set remaining
	// bytes to zero.
	if params.length == 0 {
		r.remain = 0
		r.reader.Close()
	}

	return r, nil
}

func (c *grpcStorageClient) OpenWriter(params *openWriterParams, opts ...storageOption) (*io.PipeWriter, error) {
	s := callSettings(c.settings, opts...)

	var offset int64
	errorf := params.setError
	progress := params.progress
	setObj := params.setObj

	pr, pw := io.Pipe()
	gw := newGRPCWriter(c, params, pr)
	gw.settings = s
	if s.userProject != "" {
		gw.ctx = setUserProjectMetadata(gw.ctx, s.userProject)
	}

	// This function reads the data sent to the pipe and sends sets of messages
	// on the gRPC client-stream as the buffer is filled.
	go func() {
		defer close(params.donec)

		// Loop until there is an error or the Object has been finalized.
		for {
			// Note: This blocks until either the buffer is full or EOF is read.
			recvd, doneReading, err := gw.read()
			if err != nil {
				err = checkCanceled(err)
				errorf(err)
				pr.CloseWithError(err)
				return
			}

			if params.attrs.Retention != nil {
				// TO-DO: remove once ObjectRetention is available - see b/308194853
				err = status.Errorf(codes.Unimplemented, "storage: object retention is not supported in gRPC")
				errorf(err)
				pr.CloseWithError(err)
				return
			}
			// The chunk buffer is full, but there is no end in sight. This
			// means that either:
			// 1. A resumable upload will need to be used to send
			// multiple chunks, until we are done reading data. Start a
			// resumable upload if it has not already been started.
			// 2. ChunkSize of zero may also have a full buffer, but a resumable
			// session should not be initiated in this case.
			if !doneReading && gw.upid == "" && params.chunkSize != 0 {
				err = gw.startResumableUpload()
				if err != nil {
					err = checkCanceled(err)
					errorf(err)
					pr.CloseWithError(err)
					return
				}
			}

			o, off, err := gw.uploadBuffer(recvd, offset, doneReading)
			if err != nil {
				err = checkCanceled(err)
				errorf(err)
				pr.CloseWithError(err)
				return
			}

			// At this point, the current buffer has been uploaded. For resumable
			// uploads and chunkSize = 0, capture the committed offset here in case
			// the upload was not finalized and another chunk is to be uploaded. Call
			// the progress function for resumable uploads only.
			if gw.upid != "" || gw.chunkSize == 0 {
				offset = off
			}
			if gw.upid != "" {
				progress(offset)
			}

			// When we are done reading data without errors, set the object and
			// finish.
			if doneReading {
				// Build Object from server's response.
				setObj(newObjectFromProto(o))
				return
			}
		}
	}()

	return pw, nil
}

// IAM methods.

func (c *grpcStorageClient) GetIamPolicy(ctx context.Context, resource string, version int32, opts ...storageOption) (*iampb.Policy, error) {
	// TODO: Need a way to set UserProject, potentially in X-Goog-User-Project system parameter.
	s := callSettings(c.settings, opts...)
	req := &iampb.GetIamPolicyRequest{
		Resource: bucketResourceName(globalProjectAlias, resource),
		Options: &iampb.GetPolicyOptions{
			RequestedPolicyVersion: version,
		},
	}
	var rp *iampb.Policy
	err := run(ctx, func(ctx context.Context) error {
		var err error
		rp, err = c.raw.GetIamPolicy(ctx, req, s.gax...)
		return err
	}, s.retry, s.idempotent)

	return rp, err
}

func (c *grpcStorageClient) SetIamPolicy(ctx context.Context, resource string, policy *iampb.Policy, opts ...storageOption) error {
	// TODO: Need a way to set UserProject, potentially in X-Goog-User-Project system parameter.
	s := callSettings(c.settings, opts...)

	req := &iampb.SetIamPolicyRequest{
		Resource: bucketResourceName(globalProjectAlias, resource),
		Policy:   policy,
	}

	return run(ctx, func(ctx context.Context) error {
		_, err := c.raw.SetIamPolicy(ctx, req, s.gax...)
		return err
	}, s.retry, s.idempotent)
}

func (c *grpcStorageClient) TestIamPermissions(ctx context.Context, resource string, permissions []string, opts ...storageOption) ([]string, error) {
	// TODO: Need a way to set UserProject, potentially in X-Goog-User-Project system parameter.
	s := callSettings(c.settings, opts...)
	req := &iampb.TestIamPermissionsRequest{
		Resource:    bucketResourceName(globalProjectAlias, resource),
		Permissions: permissions,
	}
	var res *iampb.TestIamPermissionsResponse
	err := run(ctx, func(ctx context.Context) error {
		var err error
		res, err = c.raw.TestIamPermissions(ctx, req, s.gax...)
		return err
	}, s.retry, s.idempotent)
	if err != nil {
		return nil, err
	}
	return res.Permissions, nil
}

// HMAC Key methods are not implemented in gRPC client.

func (c *grpcStorageClient) GetHMACKey(ctx context.Context, project, accessID string, opts ...storageOption) (*HMACKey, error) {
	return nil, errMethodNotSupported
}

func (c *grpcStorageClient) ListHMACKeys(ctx context.Context, project, serviceAccountEmail string, showDeletedKeys bool, opts ...storageOption) *HMACKeysIterator {
	it := &HMACKeysIterator{
		ctx:       ctx,
		projectID: "",
		retry:     nil,
	}
	fetch := func(_ int, _ string) (token string, err error) {
		return "", errMethodNotSupported
	}
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(
		fetch,
		func() int { return 0 },
		func() interface{} { return nil },
	)
	return it
}

func (c *grpcStorageClient) UpdateHMACKey(ctx context.Context, project, serviceAccountEmail, accessID string, attrs *HMACKeyAttrsToUpdate, opts ...storageOption) (*HMACKey, error) {
	return nil, errMethodNotSupported
}

func (c *grpcStorageClient) CreateHMACKey(ctx context.Context, project, serviceAccountEmail string, opts ...storageOption) (*HMACKey, error) {
	return nil, errMethodNotSupported
}

func (c *grpcStorageClient) DeleteHMACKey(ctx context.Context, project string, accessID string, opts ...storageOption) error {
	return errMethodNotSupported
}

// Notification methods are not implemented in gRPC client.

func (c *grpcStorageClient) ListNotifications(ctx context.Context, bucket string, opts ...storageOption) (n map[string]*Notification, err error) {
	return nil, errMethodNotSupported
}

func (c *grpcStorageClient) CreateNotification(ctx context.Context, bucket string, n *Notification, opts ...storageOption) (ret *Notification, err error) {
	return nil, errMethodNotSupported
}

func (c *grpcStorageClient) DeleteNotification(ctx context.Context, bucket string, id string, opts ...storageOption) (err error) {
	return errMethodNotSupported
}

// setUserProjectMetadata appends a project ID to the outgoing Context metadata
// via the x-goog-user-project system parameter defined at
// https://cloud.google.com/apis/docs/system-parameters. This is only for
// billing purposes, and is generally optional, except for requester-pays
// buckets.
func setUserProjectMetadata(ctx context.Context, project string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "x-goog-user-project", project)
}

type readStreamResponse struct {
	stream  storagepb.Storage_ReadObjectClient
	decoder *readResponseDecoder
}

type gRPCReader struct {
	seen, size int64
	zeroRange  bool
	stream     storagepb.Storage_ReadObjectClient
	reopen     func(seen int64) (*readStreamResponse, context.CancelFunc, error)
	leftovers  []byte
	currMsg    *readResponseDecoder // decoder for the current message
	cancel     context.CancelFunc
	settings   *settings
	checkCRC   bool   // should we check the CRC?
	wantCRC    uint32 // the CRC32c value the server sent in the header
	gotCRC     uint32 // running crc
}

// Update the running CRC with the data in the slice, if CRC checking was enabled.
func (r *gRPCReader) updateCRC(b []byte) {
	if r.checkCRC {
		r.gotCRC = crc32.Update(r.gotCRC, crc32cTable, b)
	}
}

// Checks whether the CRC matches at the conclusion of a read, if CRC checking was enabled.
func (r *gRPCReader) runCRCCheck() error {
	if r.checkCRC && r.gotCRC != r.wantCRC {
		return fmt.Errorf("storage: bad CRC on read: got %d, want %d", r.gotCRC, r.wantCRC)
	}
	return nil
}

// Read reads bytes into the user's buffer from an open gRPC stream.
func (r *gRPCReader) Read(p []byte) (int, error) {
	// The entire object has been read by this reader, check the checksum if
	// necessary and return EOF.
	if r.size == r.seen || r.zeroRange {
		if err := r.runCRCCheck(); err != nil {
			return 0, err
		}
		return 0, io.EOF
	}

	// No stream to read from, either never initialized or Close was called.
	// Note: There is a potential concurrency issue if multiple routines are
	// using the same reader. One encounters an error and the stream is closed
	// and then reopened while the other routine attempts to read from it.
	if r.stream == nil {
		return 0, fmt.Errorf("storage: reader has been closed")
	}

	var n int

	// If there is data remaining in the current message, return what was
	// available to conform to the Reader
	// interface: https://pkg.go.dev/io#Reader.
	if !r.currMsg.done {
		n = r.currMsg.readAndUpdateCRC(p, func(b []byte) {
			r.updateCRC(b)
		})
		r.seen += int64(n)
		return n, nil
	}

	// Attempt to Recv the next message on the stream.
	// This will update r.currMsg with the decoder for the new message.
	err := r.recv()
	if err != nil {
		return 0, err
	}

	// TODO: Determine if we need to capture incremental CRC32C for this
	// chunk. The Object CRC32C checksum is captured when directed to read
	// the entire Object. If directed to read a range, we may need to
	// calculate the range's checksum for verification if the checksum is
	// present in the response here.
	// TODO: Figure out if we need to support decompressive transcoding
	// https://cloud.google.com/storage/docs/transcoding.

	n = r.currMsg.readAndUpdateCRC(p, func(b []byte) {
		r.updateCRC(b)
	})
	r.seen += int64(n)
	return n, nil
}

// WriteTo writes all the data requested by the Reader into w, implementing
// io.WriterTo.
func (r *gRPCReader) WriteTo(w io.Writer) (int64, error) {
	// The entire object has been read by this reader, check the checksum if
	// necessary and return nil.
	if r.size == r.seen || r.zeroRange {
		if err := r.runCRCCheck(); err != nil {
			return 0, err
		}
		return 0, nil
	}

	// No stream to read from, either never initialized or Close was called.
	// Note: There is a potential concurrency issue if multiple routines are
	// using the same reader. One encounters an error and the stream is closed
	// and then reopened while the other routine attempts to read from it.
	if r.stream == nil {
		return 0, fmt.Errorf("storage: reader has been closed")
	}

	// Track bytes written during before call.
	var alreadySeen = r.seen

	// Write any already received message to the stream. There will be some leftovers from the
	// original NewRangeReader call.
	if r.currMsg != nil && !r.currMsg.done {
		written, err := r.currMsg.writeToAndUpdateCRC(w, func(b []byte) {
			r.updateCRC(b)
		})
		r.seen += int64(written)
		r.currMsg = nil
		if err != nil {
			return r.seen - alreadySeen, err
		}
	}

	// Loop and receive additional messages until the entire data is written.
	for {
		// Attempt to receive the next message on the stream.
		// Will terminate with io.EOF once data has all come through.
		// recv() handles stream reopening and retry logic so no need for retries here.
		err := r.recv()
		if err != nil {
			if err == io.EOF {
				// We are done; check the checksum if necessary and return.
				err = r.runCRCCheck()
			}
			return r.seen - alreadySeen, err
		}

		// TODO: Determine if we need to capture incremental CRC32C for this
		// chunk. The Object CRC32C checksum is captured when directed to read
		// the entire Object. If directed to read a range, we may need to
		// calculate the range's checksum for verification if the checksum is
		// present in the response here.
		// TODO: Figure out if we need to support decompressive transcoding
		// https://cloud.google.com/storage/docs/transcoding.
		written, err := r.currMsg.writeToAndUpdateCRC(w, func(b []byte) {
			r.updateCRC(b)
		})
		r.seen += int64(written)
		if err != nil {
			return r.seen - alreadySeen, err
		}
	}

}

// Close cancels the read stream's context in order for it to be closed and
// collected, and frees any currently in use buffers.
func (r *gRPCReader) Close() error {
	if r.cancel != nil {
		r.cancel()
	}
	r.stream = nil
	r.currMsg = nil
	return nil
}

// recv attempts to Recv the next message on the stream and extract the object
// data that it contains. In the event that a retryable error is encountered,
// the stream will be closed, reopened, and RecvMsg again.
// This will attempt to Recv until one of the following is true:
//
// * Recv is successful
// * A non-retryable error is encountered
// * The Reader's context is canceled
//
// The last error received is the one that is returned, which could be from
// an attempt to reopen the stream.
func (r *gRPCReader) recv() error {
	databufs := mem.BufferSlice{}
	err := r.stream.RecvMsg(&databufs)

	var shouldRetry = ShouldRetry
	if r.settings.retry != nil && r.settings.retry.shouldRetry != nil {
		shouldRetry = r.settings.retry.shouldRetry
	}
	if err != nil && shouldRetry(err) {
		// This will "close" the existing stream and immediately attempt to
		// reopen the stream, but will backoff if further attempts are necessary.
		// Reopening the stream Recvs the first message, so if retrying is
		// successful, r.currMsg will be updated to include the new data.
		return r.reopenStream()
	}

	if err != nil {
		return err
	}

	r.currMsg = &readResponseDecoder{databufs: databufs}
	return r.currMsg.readFullObjectResponse()
}

// ReadObjectResponse field and subfield numbers.
const (
	checksummedDataField        = protowire.Number(1)
	checksummedDataContentField = protowire.Number(1)
	checksummedDataCRC32CField  = protowire.Number(2)
	objectChecksumsField        = protowire.Number(2)
	contentRangeField           = protowire.Number(3)
	metadataField               = protowire.Number(4)
)

// readResponseDecoder is a wrapper on the raw message, used to decode one message
// without copying object data. It also has methods to write out the resulting object
// data to the user application.
type readResponseDecoder struct {
	databufs mem.BufferSlice // raw bytes of the message being processed
	// Decoding offsets
	off     uint64 // offset in the messsage relative to the data as a whole
	currBuf int    // index of the current buffer being processed
	currOff uint64 // offset in the current buffer
	// Processed data
	msg         *storagepb.ReadObjectResponse // processed response message with all fields other than object data populated
	dataOffsets bufferSliceOffsets            // offsets of the object data in the message.
	done        bool                          // true if the data has been completely read.
}

type bufferSliceOffsets struct {
	startBuf, endBuf int    // indices of start and end buffers of object data in the msg
	startOff, endOff uint64 // offsets within these buffers where the data starts and ends.
	currBuf          int    // index of current buffer being read out to the user application.
	currOff          uint64 // offset of read in current buffer.
}

// peek ahead 10 bytes from the current offset in the databufs. This will return a
// slice of the current buffer if the bytes are all in one buffer, but will copy
// the bytes into a new buffer if the distance is split across buffers. Use this
// to allow protowire methods to be used to parse tags & fixed values.
// The max length of a varint tag is 10 bytes, see
// https://protobuf.dev/programming-guides/encoding/#varints . Other int types
// are shorter.
func (d *readResponseDecoder) peek() []byte {
	b := d.databufs[d.currBuf].ReadOnlyData()
	// Check if the tag will fit in the current buffer. If not, copy the next 10
	// bytes into a new buffer to ensure that we can read the tag correctly
	// without it being divided between buffers.
	tagBuf := b[d.currOff:]
	remainingInBuf := len(tagBuf)
	// If we have less than 10 bytes remaining and are not in the final buffer,
	// copy up to 10 bytes ahead from the next buffer.
	if remainingInBuf < binary.MaxVarintLen64 && d.currBuf != len(d.databufs)-1 {
		tagBuf = d.copyNextBytes(10)
	}
	return tagBuf
}

// Copies up to next n bytes into a new buffer, or fewer if fewer bytes remain in the
// buffers overall. Does not advance offsets.
func (d *readResponseDecoder) copyNextBytes(n int) []byte {
	remaining := n
	if r := d.databufs.Len() - int(d.off); r < remaining {
		remaining = r
	}
	currBuf := d.currBuf
	currOff := d.currOff
	var buf []byte
	for remaining > 0 {
		b := d.databufs[currBuf].ReadOnlyData()
		remainingInCurr := len(b[currOff:])
		if remainingInCurr < remaining {
			buf = append(buf, b[currOff:]...)
			remaining -= remainingInCurr
			currBuf++
			currOff = 0
		} else {
			buf = append(buf, b[currOff:currOff+uint64(remaining)]...)
			remaining = 0
		}
	}
	return buf
}

// Advance current buffer & byte offset in the decoding by n bytes. Returns an error if we
// go past the end of the data.
func (d *readResponseDecoder) advanceOffset(n uint64) error {
	remaining := n
	for remaining > 0 {
		remainingInCurr := uint64(d.databufs[d.currBuf].Len()) - d.currOff
		if remainingInCurr <= remaining {
			remaining -= remainingInCurr
			d.currBuf++
			d.currOff = 0
		} else {
			d.currOff += remaining
			remaining = 0
		}
	}
	// If we have advanced past the end of the buffers, something went wrong.
	if (d.currBuf == len(d.databufs) && d.currOff > 0) || d.currBuf > len(d.databufs) {
		return errors.New("decoding: truncated message, cannot advance offset")
	}
	d.off += n
	return nil

}

// This copies object data from the message into the buffer and returns the number of
// bytes copied. The data offsets are incremented in the message. The updateCRC
// function is called on the copied bytes.
func (d *readResponseDecoder) readAndUpdateCRC(p []byte, updateCRC func([]byte)) int {
	// For a completely empty message, just return 0
	if len(d.databufs) == 0 {
		return 0
	}
	databuf := d.databufs[d.dataOffsets.currBuf]
	startOff := d.dataOffsets.currOff
	var b []byte
	if d.dataOffsets.currBuf == d.dataOffsets.endBuf {
		b = databuf.ReadOnlyData()[startOff:d.dataOffsets.endOff]
	} else {
		b = databuf.ReadOnlyData()[startOff:]
	}
	n := copy(p, b)
	updateCRC(b[:n])
	d.dataOffsets.currOff += uint64(n)

	// We've read all the data from this message. Free the underlying buffers.
	if d.dataOffsets.currBuf == d.dataOffsets.endBuf && d.dataOffsets.currOff == d.dataOffsets.endOff {
		d.done = true
		d.databufs.Free()
	}
	// We are at the end of the current buffer
	if d.dataOffsets.currBuf != d.dataOffsets.endBuf && d.dataOffsets.currOff == uint64(databuf.Len()) {
		d.dataOffsets.currOff = 0
		d.dataOffsets.currBuf++
	}
	return n
}

func (d *readResponseDecoder) writeToAndUpdateCRC(w io.Writer, updateCRC func([]byte)) (int64, error) {
	// For a completely empty message, just return 0
	if len(d.databufs) == 0 {
		return 0, nil
	}
	var written int64
	for !d.done {
		databuf := d.databufs[d.dataOffsets.currBuf]
		startOff := d.dataOffsets.currOff
		var b []byte
		if d.dataOffsets.currBuf == d.dataOffsets.endBuf {
			b = databuf.ReadOnlyData()[startOff:d.dataOffsets.endOff]
		} else {
			b = databuf.ReadOnlyData()[startOff:]
		}
		var n int
		// Write all remaining data from the current buffer
		n, err := w.Write(b)
		written += int64(n)
		updateCRC(b)
		if err != nil {
			return written, err
		}
		d.dataOffsets.currOff = 0
		// We've read all the data from this message.
		if d.dataOffsets.currBuf == d.dataOffsets.endBuf {
			d.done = true
			d.databufs.Free()
		} else {
			d.dataOffsets.currBuf++
		}
	}
	return written, nil
}

// Consume the next available tag in the input data and return the field number and type.
// Advances the relevant offsets in the data.
func (d *readResponseDecoder) consumeTag() (protowire.Number, protowire.Type, error) {
	tagBuf := d.peek()

	// Consume the next tag. This will tell us which field is next in the
	// buffer, its type, and how much space it takes up.
	fieldNum, fieldType, tagLength := protowire.ConsumeTag(tagBuf)
	if tagLength < 0 {
		return 0, 0, protowire.ParseError(tagLength)
	}
	// Update the offsets and current buffer depending on the tag length.
	if err := d.advanceOffset(uint64(tagLength)); err != nil {
		return 0, 0, fmt.Errorf("consuming tag: %w", err)
	}
	return fieldNum, fieldType, nil
}

// Consume a varint that represents the length of a bytes field. Return the length of
// the data, and advance the offsets by the length of the varint.
func (d *readResponseDecoder) consumeVarint() (uint64, error) {
	tagBuf := d.peek()

	// Consume the next tag. This will tell us which field is next in the
	// buffer, its type, and how much space it takes up.
	dataLength, tagLength := protowire.ConsumeVarint(tagBuf)
	if tagLength < 0 {
		return 0, protowire.ParseError(tagLength)
	}

	// Update the offsets and current buffer depending on the tag length.
	d.advanceOffset(uint64(tagLength))
	return dataLength, nil
}

func (d *readResponseDecoder) consumeFixed32() (uint32, error) {
	valueBuf := d.peek()

	// Consume the next tag. This will tell us which field is next in the
	// buffer, its type, and how much space it takes up.
	value, tagLength := protowire.ConsumeFixed32(valueBuf)
	if tagLength < 0 {
		return 0, protowire.ParseError(tagLength)
	}

	// Update the offsets and current buffer depending on the tag length.
	d.advanceOffset(uint64(tagLength))
	return value, nil
}

func (d *readResponseDecoder) consumeFixed64() (uint64, error) {
	valueBuf := d.peek()

	// Consume the next tag. This will tell us which field is next in the
	// buffer, its type, and how much space it takes up.
	value, tagLength := protowire.ConsumeFixed64(valueBuf)
	if tagLength < 0 {
		return 0, protowire.ParseError(tagLength)
	}

	// Update the offsets and current buffer depending on the tag length.
	d.advanceOffset(uint64(tagLength))
	return value, nil
}

// Consume any field values up to the end offset provided and don't return anything.
// This is used to skip any values which are not going to be used.
// msgEndOff is indexed in terms of the overall data across all buffers.
func (d *readResponseDecoder) consumeFieldValue(fieldNum protowire.Number, fieldType protowire.Type) error {
	// reimplement protowire.ConsumeFieldValue without the extra case for groups (which
	// are are complicted and not a thing in proto3).
	var err error
	switch fieldType {
	case protowire.VarintType:
		_, err = d.consumeVarint()
	case protowire.Fixed32Type:
		_, err = d.consumeFixed32()
	case protowire.Fixed64Type:
		_, err = d.consumeFixed64()
	case protowire.BytesType:
		_, err = d.consumeBytes()
	default:
		return fmt.Errorf("unknown field type %v in field %v", fieldType, fieldNum)
	}
	if err != nil {
		return fmt.Errorf("consuming field %v of type %v: %w", fieldNum, fieldType, err)
	}

	return nil
}

// Consume a bytes field from the input. Returns offsets for the data in the buffer slices
// and an error.
func (d *readResponseDecoder) consumeBytes() (bufferSliceOffsets, error) {
	// m is the length of the data past the tag.
	m, err := d.consumeVarint()
	if err != nil {
		return bufferSliceOffsets{}, fmt.Errorf("consuming bytes field: %w", err)
	}
	offsets := bufferSliceOffsets{
		startBuf: d.currBuf,
		startOff: d.currOff,
		currBuf:  d.currBuf,
		currOff:  d.currOff,
	}

	// Advance offsets to lengths of bytes field and capture where we end.
	d.advanceOffset(m)
	offsets.endBuf = d.currBuf
	offsets.endOff = d.currOff
	return offsets, nil
}

// Consume a bytes field from the input and copy into a new buffer if
// necessary (if the data is split across buffers in databuf).  This can be
// used to leverage proto.Unmarshal for small bytes fields (i.e. anything
// except object data).
func (d *readResponseDecoder) consumeBytesCopy() ([]byte, error) {
	// m is the length of the bytes data.
	m, err := d.consumeVarint()
	if err != nil {
		return nil, fmt.Errorf("consuming varint: %w", err)
	}
	// Copy the data into a buffer and advance the offset
	b := d.copyNextBytes(int(m))
	if err := d.advanceOffset(m); err != nil {
		return nil, fmt.Errorf("advancing offset: %w", err)
	}
	return b, nil
}

// readFullObjectResponse returns the ReadObjectResponse that is encoded in the
// wire-encoded message buffer b, or an error if the message is invalid.
// This must be used on the first recv of an object as it may contain all fields
// of ReadObjectResponse, and we use or pass on those fields to the user.
// This function is essentially identical to proto.Unmarshal, except it aliases
// the data in the input []byte. If the proto library adds a feature to
// Unmarshal that does that, this function can be dropped.
func (d *readResponseDecoder) readFullObjectResponse() error {
	msg := &storagepb.ReadObjectResponse{}

	// Loop over the entire message, extracting fields as we go. This does not
	// handle field concatenation, in which the contents of a single field
	// are split across multiple protobuf tags.
	for d.off < uint64(d.databufs.Len()) {
		fieldNum, fieldType, err := d.consumeTag()
		if err != nil {
			return fmt.Errorf("consuming next tag: %w", err)
		}

		// Unmarshal the field according to its type. Only fields that are not
		// nil will be present.
		switch {
		case fieldNum == checksummedDataField && fieldType == protowire.BytesType:
			// The ChecksummedData field was found. Initialize the struct.
			msg.ChecksummedData = &storagepb.ChecksummedData{}

			bytesFieldLen, err := d.consumeVarint()
			if err != nil {
				return fmt.Errorf("consuming bytes: %v", err)
			}

			var contentEndOff = d.off + bytesFieldLen
			for d.off < contentEndOff {
				gotNum, gotTyp, err := d.consumeTag()
				if err != nil {
					return fmt.Errorf("consuming checksummedData tag: %w", err)
				}

				switch {
				case gotNum == checksummedDataContentField && gotTyp == protowire.BytesType:
					// Get the offsets of the content bytes.
					d.dataOffsets, err = d.consumeBytes()
					if err != nil {
						return fmt.Errorf("invalid ReadObjectResponse.ChecksummedData.Content: %w", err)
					}
				case gotNum == checksummedDataCRC32CField && gotTyp == protowire.Fixed32Type:
					v, err := d.consumeFixed32()
					if err != nil {
						return fmt.Errorf("invalid ReadObjectResponse.ChecksummedData.Crc32C: %w", err)
					}
					msg.ChecksummedData.Crc32C = &v
				default:
					err := d.consumeFieldValue(gotNum, gotTyp)
					if err != nil {
						return fmt.Errorf("invalid field in ReadObjectResponse.ChecksummedData: %w", err)
					}
				}
			}
		case fieldNum == objectChecksumsField && fieldType == protowire.BytesType:
			// The field was found. Initialize the struct.
			msg.ObjectChecksums = &storagepb.ObjectChecksums{}
			// Consume the bytes and copy them into a single buffer if they are split across buffers.
			buf, err := d.consumeBytesCopy()
			if err != nil {
				return fmt.Errorf("invalid ReadObjectResponse.ObjectChecksums: %v", err)
			}
			// Unmarshal.
			if err := proto.Unmarshal(buf, msg.ObjectChecksums); err != nil {
				return err
			}
		case fieldNum == contentRangeField && fieldType == protowire.BytesType:
			msg.ContentRange = &storagepb.ContentRange{}
			buf, err := d.consumeBytesCopy()
			if err != nil {
				return fmt.Errorf("invalid ReadObjectResponse.ContentRange: %v", err)
			}
			if err := proto.Unmarshal(buf, msg.ContentRange); err != nil {
				return err
			}
		case fieldNum == metadataField && fieldType == protowire.BytesType:
			msg.Metadata = &storagepb.Object{}

			buf, err := d.consumeBytesCopy()
			if err != nil {
				return fmt.Errorf("invalid ReadObjectResponse.Metadata: %v", err)
			}

			if err := proto.Unmarshal(buf, msg.Metadata); err != nil {
				return err
			}
		default:
			err := d.consumeFieldValue(fieldNum, fieldType)
			if err != nil {
				return fmt.Errorf("invalid field in ReadObjectResponse: %w", err)
			}
		}
	}
	d.msg = msg
	return nil
}

// reopenStream "closes" the existing stream and attempts to reopen a stream and
// sets the Reader's stream and cancelStream properties in the process.
func (r *gRPCReader) reopenStream() error {
	// Close existing stream and initialize new stream with updated offset.
	r.Close()

	res, cancel, err := r.reopen(r.seen)
	if err != nil {
		return err
	}
	r.stream = res.stream
	r.currMsg = res.decoder
	r.cancel = cancel
	return nil
}

func newGRPCWriter(c *grpcStorageClient, params *openWriterParams, r io.Reader) *gRPCWriter {
	size := params.chunkSize

	// Round up chunksize to nearest 256KiB
	if size%googleapi.MinUploadChunkSize != 0 {
		size += googleapi.MinUploadChunkSize - (size % googleapi.MinUploadChunkSize)
	}

	// A completely bufferless upload is not possible as it is in JSON because
	// the buffer must be provided to the message. However use the minimum size
	// possible in this case.
	if params.chunkSize == 0 {
		size = googleapi.MinUploadChunkSize
	}

	return &gRPCWriter{
		buf:                   make([]byte, size),
		c:                     c,
		ctx:                   params.ctx,
		reader:                r,
		bucket:                params.bucket,
		attrs:                 params.attrs,
		conds:                 params.conds,
		encryptionKey:         params.encryptionKey,
		sendCRC32C:            params.sendCRC32C,
		chunkSize:             params.chunkSize,
		forceEmptyContentType: params.forceEmptyContentType,
	}
}

// gRPCWriter is a wrapper around the the gRPC client-stream API that manages
// sending chunks of data provided by the user over the stream.
type gRPCWriter struct {
	c      *grpcStorageClient
	buf    []byte
	reader io.Reader

	ctx context.Context

	bucket        string
	attrs         *ObjectAttrs
	conds         *Conditions
	encryptionKey []byte
	settings      *settings

	sendCRC32C            bool
	chunkSize             int
	forceEmptyContentType bool

	// The gRPC client-stream used for sending buffers.
	stream storagepb.Storage_BidiWriteObjectClient

	// The Resumable Upload ID started by a gRPC-based Writer.
	upid string
}

// startResumableUpload initializes a Resumable Upload with gRPC and sets the
// upload ID on the Writer.
func (w *gRPCWriter) startResumableUpload() error {
	spec, err := w.writeObjectSpec()
	if err != nil {
		return err
	}
	req := &storagepb.StartResumableWriteRequest{
		WriteObjectSpec:           spec,
		CommonObjectRequestParams: toProtoCommonObjectRequestParams(w.encryptionKey),
	}
	// TODO: Currently the checksums are only sent on the request to initialize
	// the upload, but in the future, we must also support sending it
	// on the *last* message of the stream.
	req.ObjectChecksums = toProtoChecksums(w.sendCRC32C, w.attrs)
	return run(w.ctx, func(ctx context.Context) error {
		upres, err := w.c.raw.StartResumableWrite(w.ctx, req)
		w.upid = upres.GetUploadId()
		return err
	}, w.settings.retry, w.settings.idempotent)
}

// queryProgress is a helper that queries the status of the resumable upload
// associated with the given upload ID.
func (w *gRPCWriter) queryProgress() (int64, error) {
	var persistedSize int64
	err := run(w.ctx, func(ctx context.Context) error {
		q, err := w.c.raw.QueryWriteStatus(w.ctx, &storagepb.QueryWriteStatusRequest{
			UploadId: w.upid,
		})
		persistedSize = q.GetPersistedSize()
		return err
	}, w.settings.retry, true)

	// q.GetCommittedSize() will return 0 if q is nil.
	return persistedSize, err
}

// uploadBuffer uploads the buffer at the given offset using a bi-directional
// Write stream. It will open a new stream if necessary (on the first call or
// after resuming from failure). The resulting write offset after uploading the
// buffer is returned, as well as well as the final Object if the upload is
// completed.
//
// Returns object, persisted size, and any error that is not retriable.
func (w *gRPCWriter) uploadBuffer(recvd int, start int64, doneReading bool) (*storagepb.Object, int64, error) {
	var shouldRetry = ShouldRetry
	if w.settings.retry != nil && w.settings.retry.shouldRetry != nil {
		shouldRetry = w.settings.retry.shouldRetry
	}

	var err error
	var lastWriteOfEntireObject bool

	sent := 0
	writeOffset := start

	toWrite := w.buf[:recvd]

	// Send a request with as many bytes as possible.
	// Loop until all bytes are sent.
sendBytes: // label this loop so that we can use a continue statement from a nested block
	for {
		bytesNotYetSent := recvd - sent
		remainingDataFitsInSingleReq := bytesNotYetSent <= maxPerMessageWriteSize

		if remainingDataFitsInSingleReq && doneReading {
			lastWriteOfEntireObject = true
		}

		// Send the maximum amount of bytes we can, unless we don't have that many.
		bytesToSendInCurrReq := maxPerMessageWriteSize
		if remainingDataFitsInSingleReq {
			bytesToSendInCurrReq = bytesNotYetSent
		}

		// Prepare chunk section for upload.
		data := toWrite[sent : sent+bytesToSendInCurrReq]

		req := &storagepb.BidiWriteObjectRequest{
			Data: &storagepb.BidiWriteObjectRequest_ChecksummedData{
				ChecksummedData: &storagepb.ChecksummedData{
					Content: data,
				},
			},
			WriteOffset: writeOffset,
			FinishWrite: lastWriteOfEntireObject,
			Flush:       remainingDataFitsInSingleReq && !lastWriteOfEntireObject,
			StateLookup: remainingDataFitsInSingleReq && !lastWriteOfEntireObject,
		}

		// Open a new stream if necessary and set the first_message field on
		// the request. The first message on the WriteObject stream must either
		// be the Object or the Resumable Upload ID.
		if w.stream == nil {
			hds := []string{"x-goog-request-params", fmt.Sprintf("bucket=projects/_/buckets/%s", url.QueryEscape(w.bucket))}
			ctx := gax.InsertMetadataIntoOutgoingContext(w.ctx, hds...)

			w.stream, err = w.c.raw.BidiWriteObject(ctx)
			if err != nil {
				return nil, 0, err
			}

			if w.upid != "" { // resumable upload
				req.FirstMessage = &storagepb.BidiWriteObjectRequest_UploadId{UploadId: w.upid}
			} else { // non-resumable
				spec, err := w.writeObjectSpec()
				if err != nil {
					return nil, 0, err
				}
				req.FirstMessage = &storagepb.BidiWriteObjectRequest_WriteObjectSpec{
					WriteObjectSpec: spec,
				}
				req.CommonObjectRequestParams = toProtoCommonObjectRequestParams(w.encryptionKey)
				// For a non-resumable upload, checksums must be sent in this message.
				// TODO: Currently the checksums are only sent on the first message
				// of the stream, but in the future, we must also support sending it
				// on the *last* message of the stream (instead of the first).
				req.ObjectChecksums = toProtoChecksums(w.sendCRC32C, w.attrs)
			}
		}

		err = w.stream.Send(req)
		if err == io.EOF {
			// err was io.EOF. The client-side of a stream only gets an EOF on Send
			// when the backend closes the stream and wants to return an error
			// status.

			// Receive from the stream Recv() until it returns a non-nil error
			// to receive the server's status as an error. We may get multiple
			// messages before the error due to buffering.
			err = nil
			for err == nil {
				_, err = w.stream.Recv()
			}
			// Drop the stream reference as a new one will need to be created if
			// we retry.
			w.stream = nil

			// Retriable errors mean we should start over and attempt to
			// resend the entire buffer via a new stream.
			// If not retriable, falling through will return the error received.
			if shouldRetry(err) {
				// TODO: Add test case for failure modes of querying progress.
				writeOffset, err = w.determineOffset(start)
				if err != nil {
					return nil, 0, err
				}
				sent = int(writeOffset) - int(start)

				// Continue sending requests, opening a new stream and resending
				// any bytes not yet persisted as per QueryWriteStatus
				continue sendBytes
			}
		}
		if err != nil {
			return nil, 0, err
		}

		// Update the immediate stream's sent total and the upload offset with
		// the data sent.
		sent += len(data)
		writeOffset += int64(len(data))

		// Not done sending data, do not attempt to commit it yet, loop around
		// and send more data.
		if recvd-sent > 0 {
			continue sendBytes
		}

		// The buffer has been uploaded and there is still more data to be
		// uploaded, but this is not a resumable upload session. Therefore,
		// don't check persisted data.
		if !lastWriteOfEntireObject && w.chunkSize == 0 {
			return nil, writeOffset, nil
		}

		// Done sending the data in the buffer (remainingDataFitsInSingleReq
		// should == true if we reach this code).
		// If we are done sending the whole object, close the stream and get the final
		// object. Otherwise, receive from the stream to confirm the persisted data.
		if !lastWriteOfEntireObject {
			resp, err := w.stream.Recv()

			// Retriable errors mean we should start over and attempt to
			// resend the entire buffer via a new stream.
			// If not retriable, falling through will return the error received
			// from closing the stream.
			if shouldRetry(err) {
				writeOffset, err = w.determineOffset(start)
				if err != nil {
					return nil, 0, err
				}
				sent = int(writeOffset) - int(start)

				// Drop the stream reference as a new one will need to be created.
				w.stream = nil

				continue sendBytes
			}
			if err != nil {
				return nil, 0, err
			}

			if resp.GetPersistedSize() != writeOffset {
				// Retry if not all bytes were persisted.
				writeOffset = resp.GetPersistedSize()
				sent = int(writeOffset) - int(start)
				continue sendBytes
			}
		} else {
			// If the object is done uploading, close the send stream to signal
			// to the server that we are done sending so that we can receive
			// from the stream without blocking.
			err = w.stream.CloseSend()
			if err != nil {
				// CloseSend() retries the send internally. It never returns an
				// error in the current implementation, but we check it anyway in
				// case that it does in the future.
				return nil, 0, err
			}

			// Stream receives do not block once send is closed, but we may not
			// receive the response with the object right away; loop until we
			// receive the object or error out.
			var obj *storagepb.Object
			for obj == nil {
				resp, err := w.stream.Recv()
				if shouldRetry(err) {
					writeOffset, err = w.determineOffset(start)
					if err != nil {
						return nil, 0, err
					}
					sent = int(writeOffset) - int(start)
					w.stream = nil
					continue sendBytes
				}
				if err != nil {
					return nil, 0, err
				}

				obj = resp.GetResource()
			}

			// Even though we received the object response, continue reading
			// until we receive a non-nil error, to ensure the stream does not
			// leak even if the context isn't cancelled. See:
			// https://pkg.go.dev/google.golang.org/grpc#ClientConn.NewStream
			for err == nil {
				_, err = w.stream.Recv()
			}

			return obj, writeOffset, nil
		}

		return nil, writeOffset, nil
	}
}

// determineOffset either returns the offset given to it in the case of a simple
// upload, or queries the write status in the case a resumable upload is being
// used.
func (w *gRPCWriter) determineOffset(offset int64) (int64, error) {
	// For a Resumable Upload, we must start from however much data
	// was committed.
	if w.upid != "" {
		committed, err := w.queryProgress()
		if err != nil {
			return 0, err
		}
		offset = committed
	}
	return offset, nil
}

// writeObjectSpec constructs a WriteObjectSpec proto using the Writer's
// ObjectAttrs and applies its Conditions. This is only used for gRPC.
func (w *gRPCWriter) writeObjectSpec() (*storagepb.WriteObjectSpec, error) {
	// To avoid modifying the ObjectAttrs embeded in the calling writer, deref
	// the ObjectAttrs pointer to make a copy, then assign the desired name to
	// the attribute.
	attrs := *w.attrs

	spec := &storagepb.WriteObjectSpec{
		Resource: attrs.toProtoObject(w.bucket),
	}
	// WriteObject doesn't support the generation condition, so use default.
	if err := applyCondsProto("WriteObject", defaultGen, w.conds, spec); err != nil {
		return nil, err
	}
	return spec, nil
}

// read copies the data in the reader to the given buffer and reports how much
// data was read into the buffer and if there is no more data to read (EOF).
// Furthermore, if the attrs.ContentType is unset, the first bytes of content
// will be sniffed for a matching content type unless forceEmptyContentType is enabled.
func (w *gRPCWriter) read() (int, bool, error) {
	if w.attrs.ContentType == "" && !w.forceEmptyContentType {
		w.reader, w.attrs.ContentType = gax.DetermineContentType(w.reader)
	}
	// Set n to -1 to start the Read loop.
	var n, recvd int = -1, 0
	var err error
	for err == nil && n != 0 {
		// The routine blocks here until data is received.
		n, err = w.reader.Read(w.buf[recvd:])
		recvd += n
	}
	var done bool
	if err == io.EOF {
		done = true
		err = nil
	}
	return recvd, done, err
}

func checkCanceled(err error) error {
	if status.Code(err) == codes.Canceled {
		return context.Canceled
	}

	return err
}
