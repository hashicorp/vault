package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AuditLogRoot struct {
    Entity
}
// NewAuditLogRoot instantiates a new AuditLogRoot and sets the default values.
func NewAuditLogRoot()(*AuditLogRoot) {
    m := &AuditLogRoot{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAuditLogRootFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAuditLogRootFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAuditLogRoot(), nil
}
// GetDirectoryAudits gets the directoryAudits property value. The directoryAudits property
// returns a []DirectoryAuditable when successful
func (m *AuditLogRoot) GetDirectoryAudits()([]DirectoryAuditable) {
    val, err := m.GetBackingStore().Get("directoryAudits")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryAuditable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AuditLogRoot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["directoryAudits"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryAuditFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryAuditable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DirectoryAuditable)
                }
            }
            m.SetDirectoryAudits(res)
        }
        return nil
    }
    res["provisioning"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateProvisioningObjectSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ProvisioningObjectSummaryable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ProvisioningObjectSummaryable)
                }
            }
            m.SetProvisioning(res)
        }
        return nil
    }
    res["signIns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSignInFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SignInable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SignInable)
                }
            }
            m.SetSignIns(res)
        }
        return nil
    }
    return res
}
// GetProvisioning gets the provisioning property value. The provisioning property
// returns a []ProvisioningObjectSummaryable when successful
func (m *AuditLogRoot) GetProvisioning()([]ProvisioningObjectSummaryable) {
    val, err := m.GetBackingStore().Get("provisioning")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ProvisioningObjectSummaryable)
    }
    return nil
}
// GetSignIns gets the signIns property value. The signIns property
// returns a []SignInable when successful
func (m *AuditLogRoot) GetSignIns()([]SignInable) {
    val, err := m.GetBackingStore().Get("signIns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SignInable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AuditLogRoot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetDirectoryAudits() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDirectoryAudits()))
        for i, v := range m.GetDirectoryAudits() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("directoryAudits", cast)
        if err != nil {
            return err
        }
    }
    if m.GetProvisioning() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetProvisioning()))
        for i, v := range m.GetProvisioning() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("provisioning", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSignIns() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSignIns()))
        for i, v := range m.GetSignIns() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("signIns", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDirectoryAudits sets the directoryAudits property value. The directoryAudits property
func (m *AuditLogRoot) SetDirectoryAudits(value []DirectoryAuditable)() {
    err := m.GetBackingStore().Set("directoryAudits", value)
    if err != nil {
        panic(err)
    }
}
// SetProvisioning sets the provisioning property value. The provisioning property
func (m *AuditLogRoot) SetProvisioning(value []ProvisioningObjectSummaryable)() {
    err := m.GetBackingStore().Set("provisioning", value)
    if err != nil {
        panic(err)
    }
}
// SetSignIns sets the signIns property value. The signIns property
func (m *AuditLogRoot) SetSignIns(value []SignInable)() {
    err := m.GetBackingStore().Set("signIns", value)
    if err != nil {
        panic(err)
    }
}
type AuditLogRootable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDirectoryAudits()([]DirectoryAuditable)
    GetProvisioning()([]ProvisioningObjectSummaryable)
    GetSignIns()([]SignInable)
    SetDirectoryAudits(value []DirectoryAuditable)()
    SetProvisioning(value []ProvisioningObjectSummaryable)()
    SetSignIns(value []SignInable)()
}
