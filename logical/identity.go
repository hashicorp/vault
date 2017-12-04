package logical

// Alias represents the information used by core to create implicit entity.
// Implicit entities get created when a client authenticates successfully from
// any of the authentication backends (except token backend).
type Alias struct {
	// MountType is the backend mount's type to which this identity belongs
	MountType string `json:"mount_type" structs:"mount_type" mapstructure:"mount_type"`

	// MountAccessor is the identifier of the mount entry to which this
	// identity belongs
	MountAccessor string `json:"mount_accessor" structs:"mount_accessor" mapstructure:"mount_accessor"`

	// Name is the identifier of this identity in its authentication source
	Name string `json:"name" structs:"name" mapstructure:"name"`
}
