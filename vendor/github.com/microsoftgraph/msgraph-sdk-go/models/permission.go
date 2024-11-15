package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Permission struct {
    Entity
}
// NewPermission instantiates a new Permission and sets the default values.
func NewPermission()(*Permission) {
    m := &Permission{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePermissionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePermissionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPermission(), nil
}
// GetExpirationDateTime gets the expirationDateTime property value. A format of yyyy-MM-ddTHH:mm:ssZ of DateTimeOffset indicates the expiration time of the permission. DateTime.MinValue indicates there's no expiration set for this permission. Optional.
// returns a *Time when successful
func (m *Permission) GetExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("expirationDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Permission) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["expirationDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExpirationDateTime(val)
        }
        return nil
    }
    res["grantedTo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGrantedTo(val.(IdentitySetable))
        }
        return nil
    }
    res["grantedToIdentities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]IdentitySetable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(IdentitySetable)
                }
            }
            m.SetGrantedToIdentities(res)
        }
        return nil
    }
    res["grantedToIdentitiesV2"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSharePointIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SharePointIdentitySetable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SharePointIdentitySetable)
                }
            }
            m.SetGrantedToIdentitiesV2(res)
        }
        return nil
    }
    res["grantedToV2"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSharePointIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGrantedToV2(val.(SharePointIdentitySetable))
        }
        return nil
    }
    res["hasPassword"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHasPassword(val)
        }
        return nil
    }
    res["inheritedFrom"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemReferenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInheritedFrom(val.(ItemReferenceable))
        }
        return nil
    }
    res["invitation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSharingInvitationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInvitation(val.(SharingInvitationable))
        }
        return nil
    }
    res["link"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSharingLinkFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLink(val.(SharingLinkable))
        }
        return nil
    }
    res["roles"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetRoles(res)
        }
        return nil
    }
    res["shareId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShareId(val)
        }
        return nil
    }
    return res
}
// GetGrantedTo gets the grantedTo property value. For user type permissions, the details of the users and applications for this permission. Read-only.
// returns a IdentitySetable when successful
func (m *Permission) GetGrantedTo()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("grantedTo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetGrantedToIdentities gets the grantedToIdentities property value. For type permissions, the details of the users to whom permission was granted. Read-only.
// returns a []IdentitySetable when successful
func (m *Permission) GetGrantedToIdentities()([]IdentitySetable) {
    val, err := m.GetBackingStore().Get("grantedToIdentities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]IdentitySetable)
    }
    return nil
}
// GetGrantedToIdentitiesV2 gets the grantedToIdentitiesV2 property value. For link type permissions, the details of the users to whom permission was granted. Read-only.
// returns a []SharePointIdentitySetable when successful
func (m *Permission) GetGrantedToIdentitiesV2()([]SharePointIdentitySetable) {
    val, err := m.GetBackingStore().Get("grantedToIdentitiesV2")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SharePointIdentitySetable)
    }
    return nil
}
// GetGrantedToV2 gets the grantedToV2 property value. For user type permissions, the details of the users and applications for this permission. Read-only.
// returns a SharePointIdentitySetable when successful
func (m *Permission) GetGrantedToV2()(SharePointIdentitySetable) {
    val, err := m.GetBackingStore().Get("grantedToV2")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SharePointIdentitySetable)
    }
    return nil
}
// GetHasPassword gets the hasPassword property value. Indicates whether the password is set for this permission. This property only appears in the response. Optional. Read-only. For OneDrive Personal only..
// returns a *bool when successful
func (m *Permission) GetHasPassword()(*bool) {
    val, err := m.GetBackingStore().Get("hasPassword")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetInheritedFrom gets the inheritedFrom property value. Provides a reference to the ancestor of the current permission, if it's inherited from an ancestor. Read-only.
// returns a ItemReferenceable when successful
func (m *Permission) GetInheritedFrom()(ItemReferenceable) {
    val, err := m.GetBackingStore().Get("inheritedFrom")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemReferenceable)
    }
    return nil
}
// GetInvitation gets the invitation property value. Details of any associated sharing invitation for this permission. Read-only.
// returns a SharingInvitationable when successful
func (m *Permission) GetInvitation()(SharingInvitationable) {
    val, err := m.GetBackingStore().Get("invitation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SharingInvitationable)
    }
    return nil
}
// GetLink gets the link property value. Provides the link details of the current permission, if it's a link type permission. Read-only.
// returns a SharingLinkable when successful
func (m *Permission) GetLink()(SharingLinkable) {
    val, err := m.GetBackingStore().Get("link")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SharingLinkable)
    }
    return nil
}
// GetRoles gets the roles property value. The type of permission, for example, read. See below for the full list of roles. Read-only.
// returns a []string when successful
func (m *Permission) GetRoles()([]string) {
    val, err := m.GetBackingStore().Get("roles")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetShareId gets the shareId property value. A unique token that can be used to access this shared item via the shares API. Read-only.
// returns a *string when successful
func (m *Permission) GetShareId()(*string) {
    val, err := m.GetBackingStore().Get("shareId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Permission) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("expirationDateTime", m.GetExpirationDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("grantedTo", m.GetGrantedTo())
        if err != nil {
            return err
        }
    }
    if m.GetGrantedToIdentities() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetGrantedToIdentities()))
        for i, v := range m.GetGrantedToIdentities() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("grantedToIdentities", cast)
        if err != nil {
            return err
        }
    }
    if m.GetGrantedToIdentitiesV2() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetGrantedToIdentitiesV2()))
        for i, v := range m.GetGrantedToIdentitiesV2() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("grantedToIdentitiesV2", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("grantedToV2", m.GetGrantedToV2())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("hasPassword", m.GetHasPassword())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("inheritedFrom", m.GetInheritedFrom())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("invitation", m.GetInvitation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("link", m.GetLink())
        if err != nil {
            return err
        }
    }
    if m.GetRoles() != nil {
        err = writer.WriteCollectionOfStringValues("roles", m.GetRoles())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("shareId", m.GetShareId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetExpirationDateTime sets the expirationDateTime property value. A format of yyyy-MM-ddTHH:mm:ssZ of DateTimeOffset indicates the expiration time of the permission. DateTime.MinValue indicates there's no expiration set for this permission. Optional.
func (m *Permission) SetExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("expirationDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetGrantedTo sets the grantedTo property value. For user type permissions, the details of the users and applications for this permission. Read-only.
func (m *Permission) SetGrantedTo(value IdentitySetable)() {
    err := m.GetBackingStore().Set("grantedTo", value)
    if err != nil {
        panic(err)
    }
}
// SetGrantedToIdentities sets the grantedToIdentities property value. For type permissions, the details of the users to whom permission was granted. Read-only.
func (m *Permission) SetGrantedToIdentities(value []IdentitySetable)() {
    err := m.GetBackingStore().Set("grantedToIdentities", value)
    if err != nil {
        panic(err)
    }
}
// SetGrantedToIdentitiesV2 sets the grantedToIdentitiesV2 property value. For link type permissions, the details of the users to whom permission was granted. Read-only.
func (m *Permission) SetGrantedToIdentitiesV2(value []SharePointIdentitySetable)() {
    err := m.GetBackingStore().Set("grantedToIdentitiesV2", value)
    if err != nil {
        panic(err)
    }
}
// SetGrantedToV2 sets the grantedToV2 property value. For user type permissions, the details of the users and applications for this permission. Read-only.
func (m *Permission) SetGrantedToV2(value SharePointIdentitySetable)() {
    err := m.GetBackingStore().Set("grantedToV2", value)
    if err != nil {
        panic(err)
    }
}
// SetHasPassword sets the hasPassword property value. Indicates whether the password is set for this permission. This property only appears in the response. Optional. Read-only. For OneDrive Personal only..
func (m *Permission) SetHasPassword(value *bool)() {
    err := m.GetBackingStore().Set("hasPassword", value)
    if err != nil {
        panic(err)
    }
}
// SetInheritedFrom sets the inheritedFrom property value. Provides a reference to the ancestor of the current permission, if it's inherited from an ancestor. Read-only.
func (m *Permission) SetInheritedFrom(value ItemReferenceable)() {
    err := m.GetBackingStore().Set("inheritedFrom", value)
    if err != nil {
        panic(err)
    }
}
// SetInvitation sets the invitation property value. Details of any associated sharing invitation for this permission. Read-only.
func (m *Permission) SetInvitation(value SharingInvitationable)() {
    err := m.GetBackingStore().Set("invitation", value)
    if err != nil {
        panic(err)
    }
}
// SetLink sets the link property value. Provides the link details of the current permission, if it's a link type permission. Read-only.
func (m *Permission) SetLink(value SharingLinkable)() {
    err := m.GetBackingStore().Set("link", value)
    if err != nil {
        panic(err)
    }
}
// SetRoles sets the roles property value. The type of permission, for example, read. See below for the full list of roles. Read-only.
func (m *Permission) SetRoles(value []string)() {
    err := m.GetBackingStore().Set("roles", value)
    if err != nil {
        panic(err)
    }
}
// SetShareId sets the shareId property value. A unique token that can be used to access this shared item via the shares API. Read-only.
func (m *Permission) SetShareId(value *string)() {
    err := m.GetBackingStore().Set("shareId", value)
    if err != nil {
        panic(err)
    }
}
type Permissionable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetGrantedTo()(IdentitySetable)
    GetGrantedToIdentities()([]IdentitySetable)
    GetGrantedToIdentitiesV2()([]SharePointIdentitySetable)
    GetGrantedToV2()(SharePointIdentitySetable)
    GetHasPassword()(*bool)
    GetInheritedFrom()(ItemReferenceable)
    GetInvitation()(SharingInvitationable)
    GetLink()(SharingLinkable)
    GetRoles()([]string)
    GetShareId()(*string)
    SetExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetGrantedTo(value IdentitySetable)()
    SetGrantedToIdentities(value []IdentitySetable)()
    SetGrantedToIdentitiesV2(value []SharePointIdentitySetable)()
    SetGrantedToV2(value SharePointIdentitySetable)()
    SetHasPassword(value *bool)()
    SetInheritedFrom(value ItemReferenceable)()
    SetInvitation(value SharingInvitationable)()
    SetLink(value SharingLinkable)()
    SetRoles(value []string)()
    SetShareId(value *string)()
}
