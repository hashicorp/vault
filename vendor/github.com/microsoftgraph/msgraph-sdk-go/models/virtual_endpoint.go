package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type VirtualEndpoint struct {
    Entity
}
// NewVirtualEndpoint instantiates a new VirtualEndpoint and sets the default values.
func NewVirtualEndpoint()(*VirtualEndpoint) {
    m := &VirtualEndpoint{
        Entity: *NewEntity(),
    }
    return m
}
// CreateVirtualEndpointFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateVirtualEndpointFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewVirtualEndpoint(), nil
}
// GetAuditEvents gets the auditEvents property value. A collection of Cloud PC audit events.
// returns a []CloudPcAuditEventable when successful
func (m *VirtualEndpoint) GetAuditEvents()([]CloudPcAuditEventable) {
    val, err := m.GetBackingStore().Get("auditEvents")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CloudPcAuditEventable)
    }
    return nil
}
// GetCloudPCs gets the cloudPCs property value. A collection of cloud-managed virtual desktops.
// returns a []CloudPCable when successful
func (m *VirtualEndpoint) GetCloudPCs()([]CloudPCable) {
    val, err := m.GetBackingStore().Get("cloudPCs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CloudPCable)
    }
    return nil
}
// GetDeviceImages gets the deviceImages property value. A collection of device image resources on Cloud PC.
// returns a []CloudPcDeviceImageable when successful
func (m *VirtualEndpoint) GetDeviceImages()([]CloudPcDeviceImageable) {
    val, err := m.GetBackingStore().Get("deviceImages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CloudPcDeviceImageable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *VirtualEndpoint) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["auditEvents"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCloudPcAuditEventFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CloudPcAuditEventable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CloudPcAuditEventable)
                }
            }
            m.SetAuditEvents(res)
        }
        return nil
    }
    res["cloudPCs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCloudPCFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CloudPCable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CloudPCable)
                }
            }
            m.SetCloudPCs(res)
        }
        return nil
    }
    res["deviceImages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCloudPcDeviceImageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CloudPcDeviceImageable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CloudPcDeviceImageable)
                }
            }
            m.SetDeviceImages(res)
        }
        return nil
    }
    res["galleryImages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCloudPcGalleryImageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CloudPcGalleryImageable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CloudPcGalleryImageable)
                }
            }
            m.SetGalleryImages(res)
        }
        return nil
    }
    res["onPremisesConnections"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCloudPcOnPremisesConnectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CloudPcOnPremisesConnectionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CloudPcOnPremisesConnectionable)
                }
            }
            m.SetOnPremisesConnections(res)
        }
        return nil
    }
    res["provisioningPolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCloudPcProvisioningPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CloudPcProvisioningPolicyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CloudPcProvisioningPolicyable)
                }
            }
            m.SetProvisioningPolicies(res)
        }
        return nil
    }
    res["userSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCloudPcUserSettingFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CloudPcUserSettingable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CloudPcUserSettingable)
                }
            }
            m.SetUserSettings(res)
        }
        return nil
    }
    return res
}
// GetGalleryImages gets the galleryImages property value. A collection of gallery image resources on Cloud PC.
// returns a []CloudPcGalleryImageable when successful
func (m *VirtualEndpoint) GetGalleryImages()([]CloudPcGalleryImageable) {
    val, err := m.GetBackingStore().Get("galleryImages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CloudPcGalleryImageable)
    }
    return nil
}
// GetOnPremisesConnections gets the onPremisesConnections property value. A defined collection of Azure resource information that can be used to establish Azure network connections for Cloud PCs.
// returns a []CloudPcOnPremisesConnectionable when successful
func (m *VirtualEndpoint) GetOnPremisesConnections()([]CloudPcOnPremisesConnectionable) {
    val, err := m.GetBackingStore().Get("onPremisesConnections")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CloudPcOnPremisesConnectionable)
    }
    return nil
}
// GetProvisioningPolicies gets the provisioningPolicies property value. A collection of Cloud PC provisioning policies.
// returns a []CloudPcProvisioningPolicyable when successful
func (m *VirtualEndpoint) GetProvisioningPolicies()([]CloudPcProvisioningPolicyable) {
    val, err := m.GetBackingStore().Get("provisioningPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CloudPcProvisioningPolicyable)
    }
    return nil
}
// GetUserSettings gets the userSettings property value. A collection of Cloud PC user settings.
// returns a []CloudPcUserSettingable when successful
func (m *VirtualEndpoint) GetUserSettings()([]CloudPcUserSettingable) {
    val, err := m.GetBackingStore().Get("userSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CloudPcUserSettingable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *VirtualEndpoint) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAuditEvents() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAuditEvents()))
        for i, v := range m.GetAuditEvents() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("auditEvents", cast)
        if err != nil {
            return err
        }
    }
    if m.GetCloudPCs() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCloudPCs()))
        for i, v := range m.GetCloudPCs() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("cloudPCs", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDeviceImages() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDeviceImages()))
        for i, v := range m.GetDeviceImages() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("deviceImages", cast)
        if err != nil {
            return err
        }
    }
    if m.GetGalleryImages() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetGalleryImages()))
        for i, v := range m.GetGalleryImages() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("galleryImages", cast)
        if err != nil {
            return err
        }
    }
    if m.GetOnPremisesConnections() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOnPremisesConnections()))
        for i, v := range m.GetOnPremisesConnections() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("onPremisesConnections", cast)
        if err != nil {
            return err
        }
    }
    if m.GetProvisioningPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetProvisioningPolicies()))
        for i, v := range m.GetProvisioningPolicies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("provisioningPolicies", cast)
        if err != nil {
            return err
        }
    }
    if m.GetUserSettings() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetUserSettings()))
        for i, v := range m.GetUserSettings() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("userSettings", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAuditEvents sets the auditEvents property value. A collection of Cloud PC audit events.
func (m *VirtualEndpoint) SetAuditEvents(value []CloudPcAuditEventable)() {
    err := m.GetBackingStore().Set("auditEvents", value)
    if err != nil {
        panic(err)
    }
}
// SetCloudPCs sets the cloudPCs property value. A collection of cloud-managed virtual desktops.
func (m *VirtualEndpoint) SetCloudPCs(value []CloudPCable)() {
    err := m.GetBackingStore().Set("cloudPCs", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceImages sets the deviceImages property value. A collection of device image resources on Cloud PC.
func (m *VirtualEndpoint) SetDeviceImages(value []CloudPcDeviceImageable)() {
    err := m.GetBackingStore().Set("deviceImages", value)
    if err != nil {
        panic(err)
    }
}
// SetGalleryImages sets the galleryImages property value. A collection of gallery image resources on Cloud PC.
func (m *VirtualEndpoint) SetGalleryImages(value []CloudPcGalleryImageable)() {
    err := m.GetBackingStore().Set("galleryImages", value)
    if err != nil {
        panic(err)
    }
}
// SetOnPremisesConnections sets the onPremisesConnections property value. A defined collection of Azure resource information that can be used to establish Azure network connections for Cloud PCs.
func (m *VirtualEndpoint) SetOnPremisesConnections(value []CloudPcOnPremisesConnectionable)() {
    err := m.GetBackingStore().Set("onPremisesConnections", value)
    if err != nil {
        panic(err)
    }
}
// SetProvisioningPolicies sets the provisioningPolicies property value. A collection of Cloud PC provisioning policies.
func (m *VirtualEndpoint) SetProvisioningPolicies(value []CloudPcProvisioningPolicyable)() {
    err := m.GetBackingStore().Set("provisioningPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetUserSettings sets the userSettings property value. A collection of Cloud PC user settings.
func (m *VirtualEndpoint) SetUserSettings(value []CloudPcUserSettingable)() {
    err := m.GetBackingStore().Set("userSettings", value)
    if err != nil {
        panic(err)
    }
}
type VirtualEndpointable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAuditEvents()([]CloudPcAuditEventable)
    GetCloudPCs()([]CloudPCable)
    GetDeviceImages()([]CloudPcDeviceImageable)
    GetGalleryImages()([]CloudPcGalleryImageable)
    GetOnPremisesConnections()([]CloudPcOnPremisesConnectionable)
    GetProvisioningPolicies()([]CloudPcProvisioningPolicyable)
    GetUserSettings()([]CloudPcUserSettingable)
    SetAuditEvents(value []CloudPcAuditEventable)()
    SetCloudPCs(value []CloudPCable)()
    SetDeviceImages(value []CloudPcDeviceImageable)()
    SetGalleryImages(value []CloudPcGalleryImageable)()
    SetOnPremisesConnections(value []CloudPcOnPremisesConnectionable)()
    SetProvisioningPolicies(value []CloudPcProvisioningPolicyable)()
    SetUserSettings(value []CloudPcUserSettingable)()
}
