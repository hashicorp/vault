package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type MembershipOutlierInsight struct {
    GovernanceInsight
}
// NewMembershipOutlierInsight instantiates a new MembershipOutlierInsight and sets the default values.
func NewMembershipOutlierInsight()(*MembershipOutlierInsight) {
    m := &MembershipOutlierInsight{
        GovernanceInsight: *NewGovernanceInsight(),
    }
    odataTypeValue := "#microsoft.graph.membershipOutlierInsight"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateMembershipOutlierInsightFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMembershipOutlierInsightFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMembershipOutlierInsight(), nil
}
// GetContainer gets the container property value. Navigation link to the container directory object. For example, to a group.
// returns a DirectoryObjectable when successful
func (m *MembershipOutlierInsight) GetContainer()(DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("container")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DirectoryObjectable)
    }
    return nil
}
// GetContainerId gets the containerId property value. Indicates the identifier of the container, for example, a group ID.
// returns a *string when successful
func (m *MembershipOutlierInsight) GetContainerId()(*string) {
    val, err := m.GetBackingStore().Get("containerId")
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
func (m *MembershipOutlierInsight) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.GovernanceInsight.GetFieldDeserializers()
    res["container"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContainer(val.(DirectoryObjectable))
        }
        return nil
    }
    res["containerId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContainerId(val)
        }
        return nil
    }
    res["lastModifiedBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedBy(val.(Userable))
        }
        return nil
    }
    res["member"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMember(val.(DirectoryObjectable))
        }
        return nil
    }
    res["memberId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMemberId(val)
        }
        return nil
    }
    res["outlierContainerType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseOutlierContainerType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOutlierContainerType(val.(*OutlierContainerType))
        }
        return nil
    }
    res["outlierMemberType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseOutlierMemberType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOutlierMemberType(val.(*OutlierMemberType))
        }
        return nil
    }
    return res
}
// GetLastModifiedBy gets the lastModifiedBy property value. Navigation link to a member object who modified the record. For example, to a user.
// returns a Userable when successful
func (m *MembershipOutlierInsight) GetLastModifiedBy()(Userable) {
    val, err := m.GetBackingStore().Get("lastModifiedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Userable)
    }
    return nil
}
// GetMember gets the member property value. Navigation link to a member object. For example, to a user.
// returns a DirectoryObjectable when successful
func (m *MembershipOutlierInsight) GetMember()(DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("member")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DirectoryObjectable)
    }
    return nil
}
// GetMemberId gets the memberId property value. Indicates the identifier of the user.
// returns a *string when successful
func (m *MembershipOutlierInsight) GetMemberId()(*string) {
    val, err := m.GetBackingStore().Get("memberId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOutlierContainerType gets the outlierContainerType property value. The outlierContainerType property
// returns a *OutlierContainerType when successful
func (m *MembershipOutlierInsight) GetOutlierContainerType()(*OutlierContainerType) {
    val, err := m.GetBackingStore().Get("outlierContainerType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*OutlierContainerType)
    }
    return nil
}
// GetOutlierMemberType gets the outlierMemberType property value. The outlierMemberType property
// returns a *OutlierMemberType when successful
func (m *MembershipOutlierInsight) GetOutlierMemberType()(*OutlierMemberType) {
    val, err := m.GetBackingStore().Get("outlierMemberType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*OutlierMemberType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MembershipOutlierInsight) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.GovernanceInsight.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("container", m.GetContainer())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("containerId", m.GetContainerId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("lastModifiedBy", m.GetLastModifiedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("member", m.GetMember())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("memberId", m.GetMemberId())
        if err != nil {
            return err
        }
    }
    if m.GetOutlierContainerType() != nil {
        cast := (*m.GetOutlierContainerType()).String()
        err = writer.WriteStringValue("outlierContainerType", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetOutlierMemberType() != nil {
        cast := (*m.GetOutlierMemberType()).String()
        err = writer.WriteStringValue("outlierMemberType", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetContainer sets the container property value. Navigation link to the container directory object. For example, to a group.
func (m *MembershipOutlierInsight) SetContainer(value DirectoryObjectable)() {
    err := m.GetBackingStore().Set("container", value)
    if err != nil {
        panic(err)
    }
}
// SetContainerId sets the containerId property value. Indicates the identifier of the container, for example, a group ID.
func (m *MembershipOutlierInsight) SetContainerId(value *string)() {
    err := m.GetBackingStore().Set("containerId", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedBy sets the lastModifiedBy property value. Navigation link to a member object who modified the record. For example, to a user.
func (m *MembershipOutlierInsight) SetLastModifiedBy(value Userable)() {
    err := m.GetBackingStore().Set("lastModifiedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetMember sets the member property value. Navigation link to a member object. For example, to a user.
func (m *MembershipOutlierInsight) SetMember(value DirectoryObjectable)() {
    err := m.GetBackingStore().Set("member", value)
    if err != nil {
        panic(err)
    }
}
// SetMemberId sets the memberId property value. Indicates the identifier of the user.
func (m *MembershipOutlierInsight) SetMemberId(value *string)() {
    err := m.GetBackingStore().Set("memberId", value)
    if err != nil {
        panic(err)
    }
}
// SetOutlierContainerType sets the outlierContainerType property value. The outlierContainerType property
func (m *MembershipOutlierInsight) SetOutlierContainerType(value *OutlierContainerType)() {
    err := m.GetBackingStore().Set("outlierContainerType", value)
    if err != nil {
        panic(err)
    }
}
// SetOutlierMemberType sets the outlierMemberType property value. The outlierMemberType property
func (m *MembershipOutlierInsight) SetOutlierMemberType(value *OutlierMemberType)() {
    err := m.GetBackingStore().Set("outlierMemberType", value)
    if err != nil {
        panic(err)
    }
}
type MembershipOutlierInsightable interface {
    GovernanceInsightable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetContainer()(DirectoryObjectable)
    GetContainerId()(*string)
    GetLastModifiedBy()(Userable)
    GetMember()(DirectoryObjectable)
    GetMemberId()(*string)
    GetOutlierContainerType()(*OutlierContainerType)
    GetOutlierMemberType()(*OutlierMemberType)
    SetContainer(value DirectoryObjectable)()
    SetContainerId(value *string)()
    SetLastModifiedBy(value Userable)()
    SetMember(value DirectoryObjectable)()
    SetMemberId(value *string)()
    SetOutlierContainerType(value *OutlierContainerType)()
    SetOutlierMemberType(value *OutlierMemberType)()
}
