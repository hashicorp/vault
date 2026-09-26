module github.com/hashicorp/vault/version

// The go version directive value isn't consulted when building our production binaries,
// and the vault module isn't intended to be imported into other projects.  As such the
// impact of this setting is usually rather limited.  Note however that in some cases the
// Go project introduces new semantics for handling of go.mod depending on the value.
//
// The general policy for updating it is: when the Go major version used on the branch is
// updated. If we choose not to do so at some point (e.g. because we don't want some new
// semantic related to Go module handling), this comment should be updated to explain that.
//
// Whenever this value gets updated, sdk/go.mod should be updated to the same value.
go 1.27.1

require github.com/stretchr/testify v1.12.1

require go.yaml.in/yaml/v3 v3.0.5 // indirect
