package capabilities

const (
	APICapability            = "api"
	MetaCapability           = "meta"
	APIPassThroughCapability = "passthrough"
	LinkControlCapability    = "link-control"
)

type Capability interface {
	Start() error
	Stop() error
}
