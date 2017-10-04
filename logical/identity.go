package logical

// Alias represents the information used by core to create implicit entity.
// Implicit entities get created when a client authenticates successfully from
// any of the authentication backends (except token backend).
//
// This is applicable to enterprise binaries only. Alias should be set in the
// Auth response returned by the credential backends. This structure is placed
// in the open source repository only to enable custom authetication plugins to
// be used along with enterprise binary. The custom auth plugins should make
// use of this and fill out the Alias information in the authentication
// response.
type Alias struct {
	// MountType is the backend mount's type to which this identity belongs
	// to.
	MountType string `json:"mount_type" structs:"mount_type" mapstructure:"mount_type"`

	// MountAccessor is the identifier of the mount entry to which
	// this identity
	// belongs to.
	MountAccessor string `json:"mount_accessor" structs:"mount_accessor" mapstructure:"mount_accessor"`

	// Name is the identifier of this identity in its
	// authentication source.
	Name string `json:"name" structs:"name" mapstructure:"name"`
}
