package customizations

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"net/url"
	"strings"

	"github.com/aws/smithy-go/middleware"
	"github.com/aws/smithy-go/transport/http"

	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/internal/s3shared"
	"github.com/aws/aws-sdk-go-v2/service/internal/s3shared/arn"
	s3arn "github.com/aws/aws-sdk-go-v2/service/s3/internal/arn"
)

const (
	s3AccessPoint  = "s3-accesspoint"
	s3ObjectLambda = "s3-object-lambda"
)

// processARNResource is used to process an ARN resource.
type processARNResource struct {

	// UseARNRegion indicates if region parsed from an ARN should be used.
	UseARNRegion bool

	// UseAccelerate indicates if s3 transfer acceleration is enabled
	UseAccelerate bool

	// UseDualstack instructs if s3 dualstack endpoint config is enabled
	UseDualstack bool

	// EndpointResolver used to resolve endpoints. This may be a custom endpoint resolver
	EndpointResolver EndpointResolver

	// EndpointResolverOptions used by endpoint resolver
	EndpointResolverOptions EndpointResolverOptions
}

// ID returns the middleware ID.
func (*processARNResource) ID() string { return "S3:ProcessARNResource" }

func (m *processARNResource) HandleSerialize(
	ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler,
) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	// check if arn was provided, if not skip this middleware
	arnValue, ok := s3shared.GetARNResourceFromContext(ctx)
	if !ok {
		return next.HandleSerialize(ctx, in)
	}

	req, ok := in.Request.(*http.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown request type %T", req)
	}

	// parse arn into an endpoint arn wrt to service
	resource, err := s3arn.ParseEndpointARN(arnValue)
	if err != nil {
		return out, metadata, err
	}

	// build a resource request struct
	resourceRequest := s3shared.ResourceRequest{
		Resource:      resource,
		UseARNRegion:  m.UseARNRegion,
		RequestRegion: awsmiddleware.GetRegion(ctx),
		SigningRegion: awsmiddleware.GetSigningRegion(ctx),
		PartitionID:   awsmiddleware.GetPartitionID(ctx),
	}

	// validate resource request
	if err := validateResourceRequest(resourceRequest); err != nil {
		return out, metadata, err
	}

	// switch to correct endpoint updater
	switch tv := resource.(type) {
	case arn.AccessPointARN:
		// check if accelerate
		if m.UseAccelerate {
			return out, metadata, s3shared.NewClientConfiguredForAccelerateError(tv,
				resourceRequest.PartitionID, resourceRequest.RequestRegion, nil)
		}

		// fetch arn region to resolve request
		resolveRegion := tv.Region
		// check if request region is FIPS
		if resourceRequest.UseFips() {
			// if use arn region is enabled and request signing region is not same as arn region
			if m.UseARNRegion && resourceRequest.IsCrossRegion() {
				// FIPS with cross region is not supported, the SDK must fail
				// because there is no well defined method for SDK to construct a
				// correct FIPS endpoint.
				return out, metadata,
					s3shared.NewClientConfiguredForCrossRegionFIPSError(
						tv,
						resourceRequest.PartitionID,
						resourceRequest.RequestRegion,
						nil,
					)
			}

			// if use arn region is NOT set, we should use the request region
			resolveRegion = resourceRequest.RequestRegion
		}

		// build access point request
		ctx, err = buildAccessPointRequest(ctx, accesspointOptions{
			processARNResource: *m,
			request:            req,
			resource:           tv,
			resolveRegion:      resolveRegion,
			partitionID:        resourceRequest.PartitionID,
			requestRegion:      resourceRequest.RequestRegion,
		})
		if err != nil {
			return out, metadata, err
		}

	case arn.S3ObjectLambdaAccessPointARN:
		// check if accelerate
		if m.UseAccelerate {
			return out, metadata, s3shared.NewClientConfiguredForAccelerateError(tv,
				resourceRequest.PartitionID, resourceRequest.RequestRegion, nil)
		}

		// check if dualstack
		if m.UseDualstack {
			return out, metadata, s3shared.NewClientConfiguredForDualStackError(tv,
				resourceRequest.PartitionID, resourceRequest.RequestRegion, nil)
		}

		// fetch arn region to resolve request
		resolveRegion := tv.Region

		if resourceRequest.UseFips() {
			// if use arn region is enabled and request signing region is not same as arn region
			if m.UseARNRegion && resourceRequest.IsCrossRegion() {
				// FIPS with cross region is not supported, the SDK must fail
				// because there is no well defined method for SDK to construct a
				// correct FIPS endpoint.
				return out, metadata,
					s3shared.NewClientConfiguredForCrossRegionFIPSError(
						tv,
						resourceRequest.PartitionID,
						resourceRequest.RequestRegion,
						nil,
					)
			}

			// if use arn region is NOT set, we should use the request region
			resolveRegion = resourceRequest.RequestRegion
		}

		// build access point request
		ctx, err = buildS3ObjectLambdaAccessPointRequest(ctx, accesspointOptions{
			processARNResource: *m,
			request:            req,
			resource:           tv.AccessPointARN,
			resolveRegion:      resolveRegion,
			partitionID:        resourceRequest.PartitionID,
			requestRegion:      resourceRequest.RequestRegion,
		})
		if err != nil {
			return out, metadata, err
		}

	// process outpost accesspoint ARN
	case arn.OutpostAccessPointARN:
		// check if accelerate
		if m.UseAccelerate {
			return out, metadata, s3shared.NewClientConfiguredForAccelerateError(tv,
				resourceRequest.PartitionID, resourceRequest.RequestRegion, nil)
		}

		// check if dual stack
		if m.UseDualstack {
			return out, metadata, s3shared.NewClientConfiguredForDualStackError(tv,
				resourceRequest.PartitionID, resourceRequest.RequestRegion, nil)
		}

		// check if request region is FIPS
		if resourceRequest.UseFips() {
			return out, metadata, s3shared.NewFIPSConfigurationError(tv, resourceRequest.PartitionID,
				resourceRequest.RequestRegion, nil)
		}

		// build outpost access point request
		ctx, err = buildOutpostAccessPointRequest(ctx, outpostAccessPointOptions{
			processARNResource: *m,
			resource:           tv,
			request:            req,
			partitionID:        resourceRequest.PartitionID,
			requestRegion:      resourceRequest.RequestRegion,
		})
		if err != nil {
			return out, metadata, err
		}

	default:
		return out, metadata, s3shared.NewInvalidARNError(resource, nil)
	}

	return next.HandleSerialize(ctx, in)
}

