package s3

import (
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/internal/s3shared"
	"github.com/aws/aws-sdk-go/internal/s3shared/arn"
	"github.com/aws/aws-sdk-go/private/protocol"
)

const (
	accessPointPrefixLabel    = "accesspoint"
	accountIDPrefixLabel      = "accountID"
	accessPointPrefixTemplate = "{" + accessPointPrefixLabel + "}-{" + accountIDPrefixLabel + "}."

	outpostPrefixLabel               = "outpost"
	outpostAccessPointPrefixTemplate = accessPointPrefixTemplate + "{" + outpostPrefixLabel + "}."
)

// hasCustomEndpoint returns true if endpoint is a custom endpoint
func hasCustomEndpoint(r *request.Request) bool {
	return len(aws.StringValue(r.Config.Endpoint)) > 0
}

// accessPointEndpointBuilder represents the endpoint builder for access point arn
type accessPointEndpointBuilder arn.AccessPointARN

// build builds the endpoint for corresponding access point arn
//
// For building an endpoint from access point arn, format used is:
// - Access point endpoint format : {accesspointName}-{accountId}.s3-accesspoint.{region}.{dnsSuffix}
// - example : myaccesspoint-012345678901.s3-accesspoint.us-west-2.amazonaws.com
//
// Access Point Endpoint requests are signed using "s3" as signing name.
func (a accessPointEndpointBuilder) build(req *request.Request) error {
	resolveService := arn.AccessPointARN(a).Service
	resolveRegion := arn.AccessPointARN(a).Region

	endpoint, err := resolveRegionalEndpoint(req, resolveRegion, "", resolveService)
	if err != nil {
		return s3shared.NewFailedToResolveEndpointError(arn.AccessPointARN(a),
			req.ClientInfo.PartitionID, resolveRegion, err)
	}

	endpoint.URL = endpoints.AddScheme(endpoint.URL, aws.BoolValue(req.Config.DisableSSL))

	if !hasCustomEndpoint(req) {
		if err = updateRequestEndpoint(req, endpoint.URL); err != nil {
			return err
		}

		// dual stack provided by endpoint resolver
		updateS3HostForS3AccessPoint(req)
	}

	protocol.HostPrefixBuilder{
		Prefix:   accessPointPrefixTemplate,
		LabelsFn: a.hostPrefixLabelValues,
	}.Build(req)

	// signer redirection
	redirectSigner(req, endpoint.SigningName, endpoint.SigningRegion)

	err = protocol.ValidateEndpointHost(req.Operation.Name, req.HTTPRequest.URL.Host)
	if err != nil {
		return s3shared.NewInvalidARNError(arn.AccessPointARN(a), err)
	}

	return nil
}

func (a accessPointEndpointBuilder) hostPrefixLabelValues() map[string]string {
	return map[string]string{
		accessPointPrefixLabel: arn.AccessPointARN(a).AccessPointName,
		accountIDPrefixLabel:   arn.AccessPointARN(a).AccountID,
	}
}

// s3ObjectLambdaAccessPointEndpointBuilder represents the endpoint builder for an s3 object lambda access point arn
type s3ObjectLambdaAccessPointEndpointBuilder arn.S3ObjectLambdaAccessPointARN

// build builds the endpoint for corresponding access point arn
//
// For building an endpoint from access point arn, format used is:
// - Access point endpoint format : {accesspointName}-{accountId}.s3-object-lambda.{region}.{dnsSuffix}
// - example : myaccesspoint-012345678901.s3-object-lambda.us-west-2.amazonaws.com
//
// Access Point Endpoint requests are signed using "s3-object-lambda" as signing name.
func (a s3ObjectLambdaAccessPointEndpointBuilder) build(req *request.Request) error {
	resolveRegion := arn.S3ObjectLambdaAccessPointARN(a).Region

	endpoint, err := resolveRegionalEndpoint(req, resolveRegion, "", EndpointsID)
	if err != nil {
		return s3shared.NewFailedToResolveEndpointError(arn.S3ObjectLambdaAccessPointARN(a),
			req.ClientInfo.PartitionID, resolveRegion, err)
	}

	endpoint.URL = endpoints.AddScheme(endpoint.URL, aws.BoolValue(req.Config.DisableSSL))

	endpoint.SigningName = s3ObjectsLambdaNamespace

	if !hasCustomEndpoint(req) {
		if err = updateRequestEndpoint(req, endpoint.URL); err != nil {
			return err
		}

		updateS3HostPrefixForS3ObjectLambda(req)
	}

	protocol.HostPrefixBuilder{
		Prefix:   accessPointPrefixTemplate,
		LabelsFn: a.hostPrefixLabelValues,
	}.Build(req)

	// signer redirection
	redirectSigner(req, endpoint.SigningName, endpoint.SigningRegion)

	err = protocol.ValidateEndpointHost(req.Operation.Name, req.HTTPRequest.URL.Host)
	if err != nil {
		return s3shared.NewInvalidARNError(arn.S3ObjectLambdaAccessPointARN(a), err)
	}

	return nil
}

