package wrapping

import "time"

type ResponseWrapInfo struct {
	// Setting to non-zero specifies that the response should be wrapped.
	// Specifies the desired TTL of the wrapping token.
	TTL time.Duration `json:"ttl" structs:"ttl" mapstructure:"ttl" sentinel:""`

	// The token containing the wrapped response
	Token string `json:"token" structs:"token" mapstructure:"token" sentinel:""`

	// The token accessor for the wrapped response token
	Accessor string `json:"accessor" structs:"accessor" mapstructure:"accessor"`

	// The creation time. This can be used with the TTL to figure out an
	// expected expiration.
	CreationTime time.Time `json:"creation_time" structs:"creation_time" mapstructure:"creation_time" sentinel:""`

	// If the contained response is the output of a token or approle secret-id creation call, the
	// created token's/secret-id's accessor will be accessible here
	WrappedAccessor string `json:"wrapped_accessor" structs:"wrapped_accessor" mapstructure:"wrapped_accessor" sentinel:""`

	// WrappedEntityID is the entity identifier of the caller who initiated the
	// wrapping request
	WrappedEntityID string `json:"wrapped_entity_id" structs:"wrapped_entity_id" mapstructure:"wrapped_entity_id" sentinel:""`

	// The format to use. This doesn't get returned, it's only internal.
	Format string `json:"format" structs:"format" mapstructure:"format" sentinel:""`

	// CreationPath is the original request path that was used to create
	// the wrapped response.
	CreationPath string `json:"creation_path" structs:"creation_path" mapstructure:"creation_path" sentinel:""`

	// Controls seal wrapping behavior downstream for specific use cases
	SealWrap bool `json:"seal_wrap" structs:"seal_wrap" mapstructure:"seal_wrap" sentinel:""`
}
