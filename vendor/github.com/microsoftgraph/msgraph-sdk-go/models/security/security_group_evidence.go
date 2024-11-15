package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SecurityGroupEvidence struct {
    AlertEvidence
}
// NewSecurityGroupEvidence instantiates a new SecurityGroupEvidence and sets the default values.
func NewSecurityGroupEvidence()(*SecurityGroupEvidence) {
    m := &SecurityGroupEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.securityGroupEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateSecurityGroupEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSecurityGroupEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSecurityGroupEvidence(), nil
}
// GetDisplayName gets the displayName property value. The name of the security group.
// returns a *string when successful
func (m *SecurityGroupEvidence) GetDisplayName()(*string) {
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
func (m *SecurityGroupEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
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
    res["securityGroupId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSecurityGroupId(val)
        }
        return nil
    }
    return res
}
// GetSecurityGroupId gets the securityGroupId property value. Unique identifier of the security group.
// returns a *string when successful
func (m *SecurityGroupEvidence) GetSecurityGroupId()(*string) {
    val, err := m.GetBackingStore().Get("securityGroupId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SecurityGroupEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("securityGroupId", m.GetSecurityGroupId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDisplayName sets the displayName property value. The name of the security group.
func (m *SecurityGroupEvidence) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetSecurityGroupId sets the securityGroupId property value. Unique identifier of the security group.
func (m *SecurityGroupEvidence) SetSecurityGroupId(value *string)() {
    err := m.GetBackingStore().Set("securityGroupId", value)
    if err != nil {
        panic(err)
    }
}
type SecurityGroupEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDisplayName()(*string)
    GetSecurityGroupId()(*string)
    SetDisplayName(value *string)()
    SetSecurityGroupId(value *string)()
}
