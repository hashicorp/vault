package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type OnPremisesDirectorySynchronizationFeature struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewOnPremisesDirectorySynchronizationFeature instantiates a new OnPremisesDirectorySynchronizationFeature and sets the default values.
func NewOnPremisesDirectorySynchronizationFeature()(*OnPremisesDirectorySynchronizationFeature) {
    m := &OnPremisesDirectorySynchronizationFeature{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateOnPremisesDirectorySynchronizationFeatureFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOnPremisesDirectorySynchronizationFeatureFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOnPremisesDirectorySynchronizationFeature(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetBlockCloudObjectTakeoverThroughHardMatchEnabled gets the blockCloudObjectTakeoverThroughHardMatchEnabled property value. Used to block cloud object takeover via source anchor hard match if enabled.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetBlockCloudObjectTakeoverThroughHardMatchEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("blockCloudObjectTakeoverThroughHardMatchEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBlockSoftMatchEnabled gets the blockSoftMatchEnabled property value. Use to block soft match for all objects if enabled for the  tenant. Customers are encouraged to enable this feature and keep it enabled until soft matching is required again for their tenancy. This flag should be enabled again after any soft matching has been completed and is no longer needed.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetBlockSoftMatchEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("blockSoftMatchEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBypassDirSyncOverridesEnabled gets the bypassDirSyncOverridesEnabled property value. When true, persists the values of Mobile and OtherMobile in on-premises AD during sync cycles instead of values of MobilePhone or AlternateMobilePhones in Microsoft Entra ID.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetBypassDirSyncOverridesEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("bypassDirSyncOverridesEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCloudPasswordPolicyForPasswordSyncedUsersEnabled gets the cloudPasswordPolicyForPasswordSyncedUsersEnabled property value. Used to indicate that cloud password policy applies to users whose passwords are synchronized from on-premises.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetCloudPasswordPolicyForPasswordSyncedUsersEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("cloudPasswordPolicyForPasswordSyncedUsersEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetConcurrentCredentialUpdateEnabled gets the concurrentCredentialUpdateEnabled property value. Used to enable concurrent user credentials update in OrgId.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetConcurrentCredentialUpdateEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("concurrentCredentialUpdateEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetConcurrentOrgIdProvisioningEnabled gets the concurrentOrgIdProvisioningEnabled property value. Used to enable concurrent user creation in OrgId.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetConcurrentOrgIdProvisioningEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("concurrentOrgIdProvisioningEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDeviceWritebackEnabled gets the deviceWritebackEnabled property value. Used to indicate that device write-back is enabled.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetDeviceWritebackEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("deviceWritebackEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDirectoryExtensionsEnabled gets the directoryExtensionsEnabled property value. Used to indicate that directory extensions are being synced from on-premises AD to Microsoft Entra ID.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetDirectoryExtensionsEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("directoryExtensionsEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["blockCloudObjectTakeoverThroughHardMatchEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBlockCloudObjectTakeoverThroughHardMatchEnabled(val)
        }
        return nil
    }
    res["blockSoftMatchEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBlockSoftMatchEnabled(val)
        }
        return nil
    }
    res["bypassDirSyncOverridesEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBypassDirSyncOverridesEnabled(val)
        }
        return nil
    }
    res["cloudPasswordPolicyForPasswordSyncedUsersEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCloudPasswordPolicyForPasswordSyncedUsersEnabled(val)
        }
        return nil
    }
    res["concurrentCredentialUpdateEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConcurrentCredentialUpdateEnabled(val)
        }
        return nil
    }
    res["concurrentOrgIdProvisioningEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConcurrentOrgIdProvisioningEnabled(val)
        }
        return nil
    }
    res["deviceWritebackEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceWritebackEnabled(val)
        }
        return nil
    }
    res["directoryExtensionsEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDirectoryExtensionsEnabled(val)
        }
        return nil
    }
    res["fopeConflictResolutionEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFopeConflictResolutionEnabled(val)
        }
        return nil
    }
    res["groupWriteBackEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGroupWriteBackEnabled(val)
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    res["passwordSyncEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordSyncEnabled(val)
        }
        return nil
    }
    res["passwordWritebackEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordWritebackEnabled(val)
        }
        return nil
    }
    res["quarantineUponProxyAddressesConflictEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQuarantineUponProxyAddressesConflictEnabled(val)
        }
        return nil
    }
    res["quarantineUponUpnConflictEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQuarantineUponUpnConflictEnabled(val)
        }
        return nil
    }
    res["softMatchOnUpnEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSoftMatchOnUpnEnabled(val)
        }
        return nil
    }
    res["synchronizeUpnForManagedUsersEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSynchronizeUpnForManagedUsersEnabled(val)
        }
        return nil
    }
    res["unifiedGroupWritebackEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUnifiedGroupWritebackEnabled(val)
        }
        return nil
    }
    res["userForcePasswordChangeOnLogonEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserForcePasswordChangeOnLogonEnabled(val)
        }
        return nil
    }
    res["userWritebackEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserWritebackEnabled(val)
        }
        return nil
    }
    return res
}
// GetFopeConflictResolutionEnabled gets the fopeConflictResolutionEnabled property value. Used to indicate that for a Microsoft Forefront Online Protection for Exchange (FOPE) migrated tenant, the conflicting proxy address should be migrated over.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetFopeConflictResolutionEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("fopeConflictResolutionEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetGroupWriteBackEnabled gets the groupWriteBackEnabled property value. Used to enable object-level group writeback feature for additional group types.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetGroupWriteBackEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("groupWriteBackEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPasswordSyncEnabled gets the passwordSyncEnabled property value. Used to indicate on-premise password synchronization is enabled.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetPasswordSyncEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("passwordSyncEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasswordWritebackEnabled gets the passwordWritebackEnabled property value. Used to indicate that writeback of password resets from Microsoft Entra ID to on-premises AD is enabled. This property isn't in use and updating it isn't supported.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetPasswordWritebackEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("passwordWritebackEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetQuarantineUponProxyAddressesConflictEnabled gets the quarantineUponProxyAddressesConflictEnabled property value. Used to indicate that we should quarantine objects with conflicting proxy address.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetQuarantineUponProxyAddressesConflictEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("quarantineUponProxyAddressesConflictEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetQuarantineUponUpnConflictEnabled gets the quarantineUponUpnConflictEnabled property value. Used to indicate that we should quarantine objects conflicting with duplicate userPrincipalName.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetQuarantineUponUpnConflictEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("quarantineUponUpnConflictEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSoftMatchOnUpnEnabled gets the softMatchOnUpnEnabled property value. Used to indicate that we should soft match objects based on userPrincipalName.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetSoftMatchOnUpnEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("softMatchOnUpnEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSynchronizeUpnForManagedUsersEnabled gets the synchronizeUpnForManagedUsersEnabled property value. Used to indicate that we should synchronize userPrincipalName objects for managed users with licenses.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetSynchronizeUpnForManagedUsersEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("synchronizeUpnForManagedUsersEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetUnifiedGroupWritebackEnabled gets the unifiedGroupWritebackEnabled property value. Used to indicate that Microsoft 365 Group write-back is enabled.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetUnifiedGroupWritebackEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("unifiedGroupWritebackEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetUserForcePasswordChangeOnLogonEnabled gets the userForcePasswordChangeOnLogonEnabled property value. Used to indicate that feature to force password change for a user on logon is enabled while synchronizing on-premise credentials.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetUserForcePasswordChangeOnLogonEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("userForcePasswordChangeOnLogonEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetUserWritebackEnabled gets the userWritebackEnabled property value. Used to indicate that user writeback is enabled.
// returns a *bool when successful
func (m *OnPremisesDirectorySynchronizationFeature) GetUserWritebackEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("userWritebackEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *OnPremisesDirectorySynchronizationFeature) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("blockCloudObjectTakeoverThroughHardMatchEnabled", m.GetBlockCloudObjectTakeoverThroughHardMatchEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("blockSoftMatchEnabled", m.GetBlockSoftMatchEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("bypassDirSyncOverridesEnabled", m.GetBypassDirSyncOverridesEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("cloudPasswordPolicyForPasswordSyncedUsersEnabled", m.GetCloudPasswordPolicyForPasswordSyncedUsersEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("concurrentCredentialUpdateEnabled", m.GetConcurrentCredentialUpdateEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("concurrentOrgIdProvisioningEnabled", m.GetConcurrentOrgIdProvisioningEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("deviceWritebackEnabled", m.GetDeviceWritebackEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("directoryExtensionsEnabled", m.GetDirectoryExtensionsEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("fopeConflictResolutionEnabled", m.GetFopeConflictResolutionEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("groupWriteBackEnabled", m.GetGroupWriteBackEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("passwordSyncEnabled", m.GetPasswordSyncEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("passwordWritebackEnabled", m.GetPasswordWritebackEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("quarantineUponProxyAddressesConflictEnabled", m.GetQuarantineUponProxyAddressesConflictEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("quarantineUponUpnConflictEnabled", m.GetQuarantineUponUpnConflictEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("softMatchOnUpnEnabled", m.GetSoftMatchOnUpnEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("synchronizeUpnForManagedUsersEnabled", m.GetSynchronizeUpnForManagedUsersEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("unifiedGroupWritebackEnabled", m.GetUnifiedGroupWritebackEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("userForcePasswordChangeOnLogonEnabled", m.GetUserForcePasswordChangeOnLogonEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("userWritebackEnabled", m.GetUserWritebackEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *OnPremisesDirectorySynchronizationFeature) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *OnPremisesDirectorySynchronizationFeature) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetBlockCloudObjectTakeoverThroughHardMatchEnabled sets the blockCloudObjectTakeoverThroughHardMatchEnabled property value. Used to block cloud object takeover via source anchor hard match if enabled.
func (m *OnPremisesDirectorySynchronizationFeature) SetBlockCloudObjectTakeoverThroughHardMatchEnabled(value *bool)() {
    err := m.GetBackingStore().Set("blockCloudObjectTakeoverThroughHardMatchEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetBlockSoftMatchEnabled sets the blockSoftMatchEnabled property value. Use to block soft match for all objects if enabled for the  tenant. Customers are encouraged to enable this feature and keep it enabled until soft matching is required again for their tenancy. This flag should be enabled again after any soft matching has been completed and is no longer needed.
func (m *OnPremisesDirectorySynchronizationFeature) SetBlockSoftMatchEnabled(value *bool)() {
    err := m.GetBackingStore().Set("blockSoftMatchEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetBypassDirSyncOverridesEnabled sets the bypassDirSyncOverridesEnabled property value. When true, persists the values of Mobile and OtherMobile in on-premises AD during sync cycles instead of values of MobilePhone or AlternateMobilePhones in Microsoft Entra ID.
func (m *OnPremisesDirectorySynchronizationFeature) SetBypassDirSyncOverridesEnabled(value *bool)() {
    err := m.GetBackingStore().Set("bypassDirSyncOverridesEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetCloudPasswordPolicyForPasswordSyncedUsersEnabled sets the cloudPasswordPolicyForPasswordSyncedUsersEnabled property value. Used to indicate that cloud password policy applies to users whose passwords are synchronized from on-premises.
func (m *OnPremisesDirectorySynchronizationFeature) SetCloudPasswordPolicyForPasswordSyncedUsersEnabled(value *bool)() {
    err := m.GetBackingStore().Set("cloudPasswordPolicyForPasswordSyncedUsersEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetConcurrentCredentialUpdateEnabled sets the concurrentCredentialUpdateEnabled property value. Used to enable concurrent user credentials update in OrgId.
func (m *OnPremisesDirectorySynchronizationFeature) SetConcurrentCredentialUpdateEnabled(value *bool)() {
    err := m.GetBackingStore().Set("concurrentCredentialUpdateEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetConcurrentOrgIdProvisioningEnabled sets the concurrentOrgIdProvisioningEnabled property value. Used to enable concurrent user creation in OrgId.
func (m *OnPremisesDirectorySynchronizationFeature) SetConcurrentOrgIdProvisioningEnabled(value *bool)() {
    err := m.GetBackingStore().Set("concurrentOrgIdProvisioningEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceWritebackEnabled sets the deviceWritebackEnabled property value. Used to indicate that device write-back is enabled.
func (m *OnPremisesDirectorySynchronizationFeature) SetDeviceWritebackEnabled(value *bool)() {
    err := m.GetBackingStore().Set("deviceWritebackEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetDirectoryExtensionsEnabled sets the directoryExtensionsEnabled property value. Used to indicate that directory extensions are being synced from on-premises AD to Microsoft Entra ID.
func (m *OnPremisesDirectorySynchronizationFeature) SetDirectoryExtensionsEnabled(value *bool)() {
    err := m.GetBackingStore().Set("directoryExtensionsEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetFopeConflictResolutionEnabled sets the fopeConflictResolutionEnabled property value. Used to indicate that for a Microsoft Forefront Online Protection for Exchange (FOPE) migrated tenant, the conflicting proxy address should be migrated over.
func (m *OnPremisesDirectorySynchronizationFeature) SetFopeConflictResolutionEnabled(value *bool)() {
    err := m.GetBackingStore().Set("fopeConflictResolutionEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetGroupWriteBackEnabled sets the groupWriteBackEnabled property value. Used to enable object-level group writeback feature for additional group types.
func (m *OnPremisesDirectorySynchronizationFeature) SetGroupWriteBackEnabled(value *bool)() {
    err := m.GetBackingStore().Set("groupWriteBackEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *OnPremisesDirectorySynchronizationFeature) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordSyncEnabled sets the passwordSyncEnabled property value. Used to indicate on-premise password synchronization is enabled.
func (m *OnPremisesDirectorySynchronizationFeature) SetPasswordSyncEnabled(value *bool)() {
    err := m.GetBackingStore().Set("passwordSyncEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordWritebackEnabled sets the passwordWritebackEnabled property value. Used to indicate that writeback of password resets from Microsoft Entra ID to on-premises AD is enabled. This property isn't in use and updating it isn't supported.
func (m *OnPremisesDirectorySynchronizationFeature) SetPasswordWritebackEnabled(value *bool)() {
    err := m.GetBackingStore().Set("passwordWritebackEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetQuarantineUponProxyAddressesConflictEnabled sets the quarantineUponProxyAddressesConflictEnabled property value. Used to indicate that we should quarantine objects with conflicting proxy address.
func (m *OnPremisesDirectorySynchronizationFeature) SetQuarantineUponProxyAddressesConflictEnabled(value *bool)() {
    err := m.GetBackingStore().Set("quarantineUponProxyAddressesConflictEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetQuarantineUponUpnConflictEnabled sets the quarantineUponUpnConflictEnabled property value. Used to indicate that we should quarantine objects conflicting with duplicate userPrincipalName.
func (m *OnPremisesDirectorySynchronizationFeature) SetQuarantineUponUpnConflictEnabled(value *bool)() {
    err := m.GetBackingStore().Set("quarantineUponUpnConflictEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetSoftMatchOnUpnEnabled sets the softMatchOnUpnEnabled property value. Used to indicate that we should soft match objects based on userPrincipalName.
func (m *OnPremisesDirectorySynchronizationFeature) SetSoftMatchOnUpnEnabled(value *bool)() {
    err := m.GetBackingStore().Set("softMatchOnUpnEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetSynchronizeUpnForManagedUsersEnabled sets the synchronizeUpnForManagedUsersEnabled property value. Used to indicate that we should synchronize userPrincipalName objects for managed users with licenses.
func (m *OnPremisesDirectorySynchronizationFeature) SetSynchronizeUpnForManagedUsersEnabled(value *bool)() {
    err := m.GetBackingStore().Set("synchronizeUpnForManagedUsersEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetUnifiedGroupWritebackEnabled sets the unifiedGroupWritebackEnabled property value. Used to indicate that Microsoft 365 Group write-back is enabled.
func (m *OnPremisesDirectorySynchronizationFeature) SetUnifiedGroupWritebackEnabled(value *bool)() {
    err := m.GetBackingStore().Set("unifiedGroupWritebackEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetUserForcePasswordChangeOnLogonEnabled sets the userForcePasswordChangeOnLogonEnabled property value. Used to indicate that feature to force password change for a user on logon is enabled while synchronizing on-premise credentials.
func (m *OnPremisesDirectorySynchronizationFeature) SetUserForcePasswordChangeOnLogonEnabled(value *bool)() {
    err := m.GetBackingStore().Set("userForcePasswordChangeOnLogonEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetUserWritebackEnabled sets the userWritebackEnabled property value. Used to indicate that user writeback is enabled.
func (m *OnPremisesDirectorySynchronizationFeature) SetUserWritebackEnabled(value *bool)() {
    err := m.GetBackingStore().Set("userWritebackEnabled", value)
    if err != nil {
        panic(err)
    }
}
type OnPremisesDirectorySynchronizationFeatureable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetBlockCloudObjectTakeoverThroughHardMatchEnabled()(*bool)
    GetBlockSoftMatchEnabled()(*bool)
    GetBypassDirSyncOverridesEnabled()(*bool)
    GetCloudPasswordPolicyForPasswordSyncedUsersEnabled()(*bool)
    GetConcurrentCredentialUpdateEnabled()(*bool)
    GetConcurrentOrgIdProvisioningEnabled()(*bool)
    GetDeviceWritebackEnabled()(*bool)
    GetDirectoryExtensionsEnabled()(*bool)
    GetFopeConflictResolutionEnabled()(*bool)
    GetGroupWriteBackEnabled()(*bool)
    GetOdataType()(*string)
    GetPasswordSyncEnabled()(*bool)
    GetPasswordWritebackEnabled()(*bool)
    GetQuarantineUponProxyAddressesConflictEnabled()(*bool)
    GetQuarantineUponUpnConflictEnabled()(*bool)
    GetSoftMatchOnUpnEnabled()(*bool)
    GetSynchronizeUpnForManagedUsersEnabled()(*bool)
    GetUnifiedGroupWritebackEnabled()(*bool)
    GetUserForcePasswordChangeOnLogonEnabled()(*bool)
    GetUserWritebackEnabled()(*bool)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetBlockCloudObjectTakeoverThroughHardMatchEnabled(value *bool)()
    SetBlockSoftMatchEnabled(value *bool)()
    SetBypassDirSyncOverridesEnabled(value *bool)()
    SetCloudPasswordPolicyForPasswordSyncedUsersEnabled(value *bool)()
    SetConcurrentCredentialUpdateEnabled(value *bool)()
    SetConcurrentOrgIdProvisioningEnabled(value *bool)()
    SetDeviceWritebackEnabled(value *bool)()
    SetDirectoryExtensionsEnabled(value *bool)()
    SetFopeConflictResolutionEnabled(value *bool)()
    SetGroupWriteBackEnabled(value *bool)()
    SetOdataType(value *string)()
    SetPasswordSyncEnabled(value *bool)()
    SetPasswordWritebackEnabled(value *bool)()
    SetQuarantineUponProxyAddressesConflictEnabled(value *bool)()
    SetQuarantineUponUpnConflictEnabled(value *bool)()
    SetSoftMatchOnUpnEnabled(value *bool)()
    SetSynchronizeUpnForManagedUsersEnabled(value *bool)()
    SetUnifiedGroupWritebackEnabled(value *bool)()
    SetUserForcePasswordChangeOnLogonEnabled(value *bool)()
    SetUserWritebackEnabled(value *bool)()
}
