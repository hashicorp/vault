package profile

type HttpProfile struct {
	ReqMethod  string
	ReqTimeout int
	Endpoint   string
	// Deprecated, use Scheme instead
	Protocol string
	Scheme   string
}

func NewHttpProfile() *HttpProfile {
	return &HttpProfile{
		ReqMethod:  "POST",
		ReqTimeout: 60,
		Endpoint:   "",
		Scheme:     "HTTPS",
	}
}
