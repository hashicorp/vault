package api

const (
	// ScalingPolicyTypeHorizontal indicates a policy that does horizontal scaling.
	ScalingPolicyTypeHorizontal = "horizontal"
)

// Scaling is used to query scaling-related API endpoints
type Scaling struct {
	client *Client
}

// Scaling returns a handle on the scaling endpoints.
func (c *Client) Scaling() *Scaling {
	return &Scaling{client: c}
}

func (s *Scaling) ListPolicies(q *QueryOptions) ([]*ScalingPolicyListStub, *QueryMeta, error) {
	var resp []*ScalingPolicyListStub
	qm, err := s.client.query("/v1/scaling/policies", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return resp, qm, nil
}

func (s *Scaling) GetPolicy(ID string, q *QueryOptions) (*ScalingPolicy, *QueryMeta, error) {
	var policy ScalingPolicy
	qm, err := s.client.query("/v1/scaling/policy/"+ID, &policy, q)
	if err != nil {
		return nil, nil, err
	}
	return &policy, qm, nil
}

func (p *ScalingPolicy) Canonicalize(taskGroupCount int) {
	if p.Enabled == nil {
		p.Enabled = boolToPtr(true)
	}
	if p.Min == nil {
		var m int64 = int64(taskGroupCount)
		p.Min = &m
	}
	if p.Type == "" {
		p.Type = ScalingPolicyTypeHorizontal
	}
}

// ScalingRequest is the payload for a generic scaling action
type ScalingRequest struct {
	Count   *int64
	Target  map[string]string
	Message string
	Error   bool
	Meta    map[string]interface{}
	WriteRequest
	// this is effectively a job update, so we need the ability to override policy.
	PolicyOverride bool
}

// ScalingPolicy is the user-specified API object for an autoscaling policy
type ScalingPolicy struct {
	/* fields set by user in HCL config */

	Min     *int64                 `hcl:"min,optional"`
	Max     *int64                 `hcl:"max,optional"`
	Policy  map[string]interface{} `hcl:"policy,block"`
	Enabled *bool                  `hcl:"enabled,optional"`
	Type    string                 `hcl:"type,optional"`

	/* fields set by server */

	ID          string
	Namespace   string
	Target      map[string]string
	CreateIndex uint64
	ModifyIndex uint64
}

// ScalingPolicyListStub is used to return a subset of scaling policy information
// for the scaling policy list
type ScalingPolicyListStub struct {
	ID          string
	Enabled     bool
	Type        string
	Target      map[string]string
	CreateIndex uint64
	ModifyIndex uint64
}

// JobScaleStatusResponse is used to return information about job scaling status
type JobScaleStatusResponse struct {
	JobID          string
	Namespace      string
	JobCreateIndex uint64
	JobModifyIndex uint64
	JobStopped     bool
	TaskGroups     map[string]TaskGroupScaleStatus
}

type TaskGroupScaleStatus struct {
	Desired   int
	Placed    int
	Running   int
	Healthy   int
	Unhealthy int
	Events    []ScalingEvent
}

type ScalingEvent struct {
	Count         *int64
	PreviousCount int64
	Error         bool
	Message       string
	Meta          map[string]interface{}
	EvalID        *string
	Time          uint64
	CreateIndex   uint64
}
