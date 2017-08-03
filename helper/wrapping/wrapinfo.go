package wrapping

import "time"

type ResponseWrapInfo struct {
	// Setting to non-zero specifies that the response should be wrapped.
	// Specifies the desired TTL of the wrapping token.
	TTL time.Duration `json:"ttl" structs:"ttl" mapstructure:"ttl"`

	// The token containing the wrapped response
	Token string `json:"token" structs:"token" mapstructure:"token"`

	// The creation time. This can be used with the TTL to figure out an
	// expected expiration.
	CreationTime time.Time `json:"creation_time" structs:"creation_time" mapstructure:"creation_time"`

	// If the contained response is the output of a token creation call, the
	// created token's accessor will be accessible here
	WrappedAccessor string `json:"wrapped_accessor" structs:"wrapped_accessor" mapstructure:"wrapped_accessor"`

	// The format to use. This doesn't get returned, it's only internal.
	Format string `json:"format" structs:"format" mapstructure:"format"`

	// CreationPath is the original request path that was used to create
	// the wrapped response.
	CreationPath string `json:"creation_path" structs:"creation_path" mapstructure:"creation_path"`
}
