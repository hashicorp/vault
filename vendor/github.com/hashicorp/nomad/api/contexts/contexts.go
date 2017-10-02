package contexts

// Context defines the scope in which a search for Nomad object operates
type Context string

const (
	Allocs      Context = "allocs"
	Deployments Context = "deployment"
	Evals       Context = "evals"
	Jobs        Context = "jobs"
	Nodes       Context = "nodes"
	Namespaces  Context = "namespaces"
	All         Context = "all"
)
