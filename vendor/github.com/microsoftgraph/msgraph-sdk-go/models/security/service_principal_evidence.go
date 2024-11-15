package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ServicePrincipalEvidence struct {
    AlertEvidence
}
// NewServicePrincipalEvidence instantiates a new ServicePrincipalEvidence and sets the default values.
func NewServicePrincipalEvidence()(*ServicePrincipalEvidence) {
    m := &ServicePrincipalEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.servicePrincipalEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateServicePrincipalEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateServicePrincipalEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewServicePrincipalEvidence(), nil
}
// GetAppId gets the appId property value. The appId property
// returns a *string when successful
func (m *ServicePrincipalEvidence) GetAppId()(*string) {
    val, err := m.GetBackingStore().Get("appId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAppOwnerTenantId gets the appOwnerTenantId property value. The appOwnerTenantId property
// returns a *string when successful
func (m *ServicePrincipalEvidence) GetAppOwnerTenantId()(*string) {
    val, err := m.GetBackingStore().Get("appOwnerTenantId")
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
func (m *ServicePrincipalEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["appId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppId(val)
        }
        return nil
    }
    res["appOwnerTenantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppOwnerTenantId(val)
        }
        return nil
    }
    res["servicePrincipalName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServicePrincipalName(val)
        }
        return nil
    }
    res["servicePrincipalObjectId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServicePrincipalObjectId(val)
        }
        return nil
    }
    res["servicePrincipalType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseServicePrincipalType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServicePrincipalType(val.(*ServicePrincipalType))
        }
        return nil
    }
    res["tenantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTenantId(val)
        }
        return nil
    }
    return res
}
// GetServicePrincipalName gets the servicePrincipalName property value. The servicePrincipalName property
// returns a *string when successful
func (m *ServicePrincipalEvidence) GetServicePrincipalName()(*string) {
    val, err := m.GetBackingStore().Get("servicePrincipalName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetServicePrincipalObjectId gets the servicePrincipalObjectId property value. The servicePrincipalObjectId property
// returns a *string when successful
func (m *ServicePrincipalEvidence) GetServicePrincipalObjectId()(*string) {
    val, err := m.GetBackingStore().Get("servicePrincipalObjectId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetServicePrincipalType gets the servicePrincipalType property value. The servicePrincipalType property
// returns a *ServicePrincipalType when successful
func (m *ServicePrincipalEvidence) GetServicePrincipalType()(*ServicePrincipalType) {
    val, err := m.GetBackingStore().Get("servicePrincipalType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ServicePrincipalType)
    }
    return nil
}
// GetTenantId gets the tenantId property value. The tenantId property
// returns a *string when successful
func (m *ServicePrincipalEvidence) GetTenantId()(*string) {
    val, err := m.GetBackingStore().Get("tenantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ServicePrincipalEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("appId", m.GetAppId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("appOwnerTenantId", m.GetAppOwnerTenantId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("servicePrincipalName", m.GetServicePrincipalName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("servicePrincipalObjectId", m.GetServicePrincipalObjectId())
        if err != nil {
            return err
        }
    }
    if m.GetServicePrincipalType() != nil {
        cast := (*m.GetServicePrincipalType()).String()
        err = writer.WriteStringValue("servicePrincipalType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("tenantId", m.GetTenantId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppId sets the appId property value. The appId property
func (m *ServicePrincipalEvidence) SetAppId(value *string)() {
    err := m.GetBackingStore().Set("appId", value)
    if err != nil {
        panic(err)
    }
}
// SetAppOwnerTenantId sets the appOwnerTenantId property value. The appOwnerTenantId property
func (m *ServicePrincipalEvidence) SetAppOwnerTenantId(value *string)() {
    err := m.GetBackingStore().Set("appOwnerTenantId", value)
    if err != nil {
        panic(err)
    }
}
// SetServicePrincipalName sets the servicePrincipalName property value. The servicePrincipalName property
func (m *ServicePrincipalEvidence) SetServicePrincipalName(value *string)() {
    err := m.GetBackingStore().Set("servicePrincipalName", value)
    if err != nil {
        panic(err)
    }
}
// SetServicePrincipalObjectId sets the servicePrincipalObjectId property value. The servicePrincipalObjectId property
func (m *ServicePrincipalEvidence) SetServicePrincipalObjectId(value *string)() {
    err := m.GetBackingStore().Set("servicePrincipalObjectId", value)
    if err != nil {
        panic(err)
    }
}
// SetServicePrincipalType sets the servicePrincipalType property value. The servicePrincipalType property
func (m *ServicePrincipalEvidence) SetServicePrincipalType(value *ServicePrincipalType)() {
    err := m.GetBackingStore().Set("servicePrincipalType", value)
    if err != nil {
        panic(err)
    }
}
// SetTenantId sets the tenantId property value. The tenantId property
func (m *ServicePrincipalEvidence) SetTenantId(value *string)() {
    err := m.GetBackingStore().Set("tenantId", value)
    if err != nil {
        panic(err)
    }
}
type ServicePrincipalEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppId()(*string)
    GetAppOwnerTenantId()(*string)
    GetServicePrincipalName()(*string)
    GetServicePrincipalObjectId()(*string)
    GetServicePrincipalType()(*ServicePrincipalType)
    GetTenantId()(*string)
    SetAppId(value *string)()
    SetAppOwnerTenantId(value *string)()
    SetServicePrincipalName(value *string)()
    SetServicePrincipalObjectId(value *string)()
    SetServicePrincipalType(value *ServicePrincipalType)()
    SetTenantId(value *string)()
}
