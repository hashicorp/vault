package logical

type LogInput struct {
	Type                string
	Auth                *Auth
	Request             *Request
	Response            *Response
	MountType           string
	OuterErr            error
	NonHMACReqDataKeys  []string
	NonHMACRespDataKeys []string
}

type MarshalOptions struct {
	ValueHasher func(string) string
}

type OptMarshaler interface {
	MarshalJSONWithOptions(*MarshalOptions) ([]byte, error)
}
