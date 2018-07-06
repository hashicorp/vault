package structs

import (
	"errors"
	"fmt"
	"strings"
)

const (
	errNoLeader            = "No cluster leader"
	errNoRegionPath        = "No path to region"
	errTokenNotFound       = "ACL token not found"
	errPermissionDenied    = "Permission denied"
	errNoNodeConn          = "No path to node"
	errUnknownMethod       = "Unknown rpc method"
	errUnknownNomadVersion = "Unable to determine Nomad version"
	errNodeLacksRpc        = "Node does not support RPC; requires 0.8 or later"

	// Prefix based errors that are used to check if the error is of a given
	// type. These errors should be created with the associated constructor.
	ErrUnknownAllocationPrefix = "Unknown allocation"
	ErrUnknownNodePrefix       = "Unknown node"
	ErrUnknownJobPrefix        = "Unknown job"
	ErrUnknownEvaluationPrefix = "Unknown evaluation"
	ErrUnknownDeploymentPrefix = "Unknown deployment"
)

var (
	ErrNoLeader            = errors.New(errNoLeader)
	ErrNoRegionPath        = errors.New(errNoRegionPath)
	ErrTokenNotFound       = errors.New(errTokenNotFound)
	ErrPermissionDenied    = errors.New(errPermissionDenied)
	ErrNoNodeConn          = errors.New(errNoNodeConn)
	ErrUnknownMethod       = errors.New(errUnknownMethod)
	ErrUnknownNomadVersion = errors.New(errUnknownNomadVersion)
	ErrNodeLacksRpc        = errors.New(errNodeLacksRpc)
)

// IsErrNoLeader returns whether the error is due to there being no leader.
func IsErrNoLeader(err error) bool {
	return err != nil && strings.Contains(err.Error(), errNoLeader)
}

// IsErrNoRegionPath returns whether the error is due to there being no path to
// the given region.
func IsErrNoRegionPath(err error) bool {
	return err != nil && strings.Contains(err.Error(), errNoRegionPath)
}

// IsErrTokenNotFound returns whether the error is due to the passed token not
// being resolvable.
func IsErrTokenNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), errTokenNotFound)
}

// IsErrPermissionDenied returns whether the error is due to the operation not
// being allowed due to lack of permissions.
func IsErrPermissionDenied(err error) bool {
	return err != nil && strings.Contains(err.Error(), errPermissionDenied)
}

// IsErrNoNodeConn returns whether the error is due to there being no path to
// the given node.
func IsErrNoNodeConn(err error) bool {
	return err != nil && strings.Contains(err.Error(), errNoNodeConn)
}

// IsErrUnknownMethod returns whether the error is due to the operation not
// being allowed due to lack of permissions.
func IsErrUnknownMethod(err error) bool {
	return err != nil && strings.Contains(err.Error(), errUnknownMethod)
}

// NewErrUnknownAllocation returns a new error caused by the allocation being
// unknown.
func NewErrUnknownAllocation(allocID string) error {
	return fmt.Errorf("%s %q", ErrUnknownAllocationPrefix, allocID)
}

// NewErrUnknownNode returns a new error caused by the node being unknown.
func NewErrUnknownNode(nodeID string) error {
	return fmt.Errorf("%s %q", ErrUnknownNodePrefix, nodeID)
}

// NewErrUnknownJob returns a new error caused by the job being unknown.
func NewErrUnknownJob(jobID string) error {
	return fmt.Errorf("%s %q", ErrUnknownJobPrefix, jobID)
}

// NewErrUnknownEvaluation returns a new error caused by the evaluation being
// unknown.
func NewErrUnknownEvaluation(evaluationID string) error {
	return fmt.Errorf("%s %q", ErrUnknownEvaluationPrefix, evaluationID)
}

// NewErrUnknownDeployment returns a new error caused by the deployment being
// unknown.
func NewErrUnknownDeployment(deploymentID string) error {
	return fmt.Errorf("%s %q", ErrUnknownDeploymentPrefix, deploymentID)
}

// IsErrUnknownAllocation returns whether the error is due to an unknown
// allocation.
func IsErrUnknownAllocation(err error) bool {
	return err != nil && strings.Contains(err.Error(), ErrUnknownAllocationPrefix)
}

// IsErrUnknownNode returns whether the error is due to an unknown
// node.
func IsErrUnknownNode(err error) bool {
	return err != nil && strings.Contains(err.Error(), ErrUnknownNodePrefix)
}

// IsErrUnknownJob returns whether the error is due to an unknown
// job.
func IsErrUnknownJob(err error) bool {
	return err != nil && strings.Contains(err.Error(), ErrUnknownJobPrefix)
}

// IsErrUnknownEvaluation returns whether the error is due to an unknown
// evaluation.
func IsErrUnknownEvaluation(err error) bool {
	return err != nil && strings.Contains(err.Error(), ErrUnknownEvaluationPrefix)
}

// IsErrUnknownDeployment returns whether the error is due to an unknown
// deployment.
func IsErrUnknownDeployment(err error) bool {
	return err != nil && strings.Contains(err.Error(), ErrUnknownDeploymentPrefix)
}

// IsErrUnknownNomadVersion returns whether the error is due to Nomad being
// unable to determine the version of a node.
func IsErrUnknownNomadVersion(err error) bool {
	return err != nil && strings.Contains(err.Error(), errUnknownNomadVersion)
}

// IsErrNodeLacksRpc returns whether error is due to a Nomad server being
// unable to connect to a client node because the client is too old (pre-v0.8).
func IsErrNodeLacksRpc(err error) bool {
	return err != nil && strings.Contains(err.Error(), errNodeLacksRpc)
}
