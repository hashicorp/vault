package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EducationOrganization struct {
    Entity
}
// NewEducationOrganization instantiates a new EducationOrganization and sets the default values.
func NewEducationOrganization()(*EducationOrganization) {
    m := &EducationOrganization{
        Entity: *NewEntity(),
    }
    return m
}
// CreateEducationOrganizationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEducationOrganizationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.educationSchool":
                        return NewEducationSchool(), nil
                }
            }
        }
    }
    return NewEducationOrganization(), nil
}
// GetDescription gets the description property value. Organization description.
// returns a *string when successful
func (m *EducationOrganization) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Organization display name.
// returns a *string when successful
func (m *EducationOrganization) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExternalSource gets the externalSource property value. Source where this organization was created from. Possible values are: sis, manual.
// returns a *EducationExternalSource when successful
func (m *EducationOrganization) GetExternalSource()(*EducationExternalSource) {
    val, err := m.GetBackingStore().Get("externalSource")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*EducationExternalSource)
    }
    return nil
}
// GetExternalSourceDetail gets the externalSourceDetail property value. The name of the external source this resource was generated from.
// returns a *string when successful
func (m *EducationOrganization) GetExternalSourceDetail()(*string) {
    val, err := m.GetBackingStore().Get("externalSourceDetail")
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
func (m *EducationOrganization) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
        }
        return nil
    }
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["externalSource"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseEducationExternalSource)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalSource(val.(*EducationExternalSource))
        }
        return nil
    }
    res["externalSourceDetail"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalSourceDetail(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *EducationOrganization) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetExternalSource() != nil {
        cast := (*m.GetExternalSource()).String()
        err = writer.WriteStringValue("externalSource", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("externalSourceDetail", m.GetExternalSourceDetail())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDescription sets the description property value. Organization description.
func (m *EducationOrganization) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Organization display name.
func (m *EducationOrganization) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalSource sets the externalSource property value. Source where this organization was created from. Possible values are: sis, manual.
func (m *EducationOrganization) SetExternalSource(value *EducationExternalSource)() {
    err := m.GetBackingStore().Set("externalSource", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalSourceDetail sets the externalSourceDetail property value. The name of the external source this resource was generated from.
func (m *EducationOrganization) SetExternalSourceDetail(value *string)() {
    err := m.GetBackingStore().Set("externalSourceDetail", value)
    if err != nil {
        panic(err)
    }
}
type EducationOrganizationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetExternalSource()(*EducationExternalSource)
    GetExternalSourceDetail()(*string)
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetExternalSource(value *EducationExternalSource)()
    SetExternalSourceDetail(value *string)()
}
