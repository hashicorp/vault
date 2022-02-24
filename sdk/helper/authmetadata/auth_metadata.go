package authmetadata

/*
	authmetadata is a package offering convenience and
	standardization when supporting an `auth_metadata`
	field in a plugin's configuration. This then controls
	what metadata is added to an Auth during login.

	To see an example of how to add and use it, check out
	how these structs and fields are used in the AWS auth
	method.

	Or, check out its acceptance test in this package to
	see its integration points.
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// Fields is for configuring a back-end's available
// default and additional fields. These are used for
// providing a verbose field description, and for parsing
// user input.
type Fields struct {
	// The field name as it'll be reflected in the user-facing
	// schema.
	FieldName string

	// Default is a list of the default fields that should
	// be included if a user sends "default" in their list
	// of desired fields. These fields should all have a
	// low rate of change because each change can incur a
	// write to storage.
	Default []string

	// AvailableToAdd is a list of fields not included by
	// default, that the user may include.
	AvailableToAdd []string
}

func (f *Fields) all() []string {
	return append(f.Default, f.AvailableToAdd...)
}

// FieldSchema takes the default and additionally available
// fields, and uses them to generate a verbose description
// regarding how to use the "auth_metadata" field.
func FieldSchema(fields *Fields) *framework.FieldSchema {
	return &framework.FieldSchema{
		Type:        framework.TypeCommaStringSlice,
		Description: description(fields),
		DisplayAttrs: &framework.DisplayAttributes{
			Name:  fields.FieldName,
			Value: "field1,field2",
		},
		Default: []string{"default"},
	}
}

func NewHandler(fields *Fields) *Handler {
	return &Handler{
		fields: fields,
	}
}

type Handler struct {
	// authMetadata is an explicit list of all the user's configured
	// fields that are being added to auth metadata. If it is set to
	// default or unconfigured, it will be nil. Otherwise, it will
	// hold the explicit fields set by the user.
	authMetadata []string

	// fields is a list of the configured default and available
	// fields.
	fields *Fields
}

// AuthMetadata is intended to be used on config reads.
// It gets an explicit list of all the user's configured
// fields that are being added to auth metadata.
func (h *Handler) AuthMetadata() []string {
	if h.authMetadata == nil {
		return h.fields.Default
	}
	return h.authMetadata
}

// ParseAuthMetadata is intended to be used on config create/update.
// It takes a user's selected fields (or lack thereof),
// converts it to a list of explicit fields, and adds it to the Handler
// for later storage.
func (h *Handler) ParseAuthMetadata(data *framework.FieldData) error {
	userProvidedRaw, ok := data.GetOk(h.fields.FieldName)
	if !ok {
		// Nothing further to do here.
		return nil
	}
	userProvided, ok := userProvidedRaw.([]string)
	if !ok {
		return fmt.Errorf("%s is an unexpected type of %T", userProvidedRaw, userProvidedRaw)
	}
	userProvided = strutil.RemoveDuplicates(userProvided, true)

	// If the only field the user has chosen was the default field,
	// we don't store anything so we won't have to do a storage
	// migration if the default changes.
	if len(userProvided) == 1 && userProvided[0] == "default" {
		h.authMetadata = nil
		return nil
	}

	// Validate and store the input.
	if strutil.StrListContains(userProvided, "default") {
		return fmt.Errorf("%q contains default - default can't be used in combination with other fields",
			userProvided)
	}
	if !strutil.StrListSubset(h.fields.all(), userProvided) {
		return fmt.Errorf("%q contains an unavailable field, please select from %q",
			strings.Join(userProvided, ", "), strings.Join(h.fields.all(), ", "))
	}
	h.authMetadata = userProvided
	return nil
}

// PopulateDesiredMetadata is intended to be used during login
// just before returning an auth.
// It takes the available auth metadata and,
// if the auth should have it, adds it to the auth's metadata.
func (h *Handler) PopulateDesiredMetadata(auth *logical.Auth, available map[string]string) error {
	if auth == nil {
		return errors.New("auth is nil")
	}
	if auth.Metadata == nil {
		auth.Metadata = make(map[string]string)
	}
	if auth.Alias == nil {
		auth.Alias = &logical.Alias{}
	}
	if auth.Alias.Metadata == nil {
		auth.Alias.Metadata = make(map[string]string)
	}
	fieldsToInclude := h.fields.Default
	if h.authMetadata != nil {
		fieldsToInclude = h.authMetadata
	}
	for availableField, itsValue := range available {
		if itsValue == "" {
			// Don't bother setting fields for which there is no value.
			continue
		}
		if strutil.StrListContains(fieldsToInclude, availableField) {
			auth.Metadata[availableField] = itsValue
			auth.Alias.Metadata[availableField] = itsValue
		}
	}
	return nil
}

func (h *Handler) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		AuthMetadata []string `json:"auth_metadata"`
	}{
		AuthMetadata: h.authMetadata,
	})
}

func (h *Handler) UnmarshalJSON(data []byte) error {
	jsonable := &struct {
		AuthMetadata []string `json:"auth_metadata"`
	}{
		AuthMetadata: h.authMetadata,
	}
	if err := json.Unmarshal(data, jsonable); err != nil {
		return err
	}
	h.authMetadata = jsonable.AuthMetadata
	return nil
}

func description(fields *Fields) string {
	desc := "The metadata to include on the aliases and audit logs generated by this plugin."
	if len(fields.Default) > 0 {
		desc += fmt.Sprintf(" When set to 'default', includes: %s.", strings.Join(fields.Default, ", "))
	}
	if len(fields.AvailableToAdd) > 0 {
		desc += fmt.Sprintf(" These fields are available to add: %s.", strings.Join(fields.AvailableToAdd, ", "))
	}
	desc += " Not editing this field means the 'default' fields are included." +
		" Explicitly setting this field to empty overrides the 'default' and means no metadata will be included." +
		" If not using 'default', explicit fields must be sent like: 'field1,field2'."
	return desc
}
