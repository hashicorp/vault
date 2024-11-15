package s3shared

import (
	"github.com/aws/aws-sdk-go/aws"
	awsarn "github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/internal/s3shared/arn"
)

// ResourceRequest represents the request and arn resource
type ResourceRequest struct {
	Resource arn.Resource
	Request  *request.Request
}

// ARN returns the resource ARN
func (r ResourceRequest) ARN() awsarn.ARN {
	return r.Resource.GetARN()
}

// AllowCrossRegion returns a bool value to denote if S3UseARNRegion flag is set
func (r ResourceRequest) AllowCrossRegion() bool {
	return aws.BoolValue(r.Request.Config.S3UseARNRegion)
}

// IsCrossPartition returns true if client is configured for another partition, than
// the partition that resource ARN region resolves to.
func (r ResourceRequest) IsCrossPartition() bool {
	return r.Request.ClientInfo.PartitionID != r.Resource.GetARN().Partition
}

// IsCrossRegion returns true if ARN region is different than client configured region
func (r ResourceRequest) IsCrossRegion() bool {
	return IsCrossRegion(r.Request, r.Resource.GetARN().Region)
}

// HasCustomEndpoint returns true if custom client endpoint is provided
func (r ResourceRequest) HasCustomEndpoint() bool {
	return len(aws.StringValue(r.Request.Config.Endpoint)) > 0
}

// IsCrossRegion returns true if request signing region is not same as configured region
func IsCrossRegion(req *request.Request, otherRegion string) bool {
	return req.ClientInfo.SigningRegion != otherRegion
}
