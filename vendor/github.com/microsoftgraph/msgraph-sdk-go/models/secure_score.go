package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SecureScore struct {
    Entity
}
// NewSecureScore instantiates a new SecureScore and sets the default values.
func NewSecureScore()(*SecureScore) {
    m := &SecureScore{
        Entity: *NewEntity(),
    }
    return m
}
// CreateSecureScoreFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSecureScoreFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSecureScore(), nil
}
// GetActiveUserCount gets the activeUserCount property value. Active user count of the given tenant.
// returns a *int32 when successful
func (m *SecureScore) GetActiveUserCount()(*int32) {
    val, err := m.GetBackingStore().Get("activeUserCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetAverageComparativeScores gets the averageComparativeScores property value. Average score by different scopes (for example, average by industry, average by seating) and control category (Identity, Data, Device, Apps, Infrastructure) within the scope.
// returns a []AverageComparativeScoreable when successful
func (m *SecureScore) GetAverageComparativeScores()([]AverageComparativeScoreable) {
    val, err := m.GetBackingStore().Get("averageComparativeScores")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AverageComparativeScoreable)
    }
    return nil
}
// GetAzureTenantId gets the azureTenantId property value. GUID string for tenant ID.
// returns a *string when successful
func (m *SecureScore) GetAzureTenantId()(*string) {
    val, err := m.GetBackingStore().Get("azureTenantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetControlScores gets the controlScores property value. Contains tenant scores for a set of controls.
// returns a []ControlScoreable when successful
func (m *SecureScore) GetControlScores()([]ControlScoreable) {
    val, err := m.GetBackingStore().Get("controlScores")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ControlScoreable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. When the report was created.
// returns a *Time when successful
func (m *SecureScore) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetCurrentScore gets the currentScore property value. Tenant current attained score on specified date.
// returns a *float64 when successful
func (m *SecureScore) GetCurrentScore()(*float64) {
    val, err := m.GetBackingStore().Get("currentScore")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetEnabledServices gets the enabledServices property value. Microsoft-provided services for the tenant (for example, Exchange online, Skype, Sharepoint).
// returns a []string when successful
func (m *SecureScore) GetEnabledServices()([]string) {
    val, err := m.GetBackingStore().Get("enabledServices")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SecureScore) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["activeUserCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActiveUserCount(val)
        }
        return nil
    }
    res["averageComparativeScores"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAverageComparativeScoreFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AverageComparativeScoreable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AverageComparativeScoreable)
                }
            }
            m.SetAverageComparativeScores(res)
        }
        return nil
    }
    res["azureTenantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAzureTenantId(val)
        }
        return nil
    }
    res["controlScores"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateControlScoreFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ControlScoreable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ControlScoreable)
                }
            }
            m.SetControlScores(res)
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
    res["currentScore"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCurrentScore(val)
        }
        return nil
    }
    res["enabledServices"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetEnabledServices(res)
        }
        return nil
    }
    res["licensedUserCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLicensedUserCount(val)
        }
        return nil
    }
    res["maxScore"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMaxScore(val)
        }
        return nil
    }
    res["vendorInformation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSecurityVendorInformationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVendorInformation(val.(SecurityVendorInformationable))
        }
        return nil
    }
    return res
}
// GetLicensedUserCount gets the licensedUserCount property value. Licensed user count of the given tenant.
// returns a *int32 when successful
func (m *SecureScore) GetLicensedUserCount()(*int32) {
    val, err := m.GetBackingStore().Get("licensedUserCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetMaxScore gets the maxScore property value. Tenant maximum possible score on specified date.
// returns a *float64 when successful
func (m *SecureScore) GetMaxScore()(*float64) {
    val, err := m.GetBackingStore().Get("maxScore")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetVendorInformation gets the vendorInformation property value. Complex type containing details about the security product/service vendor, provider, and subprovider (for example, vendor=Microsoft; provider=SecureScore). Required.
// returns a SecurityVendorInformationable when successful
func (m *SecureScore) GetVendorInformation()(SecurityVendorInformationable) {
    val, err := m.GetBackingStore().Get("vendorInformation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SecurityVendorInformationable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SecureScore) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("activeUserCount", m.GetActiveUserCount())
        if err != nil {
            return err
        }
    }
    if m.GetAverageComparativeScores() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAverageComparativeScores()))
        for i, v := range m.GetAverageComparativeScores() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("averageComparativeScores", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("azureTenantId", m.GetAzureTenantId())
        if err != nil {
            return err
        }
    }
    if m.GetControlScores() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetControlScores()))
        for i, v := range m.GetControlScores() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("controlScores", cast)
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
        err = writer.WriteFloat64Value("currentScore", m.GetCurrentScore())
        if err != nil {
            return err
        }
    }
    if m.GetEnabledServices() != nil {
        err = writer.WriteCollectionOfStringValues("enabledServices", m.GetEnabledServices())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("licensedUserCount", m.GetLicensedUserCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("maxScore", m.GetMaxScore())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("vendorInformation", m.GetVendorInformation())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActiveUserCount sets the activeUserCount property value. Active user count of the given tenant.
func (m *SecureScore) SetActiveUserCount(value *int32)() {
    err := m.GetBackingStore().Set("activeUserCount", value)
    if err != nil {
        panic(err)
    }
}
// SetAverageComparativeScores sets the averageComparativeScores property value. Average score by different scopes (for example, average by industry, average by seating) and control category (Identity, Data, Device, Apps, Infrastructure) within the scope.
func (m *SecureScore) SetAverageComparativeScores(value []AverageComparativeScoreable)() {
    err := m.GetBackingStore().Set("averageComparativeScores", value)
    if err != nil {
        panic(err)
    }
}
// SetAzureTenantId sets the azureTenantId property value. GUID string for tenant ID.
func (m *SecureScore) SetAzureTenantId(value *string)() {
    err := m.GetBackingStore().Set("azureTenantId", value)
    if err != nil {
        panic(err)
    }
}
// SetControlScores sets the controlScores property value. Contains tenant scores for a set of controls.
func (m *SecureScore) SetControlScores(value []ControlScoreable)() {
    err := m.GetBackingStore().Set("controlScores", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. When the report was created.
func (m *SecureScore) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetCurrentScore sets the currentScore property value. Tenant current attained score on specified date.
func (m *SecureScore) SetCurrentScore(value *float64)() {
    err := m.GetBackingStore().Set("currentScore", value)
    if err != nil {
        panic(err)
    }
}
// SetEnabledServices sets the enabledServices property value. Microsoft-provided services for the tenant (for example, Exchange online, Skype, Sharepoint).
func (m *SecureScore) SetEnabledServices(value []string)() {
    err := m.GetBackingStore().Set("enabledServices", value)
    if err != nil {
        panic(err)
    }
}
// SetLicensedUserCount sets the licensedUserCount property value. Licensed user count of the given tenant.
func (m *SecureScore) SetLicensedUserCount(value *int32)() {
    err := m.GetBackingStore().Set("licensedUserCount", value)
    if err != nil {
        panic(err)
    }
}
// SetMaxScore sets the maxScore property value. Tenant maximum possible score on specified date.
func (m *SecureScore) SetMaxScore(value *float64)() {
    err := m.GetBackingStore().Set("maxScore", value)
    if err != nil {
        panic(err)
    }
}
// SetVendorInformation sets the vendorInformation property value. Complex type containing details about the security product/service vendor, provider, and subprovider (for example, vendor=Microsoft; provider=SecureScore). Required.
func (m *SecureScore) SetVendorInformation(value SecurityVendorInformationable)() {
    err := m.GetBackingStore().Set("vendorInformation", value)
    if err != nil {
        panic(err)
    }
}
type SecureScoreable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetActiveUserCount()(*int32)
    GetAverageComparativeScores()([]AverageComparativeScoreable)
    GetAzureTenantId()(*string)
    GetControlScores()([]ControlScoreable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetCurrentScore()(*float64)
    GetEnabledServices()([]string)
    GetLicensedUserCount()(*int32)
    GetMaxScore()(*float64)
    GetVendorInformation()(SecurityVendorInformationable)
    SetActiveUserCount(value *int32)()
    SetAverageComparativeScores(value []AverageComparativeScoreable)()
    SetAzureTenantId(value *string)()
    SetControlScores(value []ControlScoreable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetCurrentScore(value *float64)()
    SetEnabledServices(value []string)()
    SetLicensedUserCount(value *int32)()
    SetMaxScore(value *float64)()
    SetVendorInformation(value SecurityVendorInformationable)()
}
