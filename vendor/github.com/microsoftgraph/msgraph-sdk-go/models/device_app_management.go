package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceAppManagement singleton entity that acts as a container for all device app management functionality.
type DeviceAppManagement struct {
    Entity
}
// NewDeviceAppManagement instantiates a new DeviceAppManagement and sets the default values.
func NewDeviceAppManagement()(*DeviceAppManagement) {
    m := &DeviceAppManagement{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceAppManagementFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceAppManagementFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceAppManagement(), nil
}
// GetAndroidManagedAppProtections gets the androidManagedAppProtections property value. Android managed app policies.
// returns a []AndroidManagedAppProtectionable when successful
func (m *DeviceAppManagement) GetAndroidManagedAppProtections()([]AndroidManagedAppProtectionable) {
    val, err := m.GetBackingStore().Get("androidManagedAppProtections")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AndroidManagedAppProtectionable)
    }
    return nil
}
// GetDefaultManagedAppProtections gets the defaultManagedAppProtections property value. Default managed app policies.
// returns a []DefaultManagedAppProtectionable when successful
func (m *DeviceAppManagement) GetDefaultManagedAppProtections()([]DefaultManagedAppProtectionable) {
    val, err := m.GetBackingStore().Get("defaultManagedAppProtections")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DefaultManagedAppProtectionable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DeviceAppManagement) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["androidManagedAppProtections"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAndroidManagedAppProtectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AndroidManagedAppProtectionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AndroidManagedAppProtectionable)
                }
            }
            m.SetAndroidManagedAppProtections(res)
        }
        return nil
    }
    res["defaultManagedAppProtections"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDefaultManagedAppProtectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DefaultManagedAppProtectionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DefaultManagedAppProtectionable)
                }
            }
            m.SetDefaultManagedAppProtections(res)
        }
        return nil
    }
    res["iosManagedAppProtections"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateIosManagedAppProtectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]IosManagedAppProtectionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(IosManagedAppProtectionable)
                }
            }
            m.SetIosManagedAppProtections(res)
        }
        return nil
    }
    res["isEnabledForMicrosoftStoreForBusiness"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsEnabledForMicrosoftStoreForBusiness(val)
        }
        return nil
    }
    res["managedAppPolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateManagedAppPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ManagedAppPolicyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ManagedAppPolicyable)
                }
            }
            m.SetManagedAppPolicies(res)
        }
        return nil
    }
    res["managedAppRegistrations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateManagedAppRegistrationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ManagedAppRegistrationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ManagedAppRegistrationable)
                }
            }
            m.SetManagedAppRegistrations(res)
        }
        return nil
    }
    res["managedAppStatuses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateManagedAppStatusFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ManagedAppStatusable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ManagedAppStatusable)
                }
            }
            m.SetManagedAppStatuses(res)
        }
        return nil
    }
    res["managedEBooks"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateManagedEBookFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ManagedEBookable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ManagedEBookable)
                }
            }
            m.SetManagedEBooks(res)
        }
        return nil
    }
    res["mdmWindowsInformationProtectionPolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMdmWindowsInformationProtectionPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MdmWindowsInformationProtectionPolicyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MdmWindowsInformationProtectionPolicyable)
                }
            }
            m.SetMdmWindowsInformationProtectionPolicies(res)
        }
        return nil
    }
    res["microsoftStoreForBusinessLanguage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMicrosoftStoreForBusinessLanguage(val)
        }
        return nil
    }
    res["microsoftStoreForBusinessLastCompletedApplicationSyncTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMicrosoftStoreForBusinessLastCompletedApplicationSyncTime(val)
        }
        return nil
    }
    res["microsoftStoreForBusinessLastSuccessfulSyncDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMicrosoftStoreForBusinessLastSuccessfulSyncDateTime(val)
        }
        return nil
    }
    res["mobileAppCategories"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMobileAppCategoryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MobileAppCategoryable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MobileAppCategoryable)
                }
            }
            m.SetMobileAppCategories(res)
        }
        return nil
    }
    res["mobileAppConfigurations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateManagedDeviceMobileAppConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ManagedDeviceMobileAppConfigurationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ManagedDeviceMobileAppConfigurationable)
                }
            }
            m.SetMobileAppConfigurations(res)
        }
        return nil
    }
    res["mobileApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMobileAppFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MobileAppable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MobileAppable)
                }
            }
            m.SetMobileApps(res)
        }
        return nil
    }
    res["targetedManagedAppConfigurations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTargetedManagedAppConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TargetedManagedAppConfigurationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TargetedManagedAppConfigurationable)
                }
            }
            m.SetTargetedManagedAppConfigurations(res)
        }
        return nil
    }
    res["vppTokens"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateVppTokenFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]VppTokenable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(VppTokenable)
                }
            }
            m.SetVppTokens(res)
        }
        return nil
    }
    res["windowsInformationProtectionPolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsInformationProtectionPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsInformationProtectionPolicyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WindowsInformationProtectionPolicyable)
                }
            }
            m.SetWindowsInformationProtectionPolicies(res)
        }
        return nil
    }
    return res
}
// GetIosManagedAppProtections gets the iosManagedAppProtections property value. iOS managed app policies.
// returns a []IosManagedAppProtectionable when successful
func (m *DeviceAppManagement) GetIosManagedAppProtections()([]IosManagedAppProtectionable) {
    val, err := m.GetBackingStore().Get("iosManagedAppProtections")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]IosManagedAppProtectionable)
    }
    return nil
}
// GetIsEnabledForMicrosoftStoreForBusiness gets the isEnabledForMicrosoftStoreForBusiness property value. Whether the account is enabled for syncing applications from the Microsoft Store for Business.
// returns a *bool when successful
func (m *DeviceAppManagement) GetIsEnabledForMicrosoftStoreForBusiness()(*bool) {
    val, err := m.GetBackingStore().Get("isEnabledForMicrosoftStoreForBusiness")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetManagedAppPolicies gets the managedAppPolicies property value. Managed app policies.
// returns a []ManagedAppPolicyable when successful
func (m *DeviceAppManagement) GetManagedAppPolicies()([]ManagedAppPolicyable) {
    val, err := m.GetBackingStore().Get("managedAppPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ManagedAppPolicyable)
    }
    return nil
}
// GetManagedAppRegistrations gets the managedAppRegistrations property value. The managed app registrations.
// returns a []ManagedAppRegistrationable when successful
func (m *DeviceAppManagement) GetManagedAppRegistrations()([]ManagedAppRegistrationable) {
    val, err := m.GetBackingStore().Get("managedAppRegistrations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ManagedAppRegistrationable)
    }
    return nil
}
// GetManagedAppStatuses gets the managedAppStatuses property value. The managed app statuses.
// returns a []ManagedAppStatusable when successful
func (m *DeviceAppManagement) GetManagedAppStatuses()([]ManagedAppStatusable) {
    val, err := m.GetBackingStore().Get("managedAppStatuses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ManagedAppStatusable)
    }
    return nil
}
// GetManagedEBooks gets the managedEBooks property value. The Managed eBook.
// returns a []ManagedEBookable when successful
func (m *DeviceAppManagement) GetManagedEBooks()([]ManagedEBookable) {
    val, err := m.GetBackingStore().Get("managedEBooks")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ManagedEBookable)
    }
    return nil
}
// GetMdmWindowsInformationProtectionPolicies gets the mdmWindowsInformationProtectionPolicies property value. Windows information protection for apps running on devices which are MDM enrolled.
// returns a []MdmWindowsInformationProtectionPolicyable when successful
func (m *DeviceAppManagement) GetMdmWindowsInformationProtectionPolicies()([]MdmWindowsInformationProtectionPolicyable) {
    val, err := m.GetBackingStore().Get("mdmWindowsInformationProtectionPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MdmWindowsInformationProtectionPolicyable)
    }
    return nil
}
// GetMicrosoftStoreForBusinessLanguage gets the microsoftStoreForBusinessLanguage property value. The locale information used to sync applications from the Microsoft Store for Business. Cultures that are specific to a country/region. The names of these cultures follow RFC 4646 (Windows Vista and later). The format is -<country/regioncode2>, where  is a lowercase two-letter code derived from ISO 639-1 and <country/regioncode2> is an uppercase two-letter code derived from ISO 3166. For example, en-US for English (United States) is a specific culture.
// returns a *string when successful
func (m *DeviceAppManagement) GetMicrosoftStoreForBusinessLanguage()(*string) {
    val, err := m.GetBackingStore().Get("microsoftStoreForBusinessLanguage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMicrosoftStoreForBusinessLastCompletedApplicationSyncTime gets the microsoftStoreForBusinessLastCompletedApplicationSyncTime property value. The last time an application sync from the Microsoft Store for Business was completed.
// returns a *Time when successful
func (m *DeviceAppManagement) GetMicrosoftStoreForBusinessLastCompletedApplicationSyncTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("microsoftStoreForBusinessLastCompletedApplicationSyncTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetMicrosoftStoreForBusinessLastSuccessfulSyncDateTime gets the microsoftStoreForBusinessLastSuccessfulSyncDateTime property value. The last time the apps from the Microsoft Store for Business were synced successfully for the account.
// returns a *Time when successful
func (m *DeviceAppManagement) GetMicrosoftStoreForBusinessLastSuccessfulSyncDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("microsoftStoreForBusinessLastSuccessfulSyncDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetMobileAppCategories gets the mobileAppCategories property value. The mobile app categories.
// returns a []MobileAppCategoryable when successful
func (m *DeviceAppManagement) GetMobileAppCategories()([]MobileAppCategoryable) {
    val, err := m.GetBackingStore().Get("mobileAppCategories")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MobileAppCategoryable)
    }
    return nil
}
// GetMobileAppConfigurations gets the mobileAppConfigurations property value. The Managed Device Mobile Application Configurations.
// returns a []ManagedDeviceMobileAppConfigurationable when successful
func (m *DeviceAppManagement) GetMobileAppConfigurations()([]ManagedDeviceMobileAppConfigurationable) {
    val, err := m.GetBackingStore().Get("mobileAppConfigurations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ManagedDeviceMobileAppConfigurationable)
    }
    return nil
}
// GetMobileApps gets the mobileApps property value. The mobile apps.
// returns a []MobileAppable when successful
func (m *DeviceAppManagement) GetMobileApps()([]MobileAppable) {
    val, err := m.GetBackingStore().Get("mobileApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MobileAppable)
    }
    return nil
}
// GetTargetedManagedAppConfigurations gets the targetedManagedAppConfigurations property value. Targeted managed app configurations.
// returns a []TargetedManagedAppConfigurationable when successful
func (m *DeviceAppManagement) GetTargetedManagedAppConfigurations()([]TargetedManagedAppConfigurationable) {
    val, err := m.GetBackingStore().Get("targetedManagedAppConfigurations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TargetedManagedAppConfigurationable)
    }
    return nil
}
// GetVppTokens gets the vppTokens property value. List of Vpp tokens for this organization.
// returns a []VppTokenable when successful
func (m *DeviceAppManagement) GetVppTokens()([]VppTokenable) {
    val, err := m.GetBackingStore().Get("vppTokens")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]VppTokenable)
    }
    return nil
}
// GetWindowsInformationProtectionPolicies gets the windowsInformationProtectionPolicies property value. Windows information protection for apps running on devices which are not MDM enrolled.
// returns a []WindowsInformationProtectionPolicyable when successful
func (m *DeviceAppManagement) GetWindowsInformationProtectionPolicies()([]WindowsInformationProtectionPolicyable) {
    val, err := m.GetBackingStore().Get("windowsInformationProtectionPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WindowsInformationProtectionPolicyable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceAppManagement) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAndroidManagedAppProtections() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAndroidManagedAppProtections()))
        for i, v := range m.GetAndroidManagedAppProtections() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("androidManagedAppProtections", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDefaultManagedAppProtections() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDefaultManagedAppProtections()))
        for i, v := range m.GetDefaultManagedAppProtections() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("defaultManagedAppProtections", cast)
        if err != nil {
            return err
        }
    }
    if m.GetIosManagedAppProtections() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetIosManagedAppProtections()))
        for i, v := range m.GetIosManagedAppProtections() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("iosManagedAppProtections", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isEnabledForMicrosoftStoreForBusiness", m.GetIsEnabledForMicrosoftStoreForBusiness())
        if err != nil {
            return err
        }
    }
    if m.GetManagedAppPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetManagedAppPolicies()))
        for i, v := range m.GetManagedAppPolicies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("managedAppPolicies", cast)
        if err != nil {
            return err
        }
    }
    if m.GetManagedAppRegistrations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetManagedAppRegistrations()))
        for i, v := range m.GetManagedAppRegistrations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("managedAppRegistrations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetManagedAppStatuses() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetManagedAppStatuses()))
        for i, v := range m.GetManagedAppStatuses() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("managedAppStatuses", cast)
        if err != nil {
            return err
        }
    }
    if m.GetManagedEBooks() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetManagedEBooks()))
        for i, v := range m.GetManagedEBooks() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("managedEBooks", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMdmWindowsInformationProtectionPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMdmWindowsInformationProtectionPolicies()))
        for i, v := range m.GetMdmWindowsInformationProtectionPolicies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("mdmWindowsInformationProtectionPolicies", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("microsoftStoreForBusinessLanguage", m.GetMicrosoftStoreForBusinessLanguage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("microsoftStoreForBusinessLastCompletedApplicationSyncTime", m.GetMicrosoftStoreForBusinessLastCompletedApplicationSyncTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("microsoftStoreForBusinessLastSuccessfulSyncDateTime", m.GetMicrosoftStoreForBusinessLastSuccessfulSyncDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetMobileAppCategories() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMobileAppCategories()))
        for i, v := range m.GetMobileAppCategories() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("mobileAppCategories", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMobileAppConfigurations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMobileAppConfigurations()))
        for i, v := range m.GetMobileAppConfigurations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("mobileAppConfigurations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMobileApps() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMobileApps()))
        for i, v := range m.GetMobileApps() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("mobileApps", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTargetedManagedAppConfigurations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTargetedManagedAppConfigurations()))
        for i, v := range m.GetTargetedManagedAppConfigurations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("targetedManagedAppConfigurations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetVppTokens() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetVppTokens()))
        for i, v := range m.GetVppTokens() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("vppTokens", cast)
        if err != nil {
            return err
        }
    }
    if m.GetWindowsInformationProtectionPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetWindowsInformationProtectionPolicies()))
        for i, v := range m.GetWindowsInformationProtectionPolicies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("windowsInformationProtectionPolicies", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAndroidManagedAppProtections sets the androidManagedAppProtections property value. Android managed app policies.
func (m *DeviceAppManagement) SetAndroidManagedAppProtections(value []AndroidManagedAppProtectionable)() {
    err := m.GetBackingStore().Set("androidManagedAppProtections", value)
    if err != nil {
        panic(err)
    }
}
// SetDefaultManagedAppProtections sets the defaultManagedAppProtections property value. Default managed app policies.
func (m *DeviceAppManagement) SetDefaultManagedAppProtections(value []DefaultManagedAppProtectionable)() {
    err := m.GetBackingStore().Set("defaultManagedAppProtections", value)
    if err != nil {
        panic(err)
    }
}
// SetIosManagedAppProtections sets the iosManagedAppProtections property value. iOS managed app policies.
func (m *DeviceAppManagement) SetIosManagedAppProtections(value []IosManagedAppProtectionable)() {
    err := m.GetBackingStore().Set("iosManagedAppProtections", value)
    if err != nil {
        panic(err)
    }
}
// SetIsEnabledForMicrosoftStoreForBusiness sets the isEnabledForMicrosoftStoreForBusiness property value. Whether the account is enabled for syncing applications from the Microsoft Store for Business.
func (m *DeviceAppManagement) SetIsEnabledForMicrosoftStoreForBusiness(value *bool)() {
    err := m.GetBackingStore().Set("isEnabledForMicrosoftStoreForBusiness", value)
    if err != nil {
        panic(err)
    }
}
// SetManagedAppPolicies sets the managedAppPolicies property value. Managed app policies.
func (m *DeviceAppManagement) SetManagedAppPolicies(value []ManagedAppPolicyable)() {
    err := m.GetBackingStore().Set("managedAppPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetManagedAppRegistrations sets the managedAppRegistrations property value. The managed app registrations.
func (m *DeviceAppManagement) SetManagedAppRegistrations(value []ManagedAppRegistrationable)() {
    err := m.GetBackingStore().Set("managedAppRegistrations", value)
    if err != nil {
        panic(err)
    }
}
// SetManagedAppStatuses sets the managedAppStatuses property value. The managed app statuses.
func (m *DeviceAppManagement) SetManagedAppStatuses(value []ManagedAppStatusable)() {
    err := m.GetBackingStore().Set("managedAppStatuses", value)
    if err != nil {
        panic(err)
    }
}
// SetManagedEBooks sets the managedEBooks property value. The Managed eBook.
func (m *DeviceAppManagement) SetManagedEBooks(value []ManagedEBookable)() {
    err := m.GetBackingStore().Set("managedEBooks", value)
    if err != nil {
        panic(err)
    }
}
// SetMdmWindowsInformationProtectionPolicies sets the mdmWindowsInformationProtectionPolicies property value. Windows information protection for apps running on devices which are MDM enrolled.
func (m *DeviceAppManagement) SetMdmWindowsInformationProtectionPolicies(value []MdmWindowsInformationProtectionPolicyable)() {
    err := m.GetBackingStore().Set("mdmWindowsInformationProtectionPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetMicrosoftStoreForBusinessLanguage sets the microsoftStoreForBusinessLanguage property value. The locale information used to sync applications from the Microsoft Store for Business. Cultures that are specific to a country/region. The names of these cultures follow RFC 4646 (Windows Vista and later). The format is -<country/regioncode2>, where  is a lowercase two-letter code derived from ISO 639-1 and <country/regioncode2> is an uppercase two-letter code derived from ISO 3166. For example, en-US for English (United States) is a specific culture.
func (m *DeviceAppManagement) SetMicrosoftStoreForBusinessLanguage(value *string)() {
    err := m.GetBackingStore().Set("microsoftStoreForBusinessLanguage", value)
    if err != nil {
        panic(err)
    }
}
// SetMicrosoftStoreForBusinessLastCompletedApplicationSyncTime sets the microsoftStoreForBusinessLastCompletedApplicationSyncTime property value. The last time an application sync from the Microsoft Store for Business was completed.
func (m *DeviceAppManagement) SetMicrosoftStoreForBusinessLastCompletedApplicationSyncTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("microsoftStoreForBusinessLastCompletedApplicationSyncTime", value)
    if err != nil {
        panic(err)
    }
}
// SetMicrosoftStoreForBusinessLastSuccessfulSyncDateTime sets the microsoftStoreForBusinessLastSuccessfulSyncDateTime property value. The last time the apps from the Microsoft Store for Business were synced successfully for the account.
func (m *DeviceAppManagement) SetMicrosoftStoreForBusinessLastSuccessfulSyncDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("microsoftStoreForBusinessLastSuccessfulSyncDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetMobileAppCategories sets the mobileAppCategories property value. The mobile app categories.
func (m *DeviceAppManagement) SetMobileAppCategories(value []MobileAppCategoryable)() {
    err := m.GetBackingStore().Set("mobileAppCategories", value)
    if err != nil {
        panic(err)
    }
}
// SetMobileAppConfigurations sets the mobileAppConfigurations property value. The Managed Device Mobile Application Configurations.
func (m *DeviceAppManagement) SetMobileAppConfigurations(value []ManagedDeviceMobileAppConfigurationable)() {
    err := m.GetBackingStore().Set("mobileAppConfigurations", value)
    if err != nil {
        panic(err)
    }
}
// SetMobileApps sets the mobileApps property value. The mobile apps.
func (m *DeviceAppManagement) SetMobileApps(value []MobileAppable)() {
    err := m.GetBackingStore().Set("mobileApps", value)
    if err != nil {
        panic(err)
    }
}
// SetTargetedManagedAppConfigurations sets the targetedManagedAppConfigurations property value. Targeted managed app configurations.
func (m *DeviceAppManagement) SetTargetedManagedAppConfigurations(value []TargetedManagedAppConfigurationable)() {
    err := m.GetBackingStore().Set("targetedManagedAppConfigurations", value)
    if err != nil {
        panic(err)
    }
}
// SetVppTokens sets the vppTokens property value. List of Vpp tokens for this organization.
func (m *DeviceAppManagement) SetVppTokens(value []VppTokenable)() {
    err := m.GetBackingStore().Set("vppTokens", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsInformationProtectionPolicies sets the windowsInformationProtectionPolicies property value. Windows information protection for apps running on devices which are not MDM enrolled.
func (m *DeviceAppManagement) SetWindowsInformationProtectionPolicies(value []WindowsInformationProtectionPolicyable)() {
    err := m.GetBackingStore().Set("windowsInformationProtectionPolicies", value)
    if err != nil {
        panic(err)
    }
}
type DeviceAppManagementable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAndroidManagedAppProtections()([]AndroidManagedAppProtectionable)
    GetDefaultManagedAppProtections()([]DefaultManagedAppProtectionable)
    GetIosManagedAppProtections()([]IosManagedAppProtectionable)
    GetIsEnabledForMicrosoftStoreForBusiness()(*bool)
    GetManagedAppPolicies()([]ManagedAppPolicyable)
    GetManagedAppRegistrations()([]ManagedAppRegistrationable)
    GetManagedAppStatuses()([]ManagedAppStatusable)
    GetManagedEBooks()([]ManagedEBookable)
    GetMdmWindowsInformationProtectionPolicies()([]MdmWindowsInformationProtectionPolicyable)
    GetMicrosoftStoreForBusinessLanguage()(*string)
    GetMicrosoftStoreForBusinessLastCompletedApplicationSyncTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetMicrosoftStoreForBusinessLastSuccessfulSyncDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetMobileAppCategories()([]MobileAppCategoryable)
    GetMobileAppConfigurations()([]ManagedDeviceMobileAppConfigurationable)
    GetMobileApps()([]MobileAppable)
    GetTargetedManagedAppConfigurations()([]TargetedManagedAppConfigurationable)
    GetVppTokens()([]VppTokenable)
    GetWindowsInformationProtectionPolicies()([]WindowsInformationProtectionPolicyable)
    SetAndroidManagedAppProtections(value []AndroidManagedAppProtectionable)()
    SetDefaultManagedAppProtections(value []DefaultManagedAppProtectionable)()
    SetIosManagedAppProtections(value []IosManagedAppProtectionable)()
    SetIsEnabledForMicrosoftStoreForBusiness(value *bool)()
    SetManagedAppPolicies(value []ManagedAppPolicyable)()
    SetManagedAppRegistrations(value []ManagedAppRegistrationable)()
    SetManagedAppStatuses(value []ManagedAppStatusable)()
    SetManagedEBooks(value []ManagedEBookable)()
    SetMdmWindowsInformationProtectionPolicies(value []MdmWindowsInformationProtectionPolicyable)()
    SetMicrosoftStoreForBusinessLanguage(value *string)()
    SetMicrosoftStoreForBusinessLastCompletedApplicationSyncTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetMicrosoftStoreForBusinessLastSuccessfulSyncDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetMobileAppCategories(value []MobileAppCategoryable)()
    SetMobileAppConfigurations(value []ManagedDeviceMobileAppConfigurationable)()
    SetMobileApps(value []MobileAppable)()
    SetTargetedManagedAppConfigurations(value []TargetedManagedAppConfigurationable)()
    SetVppTokens(value []VppTokenable)()
    SetWindowsInformationProtectionPolicies(value []WindowsInformationProtectionPolicyable)()
}
