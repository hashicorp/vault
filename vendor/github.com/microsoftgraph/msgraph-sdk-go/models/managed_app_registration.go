package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ManagedAppRegistration the ManagedAppEntity is the base entity type for all other entity types under app management workflow.
type ManagedAppRegistration struct {
    Entity
}
// NewManagedAppRegistration instantiates a new ManagedAppRegistration and sets the default values.
func NewManagedAppRegistration()(*ManagedAppRegistration) {
    m := &ManagedAppRegistration{
        Entity: *NewEntity(),
    }
    return m
}
// CreateManagedAppRegistrationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateManagedAppRegistrationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.androidManagedAppRegistration":
                        return NewAndroidManagedAppRegistration(), nil
                    case "#microsoft.graph.iosManagedAppRegistration":
                        return NewIosManagedAppRegistration(), nil
                }
            }
        }
    }
    return NewManagedAppRegistration(), nil
}
// GetAppIdentifier gets the appIdentifier property value. The app package Identifier
// returns a MobileAppIdentifierable when successful
func (m *ManagedAppRegistration) GetAppIdentifier()(MobileAppIdentifierable) {
    val, err := m.GetBackingStore().Get("appIdentifier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MobileAppIdentifierable)
    }
    return nil
}
// GetApplicationVersion gets the applicationVersion property value. App version
// returns a *string when successful
func (m *ManagedAppRegistration) GetApplicationVersion()(*string) {
    val, err := m.GetBackingStore().Get("applicationVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAppliedPolicies gets the appliedPolicies property value. Zero or more policys already applied on the registered app when it last synchronized with managment service.
// returns a []ManagedAppPolicyable when successful
func (m *ManagedAppRegistration) GetAppliedPolicies()([]ManagedAppPolicyable) {
    val, err := m.GetBackingStore().Get("appliedPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ManagedAppPolicyable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Date and time of creation
// returns a *Time when successful
func (m *ManagedAppRegistration) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDeviceName gets the deviceName property value. Host device name
// returns a *string when successful
func (m *ManagedAppRegistration) GetDeviceName()(*string) {
    val, err := m.GetBackingStore().Get("deviceName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceTag gets the deviceTag property value. App management SDK generated tag, which helps relate apps hosted on the same device. Not guaranteed to relate apps in all conditions.
// returns a *string when successful
func (m *ManagedAppRegistration) GetDeviceTag()(*string) {
    val, err := m.GetBackingStore().Get("deviceTag")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceType gets the deviceType property value. Host device type
// returns a *string when successful
func (m *ManagedAppRegistration) GetDeviceType()(*string) {
    val, err := m.GetBackingStore().Get("deviceType")
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
func (m *ManagedAppRegistration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["appIdentifier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMobileAppIdentifierFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppIdentifier(val.(MobileAppIdentifierable))
        }
        return nil
    }
    res["applicationVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationVersion(val)
        }
        return nil
    }
    res["appliedPolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetAppliedPolicies(res)
        }
        return nil
    }
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
        }
        return nil
    }
    res["deviceName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceName(val)
        }
        return nil
    }
    res["deviceTag"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceTag(val)
        }
        return nil
    }
    res["deviceType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceType(val)
        }
        return nil
    }
    res["flaggedReasons"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseManagedAppFlaggedReason)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ManagedAppFlaggedReason, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*ManagedAppFlaggedReason))
                }
            }
            m.SetFlaggedReasons(res)
        }
        return nil
    }
    res["intendedPolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetIntendedPolicies(res)
        }
        return nil
    }
    res["lastSyncDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastSyncDateTime(val)
        }
        return nil
    }
    res["managementSdkVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManagementSdkVersion(val)
        }
        return nil
    }
    res["operations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateManagedAppOperationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ManagedAppOperationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ManagedAppOperationable)
                }
            }
            m.SetOperations(res)
        }
        return nil
    }
    res["platformVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPlatformVersion(val)
        }
        return nil
    }
    res["userId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserId(val)
        }
        return nil
    }
    res["version"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVersion(val)
        }
        return nil
    }
    return res
}
// GetFlaggedReasons gets the flaggedReasons property value. Zero or more reasons an app registration is flagged. E.g. app running on rooted device
// returns a []ManagedAppFlaggedReason when successful
func (m *ManagedAppRegistration) GetFlaggedReasons()([]ManagedAppFlaggedReason) {
    val, err := m.GetBackingStore().Get("flaggedReasons")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ManagedAppFlaggedReason)
    }
    return nil
}
// GetIntendedPolicies gets the intendedPolicies property value. Zero or more policies admin intended for the app as of now.
// returns a []ManagedAppPolicyable when successful
func (m *ManagedAppRegistration) GetIntendedPolicies()([]ManagedAppPolicyable) {
    val, err := m.GetBackingStore().Get("intendedPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ManagedAppPolicyable)
    }
    return nil
}
// GetLastSyncDateTime gets the lastSyncDateTime property value. Date and time of last the app synced with management service.
// returns a *Time when successful
func (m *ManagedAppRegistration) GetLastSyncDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastSyncDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetManagementSdkVersion gets the managementSdkVersion property value. App management SDK version
// returns a *string when successful
func (m *ManagedAppRegistration) GetManagementSdkVersion()(*string) {
    val, err := m.GetBackingStore().Get("managementSdkVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOperations gets the operations property value. Zero or more long running operations triggered on the app registration.
// returns a []ManagedAppOperationable when successful
func (m *ManagedAppRegistration) GetOperations()([]ManagedAppOperationable) {
    val, err := m.GetBackingStore().Get("operations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ManagedAppOperationable)
    }
    return nil
}
// GetPlatformVersion gets the platformVersion property value. Operating System version
// returns a *string when successful
func (m *ManagedAppRegistration) GetPlatformVersion()(*string) {
    val, err := m.GetBackingStore().Get("platformVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserId gets the userId property value. The user Id to who this app registration belongs.
// returns a *string when successful
func (m *ManagedAppRegistration) GetUserId()(*string) {
    val, err := m.GetBackingStore().Get("userId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetVersion gets the version property value. Version of the entity.
// returns a *string when successful
func (m *ManagedAppRegistration) GetVersion()(*string) {
    val, err := m.GetBackingStore().Get("version")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ManagedAppRegistration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("appIdentifier", m.GetAppIdentifier())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("applicationVersion", m.GetApplicationVersion())
        if err != nil {
            return err
        }
    }
    if m.GetAppliedPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppliedPolicies()))
        for i, v := range m.GetAppliedPolicies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("appliedPolicies", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceName", m.GetDeviceName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceTag", m.GetDeviceTag())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("deviceType", m.GetDeviceType())
        if err != nil {
            return err
        }
    }
    if m.GetFlaggedReasons() != nil {
        err = writer.WriteCollectionOfStringValues("flaggedReasons", SerializeManagedAppFlaggedReason(m.GetFlaggedReasons()))
        if err != nil {
            return err
        }
    }
    if m.GetIntendedPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetIntendedPolicies()))
        for i, v := range m.GetIntendedPolicies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("intendedPolicies", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastSyncDateTime", m.GetLastSyncDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("managementSdkVersion", m.GetManagementSdkVersion())
        if err != nil {
            return err
        }
    }
    if m.GetOperations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOperations()))
        for i, v := range m.GetOperations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("operations", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("platformVersion", m.GetPlatformVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userId", m.GetUserId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("version", m.GetVersion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppIdentifier sets the appIdentifier property value. The app package Identifier
func (m *ManagedAppRegistration) SetAppIdentifier(value MobileAppIdentifierable)() {
    err := m.GetBackingStore().Set("appIdentifier", value)
    if err != nil {
        panic(err)
    }
}
// SetApplicationVersion sets the applicationVersion property value. App version
func (m *ManagedAppRegistration) SetApplicationVersion(value *string)() {
    err := m.GetBackingStore().Set("applicationVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetAppliedPolicies sets the appliedPolicies property value. Zero or more policys already applied on the registered app when it last synchronized with managment service.
func (m *ManagedAppRegistration) SetAppliedPolicies(value []ManagedAppPolicyable)() {
    err := m.GetBackingStore().Set("appliedPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Date and time of creation
func (m *ManagedAppRegistration) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceName sets the deviceName property value. Host device name
func (m *ManagedAppRegistration) SetDeviceName(value *string)() {
    err := m.GetBackingStore().Set("deviceName", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceTag sets the deviceTag property value. App management SDK generated tag, which helps relate apps hosted on the same device. Not guaranteed to relate apps in all conditions.
func (m *ManagedAppRegistration) SetDeviceTag(value *string)() {
    err := m.GetBackingStore().Set("deviceTag", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceType sets the deviceType property value. Host device type
func (m *ManagedAppRegistration) SetDeviceType(value *string)() {
    err := m.GetBackingStore().Set("deviceType", value)
    if err != nil {
        panic(err)
    }
}
// SetFlaggedReasons sets the flaggedReasons property value. Zero or more reasons an app registration is flagged. E.g. app running on rooted device
func (m *ManagedAppRegistration) SetFlaggedReasons(value []ManagedAppFlaggedReason)() {
    err := m.GetBackingStore().Set("flaggedReasons", value)
    if err != nil {
        panic(err)
    }
}
// SetIntendedPolicies sets the intendedPolicies property value. Zero or more policies admin intended for the app as of now.
func (m *ManagedAppRegistration) SetIntendedPolicies(value []ManagedAppPolicyable)() {
    err := m.GetBackingStore().Set("intendedPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetLastSyncDateTime sets the lastSyncDateTime property value. Date and time of last the app synced with management service.
func (m *ManagedAppRegistration) SetLastSyncDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastSyncDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetManagementSdkVersion sets the managementSdkVersion property value. App management SDK version
func (m *ManagedAppRegistration) SetManagementSdkVersion(value *string)() {
    err := m.GetBackingStore().Set("managementSdkVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetOperations sets the operations property value. Zero or more long running operations triggered on the app registration.
func (m *ManagedAppRegistration) SetOperations(value []ManagedAppOperationable)() {
    err := m.GetBackingStore().Set("operations", value)
    if err != nil {
        panic(err)
    }
}
// SetPlatformVersion sets the platformVersion property value. Operating System version
func (m *ManagedAppRegistration) SetPlatformVersion(value *string)() {
    err := m.GetBackingStore().Set("platformVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetUserId sets the userId property value. The user Id to who this app registration belongs.
func (m *ManagedAppRegistration) SetUserId(value *string)() {
    err := m.GetBackingStore().Set("userId", value)
    if err != nil {
        panic(err)
    }
}
// SetVersion sets the version property value. Version of the entity.
func (m *ManagedAppRegistration) SetVersion(value *string)() {
    err := m.GetBackingStore().Set("version", value)
    if err != nil {
        panic(err)
    }
}
type ManagedAppRegistrationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppIdentifier()(MobileAppIdentifierable)
    GetApplicationVersion()(*string)
    GetAppliedPolicies()([]ManagedAppPolicyable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDeviceName()(*string)
    GetDeviceTag()(*string)
    GetDeviceType()(*string)
    GetFlaggedReasons()([]ManagedAppFlaggedReason)
    GetIntendedPolicies()([]ManagedAppPolicyable)
    GetLastSyncDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetManagementSdkVersion()(*string)
    GetOperations()([]ManagedAppOperationable)
    GetPlatformVersion()(*string)
    GetUserId()(*string)
    GetVersion()(*string)
    SetAppIdentifier(value MobileAppIdentifierable)()
    SetApplicationVersion(value *string)()
    SetAppliedPolicies(value []ManagedAppPolicyable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDeviceName(value *string)()
    SetDeviceTag(value *string)()
    SetDeviceType(value *string)()
    SetFlaggedReasons(value []ManagedAppFlaggedReason)()
    SetIntendedPolicies(value []ManagedAppPolicyable)()
    SetLastSyncDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetManagementSdkVersion(value *string)()
    SetOperations(value []ManagedAppOperationable)()
    SetPlatformVersion(value *string)()
    SetUserId(value *string)()
    SetVersion(value *string)()
}
