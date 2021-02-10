package contexts

// Context defines the scope in which a search for Nomad object operates
type Context string

const (
	Allocs          Context = "allocs"
	Deployments     Context = "deployment"
	Evals           Context = "evals"
	Jobs            Context = "jobs"
	Nodes           Context = "nodes"
	Namespaces      Context = "namespaces"
	Quotas          Context = "quotas"
	Recommendations Context = "recommendations"
	ScalingPolicies Context = "scaling_policy"
	Plugins         Context = "plugins"
	Volumes         Context = "volumes"
	All             Context = "all"
)
