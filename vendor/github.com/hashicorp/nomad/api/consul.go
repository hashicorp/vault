package api

// Consul represents configuration related to consul.
type Consul struct {
	// (Enterprise-only) Namespace represents a Consul namespace.
	Namespace string `mapstructure:"namespace" hcl:"namespace,optional"`
}

// Canonicalize Consul into a canonical form. The Canonicalize structs containing
// a Consul should ensure it is not nil.
func (c *Consul) Canonicalize() {
	// Nothing to do here.
	//
	// If Namespace is nil, that is a choice of the job submitter that
	// we should inherit from higher up (i.e. job<-group). Likewise, if
	// Namespace is set but empty, that is a choice to use the default consul
	// namespace.
}

// Copy creates a deep copy of c.
func (c *Consul) Copy() *Consul {
	return &Consul{
		Namespace: c.Namespace,
	}
}

// MergeNamespace sets Namespace to namespace if not already configured.
// This is used to inherit the job-level consul_namespace if the group-level
// namespace is not explicitly configured.
func (c *Consul) MergeNamespace(namespace *string) {
	// only inherit namespace from above if not already set
	if c.Namespace == "" && namespace != nil {
		c.Namespace = *namespace
	}
}
