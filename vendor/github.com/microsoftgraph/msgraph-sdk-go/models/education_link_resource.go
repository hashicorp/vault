package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EducationLinkResource struct {
    EducationResource
}
// NewEducationLinkResource instantiates a new EducationLinkResource and sets the default values.
func NewEducationLinkResource()(*EducationLinkResource) {
    m := &EducationLinkResource{
        EducationResource: *NewEducationResource(),
    }
    odataTypeValue := "#microsoft.graph.educationLinkResource"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEducationLinkResourceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEducationLinkResourceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEducationLinkResource(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EducationLinkResource) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EducationResource.GetFieldDeserializers()
    res["link"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLink(val)
        }
        return nil
    }
    return res
}
// GetLink gets the link property value. URL to the resource.
// returns a *string when successful
func (m *EducationLinkResource) GetLink()(*string) {
    val, err := m.GetBackingStore().Get("link")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EducationLinkResource) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EducationResource.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("link", m.GetLink())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetLink sets the link property value. URL to the resource.
func (m *EducationLinkResource) SetLink(value *string)() {
    err := m.GetBackingStore().Set("link", value)
    if err != nil {
        panic(err)
    }
}
type EducationLinkResourceable interface {
    EducationResourceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetLink()(*string)
    SetLink(value *string)()
}
