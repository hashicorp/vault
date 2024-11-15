package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type ConditionalAccessConditionSet struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewConditionalAccessConditionSet instantiates a new ConditionalAccessConditionSet and sets the default values.
func NewConditionalAccessConditionSet()(*ConditionalAccessConditionSet) {
    m := &ConditionalAccessConditionSet{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateConditionalAccessConditionSetFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateConditionalAccessConditionSetFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConditionalAccessConditionSet(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *ConditionalAccessConditionSet) GetAdditionalData()(map[string]any) {
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
// GetApplications gets the applications property value. Applications and user actions included in and excluded from the policy. Required.
// returns a ConditionalAccessApplicationsable when successful
func (m *ConditionalAccessConditionSet) GetApplications()(ConditionalAccessApplicationsable) {
    val, err := m.GetBackingStore().Get("applications")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConditionalAccessApplicationsable)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *ConditionalAccessConditionSet) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetClientApplications gets the clientApplications property value. Client applications (service principals and workload identities) included in and excluded from the policy. Either users or clientApplications is required.
// returns a ConditionalAccessClientApplicationsable when successful
func (m *ConditionalAccessConditionSet) GetClientApplications()(ConditionalAccessClientApplicationsable) {
    val, err := m.GetBackingStore().Get("clientApplications")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConditionalAccessClientApplicationsable)
    }
    return nil
}
// GetClientAppTypes gets the clientAppTypes property value. Client application types included in the policy. Possible values are: all, browser, mobileAppsAndDesktopClients, exchangeActiveSync, easSupported, other. Required.  The easUnsupported enumeration member will be deprecated in favor of exchangeActiveSync which includes EAS supported and unsupported platforms.
// returns a []ConditionalAccessClientApp when successful
func (m *ConditionalAccessConditionSet) GetClientAppTypes()([]ConditionalAccessClientApp) {
    val, err := m.GetBackingStore().Get("clientAppTypes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ConditionalAccessClientApp)
    }
    return nil
}
// GetDevices gets the devices property value. Devices in the policy.
// returns a ConditionalAccessDevicesable when successful
func (m *ConditionalAccessConditionSet) GetDevices()(ConditionalAccessDevicesable) {
    val, err := m.GetBackingStore().Get("devices")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConditionalAccessDevicesable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ConditionalAccessConditionSet) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["applications"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConditionalAccessApplicationsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplications(val.(ConditionalAccessApplicationsable))
        }
        return nil
    }
    res["clientApplications"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConditionalAccessClientApplicationsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClientApplications(val.(ConditionalAccessClientApplicationsable))
        }
        return nil
    }
    res["clientAppTypes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseConditionalAccessClientApp)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ConditionalAccessClientApp, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*ConditionalAccessClientApp))
                }
            }
            m.SetClientAppTypes(res)
        }
        return nil
    }
    res["devices"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConditionalAccessDevicesFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDevices(val.(ConditionalAccessDevicesable))
        }
        return nil
    }
    res["insiderRiskLevels"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseConditionalAccessInsiderRiskLevels)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInsiderRiskLevels(val.(*ConditionalAccessInsiderRiskLevels))
        }
        return nil
    }
    res["locations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConditionalAccessLocationsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLocations(val.(ConditionalAccessLocationsable))
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
    res["platforms"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConditionalAccessPlatformsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPlatforms(val.(ConditionalAccessPlatformsable))
        }
        return nil
    }
    res["servicePrincipalRiskLevels"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseRiskLevel)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RiskLevel, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*RiskLevel))
                }
            }
            m.SetServicePrincipalRiskLevels(res)
        }
        return nil
    }
    res["signInRiskLevels"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseRiskLevel)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RiskLevel, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*RiskLevel))
                }
            }
            m.SetSignInRiskLevels(res)
        }
        return nil
    }
    res["userRiskLevels"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseRiskLevel)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RiskLevel, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*RiskLevel))
                }
            }
            m.SetUserRiskLevels(res)
        }
        return nil
    }
    res["users"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConditionalAccessUsersFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUsers(val.(ConditionalAccessUsersable))
        }
        return nil
    }
    return res
}
// GetInsiderRiskLevels gets the insiderRiskLevels property value. Insider risk levels included in the policy. The possible values are: minor, moderate, elevated, unknownFutureValue.
// returns a *ConditionalAccessInsiderRiskLevels when successful
func (m *ConditionalAccessConditionSet) GetInsiderRiskLevels()(*ConditionalAccessInsiderRiskLevels) {
    val, err := m.GetBackingStore().Get("insiderRiskLevels")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ConditionalAccessInsiderRiskLevels)
    }
    return nil
}
// GetLocations gets the locations property value. Locations included in and excluded from the policy.
// returns a ConditionalAccessLocationsable when successful
func (m *ConditionalAccessConditionSet) GetLocations()(ConditionalAccessLocationsable) {
    val, err := m.GetBackingStore().Get("locations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConditionalAccessLocationsable)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *ConditionalAccessConditionSet) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPlatforms gets the platforms property value. Platforms included in and excluded from the policy.
// returns a ConditionalAccessPlatformsable when successful
func (m *ConditionalAccessConditionSet) GetPlatforms()(ConditionalAccessPlatformsable) {
    val, err := m.GetBackingStore().Get("platforms")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConditionalAccessPlatformsable)
    }
    return nil
}
// GetServicePrincipalRiskLevels gets the servicePrincipalRiskLevels property value. Service principal risk levels included in the policy. Possible values are: low, medium, high, none, unknownFutureValue.
// returns a []RiskLevel when successful
func (m *ConditionalAccessConditionSet) GetServicePrincipalRiskLevels()([]RiskLevel) {
    val, err := m.GetBackingStore().Get("servicePrincipalRiskLevels")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RiskLevel)
    }
    return nil
}
// GetSignInRiskLevels gets the signInRiskLevels property value. Sign-in risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue. Required.
// returns a []RiskLevel when successful
func (m *ConditionalAccessConditionSet) GetSignInRiskLevels()([]RiskLevel) {
    val, err := m.GetBackingStore().Get("signInRiskLevels")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RiskLevel)
    }
    return nil
}
// GetUserRiskLevels gets the userRiskLevels property value. User risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue. Required.
// returns a []RiskLevel when successful
func (m *ConditionalAccessConditionSet) GetUserRiskLevels()([]RiskLevel) {
    val, err := m.GetBackingStore().Get("userRiskLevels")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RiskLevel)
    }
    return nil
}
// GetUsers gets the users property value. Users, groups, and roles included in and excluded from the policy. Either users or clientApplications is required.
// returns a ConditionalAccessUsersable when successful
func (m *ConditionalAccessConditionSet) GetUsers()(ConditionalAccessUsersable) {
    val, err := m.GetBackingStore().Get("users")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConditionalAccessUsersable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ConditionalAccessConditionSet) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("applications", m.GetApplications())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("clientApplications", m.GetClientApplications())
        if err != nil {
            return err
        }
    }
    if m.GetClientAppTypes() != nil {
        err := writer.WriteCollectionOfStringValues("clientAppTypes", SerializeConditionalAccessClientApp(m.GetClientAppTypes()))
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("devices", m.GetDevices())
        if err != nil {
            return err
        }
    }
    if m.GetInsiderRiskLevels() != nil {
        cast := (*m.GetInsiderRiskLevels()).String()
        err := writer.WriteStringValue("insiderRiskLevels", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("locations", m.GetLocations())
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
        err := writer.WriteObjectValue("platforms", m.GetPlatforms())
        if err != nil {
            return err
        }
    }
    if m.GetServicePrincipalRiskLevels() != nil {
        err := writer.WriteCollectionOfStringValues("servicePrincipalRiskLevels", SerializeRiskLevel(m.GetServicePrincipalRiskLevels()))
        if err != nil {
            return err
        }
    }
    if m.GetSignInRiskLevels() != nil {
        err := writer.WriteCollectionOfStringValues("signInRiskLevels", SerializeRiskLevel(m.GetSignInRiskLevels()))
        if err != nil {
            return err
        }
    }
    if m.GetUserRiskLevels() != nil {
        err := writer.WriteCollectionOfStringValues("userRiskLevels", SerializeRiskLevel(m.GetUserRiskLevels()))
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("users", m.GetUsers())
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
func (m *ConditionalAccessConditionSet) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetApplications sets the applications property value. Applications and user actions included in and excluded from the policy. Required.
func (m *ConditionalAccessConditionSet) SetApplications(value ConditionalAccessApplicationsable)() {
    err := m.GetBackingStore().Set("applications", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *ConditionalAccessConditionSet) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetClientApplications sets the clientApplications property value. Client applications (service principals and workload identities) included in and excluded from the policy. Either users or clientApplications is required.
func (m *ConditionalAccessConditionSet) SetClientApplications(value ConditionalAccessClientApplicationsable)() {
    err := m.GetBackingStore().Set("clientApplications", value)
    if err != nil {
        panic(err)
    }
}
// SetClientAppTypes sets the clientAppTypes property value. Client application types included in the policy. Possible values are: all, browser, mobileAppsAndDesktopClients, exchangeActiveSync, easSupported, other. Required.  The easUnsupported enumeration member will be deprecated in favor of exchangeActiveSync which includes EAS supported and unsupported platforms.
func (m *ConditionalAccessConditionSet) SetClientAppTypes(value []ConditionalAccessClientApp)() {
    err := m.GetBackingStore().Set("clientAppTypes", value)
    if err != nil {
        panic(err)
    }
}
// SetDevices sets the devices property value. Devices in the policy.
func (m *ConditionalAccessConditionSet) SetDevices(value ConditionalAccessDevicesable)() {
    err := m.GetBackingStore().Set("devices", value)
    if err != nil {
        panic(err)
    }
}
// SetInsiderRiskLevels sets the insiderRiskLevels property value. Insider risk levels included in the policy. The possible values are: minor, moderate, elevated, unknownFutureValue.
func (m *ConditionalAccessConditionSet) SetInsiderRiskLevels(value *ConditionalAccessInsiderRiskLevels)() {
    err := m.GetBackingStore().Set("insiderRiskLevels", value)
    if err != nil {
        panic(err)
    }
}
// SetLocations sets the locations property value. Locations included in and excluded from the policy.
func (m *ConditionalAccessConditionSet) SetLocations(value ConditionalAccessLocationsable)() {
    err := m.GetBackingStore().Set("locations", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ConditionalAccessConditionSet) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPlatforms sets the platforms property value. Platforms included in and excluded from the policy.
func (m *ConditionalAccessConditionSet) SetPlatforms(value ConditionalAccessPlatformsable)() {
    err := m.GetBackingStore().Set("platforms", value)
    if err != nil {
        panic(err)
    }
}
// SetServicePrincipalRiskLevels sets the servicePrincipalRiskLevels property value. Service principal risk levels included in the policy. Possible values are: low, medium, high, none, unknownFutureValue.
func (m *ConditionalAccessConditionSet) SetServicePrincipalRiskLevels(value []RiskLevel)() {
    err := m.GetBackingStore().Set("servicePrincipalRiskLevels", value)
    if err != nil {
        panic(err)
    }
}
// SetSignInRiskLevels sets the signInRiskLevels property value. Sign-in risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue. Required.
func (m *ConditionalAccessConditionSet) SetSignInRiskLevels(value []RiskLevel)() {
    err := m.GetBackingStore().Set("signInRiskLevels", value)
    if err != nil {
        panic(err)
    }
}
// SetUserRiskLevels sets the userRiskLevels property value. User risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue. Required.
func (m *ConditionalAccessConditionSet) SetUserRiskLevels(value []RiskLevel)() {
    err := m.GetBackingStore().Set("userRiskLevels", value)
    if err != nil {
        panic(err)
    }
}
// SetUsers sets the users property value. Users, groups, and roles included in and excluded from the policy. Either users or clientApplications is required.
func (m *ConditionalAccessConditionSet) SetUsers(value ConditionalAccessUsersable)() {
    err := m.GetBackingStore().Set("users", value)
    if err != nil {
        panic(err)
    }
}
type ConditionalAccessConditionSetable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApplications()(ConditionalAccessApplicationsable)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetClientApplications()(ConditionalAccessClientApplicationsable)
    GetClientAppTypes()([]ConditionalAccessClientApp)
    GetDevices()(ConditionalAccessDevicesable)
    GetInsiderRiskLevels()(*ConditionalAccessInsiderRiskLevels)
    GetLocations()(ConditionalAccessLocationsable)
    GetOdataType()(*string)
    GetPlatforms()(ConditionalAccessPlatformsable)
    GetServicePrincipalRiskLevels()([]RiskLevel)
    GetSignInRiskLevels()([]RiskLevel)
    GetUserRiskLevels()([]RiskLevel)
    GetUsers()(ConditionalAccessUsersable)
    SetApplications(value ConditionalAccessApplicationsable)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetClientApplications(value ConditionalAccessClientApplicationsable)()
    SetClientAppTypes(value []ConditionalAccessClientApp)()
    SetDevices(value ConditionalAccessDevicesable)()
    SetInsiderRiskLevels(value *ConditionalAccessInsiderRiskLevels)()
    SetLocations(value ConditionalAccessLocationsable)()
    SetOdataType(value *string)()
    SetPlatforms(value ConditionalAccessPlatformsable)()
    SetServicePrincipalRiskLevels(value []RiskLevel)()
    SetSignInRiskLevels(value []RiskLevel)()
    SetUserRiskLevels(value []RiskLevel)()
    SetUsers(value ConditionalAccessUsersable)()
}
