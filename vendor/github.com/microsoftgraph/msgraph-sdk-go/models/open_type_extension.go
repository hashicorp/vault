package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type OpenTypeExtension struct {
    Extension
}
// NewOpenTypeExtension instantiates a new OpenTypeExtension and sets the default values.
func NewOpenTypeExtension()(*OpenTypeExtension) {
    m := &OpenTypeExtension{
        Extension: *NewExtension(),
    }
    odataTypeValue := "#microsoft.graph.openTypeExtension"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateOpenTypeExtensionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOpenTypeExtensionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOpenTypeExtension(), nil
}
// GetExtensionName gets the extensionName property value. A unique text identifier for an open type data extension. Optional.
// returns a *string when successful
func (m *OpenTypeExtension) GetExtensionName()(*string) {
    val, err := m.GetBackingStore().Get("extensionName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OpenTypeExtension) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Extension.GetFieldDeserializers()
    res["extensionName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExtensionName(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *OpenTypeExtension) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Extension.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("extensionName", m.GetExtensionName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetExtensionName sets the extensionName property value. A unique text identifier for an open type data extension. Optional.
func (m *OpenTypeExtension) SetExtensionName(value *string)() {
    err := m.GetBackingStore().Set("extensionName", value)
    if err != nil {
        panic(err)
    }
}
type OpenTypeExtensionable interface {
    Extensionable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetExtensionName()(*string)
    SetExtensionName(value *string)()
}