func (a s3ObjectLambdaAccessPointEndpointBuilder) hostPrefixLabelValues() map[string]string {
	return map[string]string{
		accessPointPrefixLabel: arn.S3ObjectLambdaAccessPointARN(a).AccessPointName,
		accountIDPrefixLabel:   arn.S3ObjectLambdaAccessPointARN(a).AccountID,
	}
}

// outpostAccessPointEndpointBuilder represents the Endpoint builder for outpost access point arn.
type outpostAccessPointEndpointBuilder arn.OutpostAccessPointARN

// build builds an endpoint corresponding to the outpost access point arn.
//
// For building an endpoint from outpost access point arn, format used is:
// - Outpost access point endpoint format : {accesspointName}-{accountId}.{outpostId}.s3-outposts.{region}.{dnsSuffix}
// - example : myaccesspoint-012345678901.op-01234567890123456.s3-outposts.us-west-2.amazonaws.com
//
// Outpost AccessPoint Endpoint request are signed using "s3-outposts" as signing name.
func (o outpostAccessPointEndpointBuilder) build(req *request.Request) error {
	resolveRegion := o.Region
	resolveService := o.Service

	endpointsID := resolveService
	if resolveService == s3OutpostsNamespace {
		endpointsID = "s3"
	}

	endpoint, err := resolveRegionalEndpoint(req, resolveRegion, "", endpointsID)
	if err != nil {
		return s3shared.NewFailedToResolveEndpointError(o,
			req.ClientInfo.PartitionID, resolveRegion, err)
	}

	endpoint.URL = endpoints.AddScheme(endpoint.URL, aws.BoolValue(req.Config.DisableSSL))

	if !hasCustomEndpoint(req) {
		if err = updateRequestEndpoint(req, endpoint.URL); err != nil {
			return err
		}
		updateHostPrefix(req, endpointsID, resolveService)
	}

	protocol.HostPrefixBuilder{
		Prefix:   outpostAccessPointPrefixTemplate,
		LabelsFn: o.hostPrefixLabelValues,
	}.Build(req)

	// set the signing region, name to resolved names from ARN
	redirectSigner(req, resolveService, resolveRegion)

	err = protocol.ValidateEndpointHost(req.Operation.Name, req.HTTPRequest.URL.Host)
	if err != nil {
		return s3shared.NewInvalidARNError(o, err)
	}

	return nil
}

func (o outpostAccessPointEndpointBuilder) hostPrefixLabelValues() map[string]string {
	return map[string]string{
		accessPointPrefixLabel: o.AccessPointName,
		accountIDPrefixLabel:   o.AccountID,
		outpostPrefixLabel:     o.OutpostID,
	}
}

func resolveRegionalEndpoint(r *request.Request, region, resolvedRegion, endpointsID string) (endpoints.ResolvedEndpoint, error) {
	return r.Config.EndpointResolver.EndpointFor(endpointsID, region, func(opts *endpoints.Options) {
		opts.DisableSSL = aws.BoolValue(r.Config.DisableSSL)
		opts.UseDualStack = aws.BoolValue(r.Config.UseDualStack)
		opts.UseDualStackEndpoint = r.Config.UseDualStackEndpoint
		opts.UseFIPSEndpoint = r.Config.UseFIPSEndpoint
		opts.S3UsEast1RegionalEndpoint = endpoints.RegionalS3UsEast1Endpoint
		opts.ResolvedRegion = resolvedRegion
		opts.Logger = r.Config.Logger
		opts.LogDeprecated = r.Config.LogLevel.Matches(aws.LogDebugWithDeprecated)
	})
}

func updateRequestEndpoint(r *request.Request, endpoint string) (err error) {
	r.HTTPRequest.URL, err = url.Parse(endpoint + r.Operation.HTTPPath)
	if err != nil {
		return awserr.New(request.ErrCodeSerialization,
			"failed to parse endpoint URL", err)
	}

	return nil
}

// redirectSigner sets signing name, signing region for a request
func redirectSigner(req *request.Request, signingName string, signingRegion string) {
	req.ClientInfo.SigningName = signingName
	req.ClientInfo.SigningRegion = signingRegion
}

func updateS3HostForS3AccessPoint(req *request.Request) {
	updateHostPrefix(req, "s3", s3AccessPointNamespace)
}

func updateS3HostPrefixForS3ObjectLambda(req *request.Request) {
	updateHostPrefix(req, "s3", s3ObjectsLambdaNamespace)
}

func updateHostPrefix(req *request.Request, oldEndpointPrefix, newEndpointPrefix string) {
	host := req.HTTPRequest.URL.Host
	if strings.HasPrefix(host, oldEndpointPrefix) {
		// replace service hostlabel oldEndpointPrefix to newEndpointPrefix
		req.HTTPRequest.URL.Host = newEndpointPrefix + host[len(oldEndpointPrefix):]
	}
}
