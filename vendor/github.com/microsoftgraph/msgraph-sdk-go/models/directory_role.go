package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DirectoryRole struct {
    DirectoryObject
}
// NewDirectoryRole instantiates a new DirectoryRole and sets the default values.
func NewDirectoryRole()(*DirectoryRole) {
    m := &DirectoryRole{
        DirectoryObject: *NewDirectoryObject(),
    }
    odataTypeValue := "#microsoft.graph.directoryRole"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateDirectoryRoleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDirectoryRoleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDirectoryRole(), nil
}
// GetDescription gets the description property value. The description for the directory role. Read-only. Supports $filter (eq), $search, $select.
// returns a *string when successful
func (m *DirectoryRole) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The display name for the directory role. Read-only. Supports $filter (eq), $search, $select.
// returns a *string when successful
func (m *DirectoryRole) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
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
func (m *DirectoryRole) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DirectoryObject.GetFieldDeserializers()
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
    res["members"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DirectoryObjectable)
                }
            }
            m.SetMembers(res)
        }
        return nil
    }
    res["roleTemplateId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRoleTemplateId(val)
        }
        return nil
    }
    res["scopedMembers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateScopedRoleMembershipFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ScopedRoleMembershipable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ScopedRoleMembershipable)
                }
            }
            m.SetScopedMembers(res)
        }
        return nil
    }
    return res
}
// GetMembers gets the members property value. Users that are members of this directory role. HTTP Methods: GET, POST, DELETE. Read-only. Nullable. Supports $expand.
// returns a []DirectoryObjectable when successful
func (m *DirectoryRole) GetMembers()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("members")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetRoleTemplateId gets the roleTemplateId property value. The id of the directoryRoleTemplate that this role is based on. The property must be specified when activating a directory role in a tenant with a POST operation. After the directory role has been activated, the property is read only. Supports $filter (eq), $select.
// returns a *string when successful
func (m *DirectoryRole) GetRoleTemplateId()(*string) {
    val, err := m.GetBackingStore().Get("roleTemplateId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetScopedMembers gets the scopedMembers property value. Members of this directory role that are scoped to administrative units. Read-only. Nullable.
// returns a []ScopedRoleMembershipable when successful
func (m *DirectoryRole) GetScopedMembers()([]ScopedRoleMembershipable) {
    val, err := m.GetBackingStore().Get("scopedMembers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ScopedRoleMembershipable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DirectoryRole) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DirectoryObject.Serialize(writer)
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
    if m.GetMembers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMembers()))
        for i, v := range m.GetMembers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("members", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("roleTemplateId", m.GetRoleTemplateId())
        if err != nil {
            return err
        }
    }
    if m.GetScopedMembers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetScopedMembers()))
        for i, v := range m.GetScopedMembers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("scopedMembers", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDescription sets the description property value. The description for the directory role. Read-only. Supports $filter (eq), $search, $select.
func (m *DirectoryRole) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The display name for the directory role. Read-only. Supports $filter (eq), $search, $select.
func (m *DirectoryRole) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetMembers sets the members property value. Users that are members of this directory role. HTTP Methods: GET, POST, DELETE. Read-only. Nullable. Supports $expand.
func (m *DirectoryRole) SetMembers(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("members", value)
    if err != nil {
        panic(err)
    }
}
// SetRoleTemplateId sets the roleTemplateId property value. The id of the directoryRoleTemplate that this role is based on. The property must be specified when activating a directory role in a tenant with a POST operation. After the directory role has been activated, the property is read only. Supports $filter (eq), $select.
func (m *DirectoryRole) SetRoleTemplateId(value *string)() {
    err := m.GetBackingStore().Set("roleTemplateId", value)
    if err != nil {
        panic(err)
    }
}
// SetScopedMembers sets the scopedMembers property value. Members of this directory role that are scoped to administrative units. Read-only. Nullable.
func (m *DirectoryRole) SetScopedMembers(value []ScopedRoleMembershipable)() {
    err := m.GetBackingStore().Set("scopedMembers", value)
    if err != nil {
        panic(err)
    }
}
type DirectoryRoleable interface {
    DirectoryObjectable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetMembers()([]DirectoryObjectable)
    GetRoleTemplateId()(*string)
    GetScopedMembers()([]ScopedRoleMembershipable)
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetMembers(value []DirectoryObjectable)()
    SetRoleTemplateId(value *string)()
    SetScopedMembers(value []ScopedRoleMembershipable)()
}