// validate if s3 resource and request config is compatible.
func validateResourceRequest(resourceRequest s3shared.ResourceRequest) error {
	// check if resourceRequest leads to a cross partition error
	v, err := resourceRequest.IsCrossPartition()
	if err != nil {
		return err
	}
	if v {
		// if cross partition
		return s3shared.NewClientPartitionMismatchError(resourceRequest.Resource,
			resourceRequest.PartitionID, resourceRequest.RequestRegion, nil)
	}

	// check if resourceRequest leads to a cross region error
	if !resourceRequest.AllowCrossRegion() && resourceRequest.IsCrossRegion() {
		// if cross region, but not use ARN region is not enabled
		return s3shared.NewClientRegionMismatchError(resourceRequest.Resource,
			resourceRequest.PartitionID, resourceRequest.RequestRegion, nil)
	}

	return nil
}

// === Accesspoint ==========

type accesspointOptions struct {
	processARNResource
	request       *http.Request
	resource      arn.AccessPointARN
	resolveRegion string
	partitionID   string
	requestRegion string
}

func buildAccessPointRequest(ctx context.Context, options accesspointOptions) (context.Context, error) {
	tv := options.resource
	req := options.request
	resolveRegion := options.resolveRegion

	resolveService := tv.Service

	// resolve endpoint
	endpoint, err := options.EndpointResolver.ResolveEndpoint(resolveRegion, options.EndpointResolverOptions)
	if err != nil {
		return ctx, s3shared.NewFailedToResolveEndpointError(
			tv,
			options.partitionID,
			options.requestRegion,
			err,
		)
	}

	// assign resolved endpoint url to request url
	req.URL, err = url.Parse(endpoint.URL)
	if err != nil {
		return ctx, fmt.Errorf("failed to parse endpoint URL: %w", err)
	}

	if len(endpoint.SigningName) != 0 && endpoint.Source == aws.EndpointSourceCustom {
		ctx = awsmiddleware.SetSigningName(ctx, endpoint.SigningName)
	} else {
		// Must sign with s3-object-lambda
		ctx = awsmiddleware.SetSigningName(ctx, resolveService)
	}

	if len(endpoint.SigningRegion) != 0 {
		ctx = awsmiddleware.SetSigningRegion(ctx, endpoint.SigningRegion)
	} else {
		ctx = awsmiddleware.SetSigningRegion(ctx, resolveRegion)
	}

	// update serviceID to "s3-accesspoint"
	ctx = awsmiddleware.SetServiceID(ctx, s3AccessPoint)

	// disable host prefix behavior
	ctx = http.DisableEndpointHostPrefix(ctx, true)

	// remove the serialized arn in place of /{Bucket}
	ctx = setBucketToRemoveOnContext(ctx, tv.String())

	// skip arn processing, if arn region resolves to a immutable endpoint
	if endpoint.HostnameImmutable {
		return ctx, nil
	}

	updateS3HostForS3AccessPoint(req)

	ctx, err = buildAccessPointHostPrefix(ctx, req, tv)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func buildS3ObjectLambdaAccessPointRequest(ctx context.Context, options accesspointOptions) (context.Context, error) {
	tv := options.resource
	req := options.request
	resolveRegion := options.resolveRegion

	resolveService := tv.Service

	// resolve endpoint
	endpoint, err := options.EndpointResolver.ResolveEndpoint(resolveRegion, options.EndpointResolverOptions)
	if err != nil {
		return ctx, s3shared.NewFailedToResolveEndpointError(
			tv,
			options.partitionID,
			options.requestRegion,
			err,
		)
	}

	// assign resolved endpoint url to request url
	req.URL, err = url.Parse(endpoint.URL)
	if err != nil {
		return ctx, fmt.Errorf("failed to parse endpoint URL: %w", err)
	}

	if len(endpoint.SigningName) != 0 && endpoint.Source == aws.EndpointSourceCustom {
		ctx = awsmiddleware.SetSigningName(ctx, endpoint.SigningName)
	} else {
		// Must sign with s3-object-lambda
		ctx = awsmiddleware.SetSigningName(ctx, resolveService)
	}

	if len(endpoint.SigningRegion) != 0 {
		ctx = awsmiddleware.SetSigningRegion(ctx, endpoint.SigningRegion)
	} else {
		ctx = awsmiddleware.SetSigningRegion(ctx, resolveRegion)
	}

	// update serviceID to "s3-object-lambda"
	ctx = awsmiddleware.SetServiceID(ctx, s3ObjectLambda)

	// disable host prefix behavior
	ctx = http.DisableEndpointHostPrefix(ctx, true)

	// remove the serialized arn in place of /{Bucket}
	ctx = setBucketToRemoveOnContext(ctx, tv.String())

	// skip arn processing, if arn region resolves to a immutable endpoint
	if endpoint.HostnameImmutable {
		return ctx, nil
	}

	if endpoint.Source == aws.EndpointSourceServiceMetadata {
		updateS3HostForS3ObjectLambda(req)
	}

	ctx, err = buildAccessPointHostPrefix(ctx, req, tv)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func buildAccessPointHostPrefix(ctx context.Context, req *http.Request, tv arn.AccessPointARN) (context.Context, error) {
	// add host prefix for access point
	accessPointHostPrefix := tv.AccessPointName + "-" + tv.AccountID + "."
	req.URL.Host = accessPointHostPrefix + req.URL.Host
	if len(req.Host) > 0 {
		req.Host = accessPointHostPrefix + req.Host
	}

	// validate the endpoint host
	if err := http.ValidateEndpointHost(req.URL.Host); err != nil {
		return ctx, s3shared.NewInvalidARNError(tv, err)
	}

	return ctx, nil
}

// ====== Outpost Accesspoint ========

type outpostAccessPointOptions struct {
	processARNResource
	request       *http.Request
	resource      arn.OutpostAccessPointARN
	partitionID   string
	requestRegion string
}

func buildOutpostAccessPointRequest(ctx context.Context, options outpostAccessPointOptions) (context.Context, error) {
	tv := options.resource
	req := options.request

	resolveRegion := tv.Region
	resolveService := tv.Service
	endpointsID := resolveService
	if strings.EqualFold(resolveService, "s3-outposts") {
		// assign endpoints ID as "S3"
		endpointsID = "s3"
	}

	// resolve regional endpoint for resolved region.
	endpoint, err := options.EndpointResolver.ResolveEndpoint(resolveRegion, options.EndpointResolverOptions)
	if err != nil {
		return ctx, s3shared.NewFailedToResolveEndpointError(
			tv,
			options.partitionID,
			options.requestRegion,
			err,
		)
	}

	// assign resolved endpoint url to request url
	req.URL, err = url.Parse(endpoint.URL)
	if err != nil {
		return ctx, fmt.Errorf("failed to parse endpoint URL: %w", err)
	}

	// assign resolved service from arn as signing name
	if len(endpoint.SigningName) != 0 && endpoint.Source == aws.EndpointSourceCustom {
		ctx = awsmiddleware.SetSigningName(ctx, endpoint.SigningName)
	} else {
		ctx = awsmiddleware.SetSigningName(ctx, resolveService)
	}

	if len(endpoint.SigningRegion) != 0 {
		// redirect signer to use resolved endpoint signing name and region
		ctx = awsmiddleware.SetSigningRegion(ctx, endpoint.SigningRegion)
	} else {
		ctx = awsmiddleware.SetSigningRegion(ctx, resolveRegion)
	}

	// update serviceID to resolved service id
	ctx = awsmiddleware.SetServiceID(ctx, resolveService)

	// disable host prefix behavior
	ctx = http.DisableEndpointHostPrefix(ctx, true)

	// remove the serialized arn in place of /{Bucket}
	ctx = setBucketToRemoveOnContext(ctx, tv.String())

	// skip further customizations, if arn region resolves to a immutable endpoint
	if endpoint.HostnameImmutable {
		return ctx, nil
	}

	updateHostPrefix(req, endpointsID, resolveService)

	// add host prefix for s3-outposts
	outpostAPHostPrefix := tv.AccessPointName + "-" + tv.AccountID + "." + tv.OutpostID + "."
	req.URL.Host = outpostAPHostPrefix + req.URL.Host
	if len(req.Host) > 0 {
		req.Host = outpostAPHostPrefix + req.Host
	}

	// validate the endpoint host
	if err := http.ValidateEndpointHost(req.URL.Host); err != nil {
		return ctx, s3shared.NewInvalidARNError(tv, err)
	}

	return ctx, nil
}
