package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ComplianceManagementPartner compliance management partner for all platforms
type ComplianceManagementPartner struct {
    Entity
}
// NewComplianceManagementPartner instantiates a new ComplianceManagementPartner and sets the default values.
func NewComplianceManagementPartner()(*ComplianceManagementPartner) {
    m := &ComplianceManagementPartner{
        Entity: *NewEntity(),
    }
    return m
}
// CreateComplianceManagementPartnerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateComplianceManagementPartnerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewComplianceManagementPartner(), nil
}
// GetAndroidEnrollmentAssignments gets the androidEnrollmentAssignments property value. User groups which enroll Android devices through partner.
// returns a []ComplianceManagementPartnerAssignmentable when successful
func (m *ComplianceManagementPartner) GetAndroidEnrollmentAssignments()([]ComplianceManagementPartnerAssignmentable) {
    val, err := m.GetBackingStore().Get("androidEnrollmentAssignments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ComplianceManagementPartnerAssignmentable)
    }
    return nil
}
// GetAndroidOnboarded gets the androidOnboarded property value. Partner onboarded for Android devices.
// returns a *bool when successful
func (m *ComplianceManagementPartner) GetAndroidOnboarded()(*bool) {
    val, err := m.GetBackingStore().Get("androidOnboarded")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Partner display name
// returns a *string when successful
func (m *ComplianceManagementPartner) GetDisplayName()(*string) {
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
func (m *ComplianceManagementPartner) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["androidEnrollmentAssignments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateComplianceManagementPartnerAssignmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ComplianceManagementPartnerAssignmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ComplianceManagementPartnerAssignmentable)
                }
            }
            m.SetAndroidEnrollmentAssignments(res)
        }
        return nil
    }
    res["androidOnboarded"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAndroidOnboarded(val)
        }
        return nil
    }
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
    res["iosEnrollmentAssignments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateComplianceManagementPartnerAssignmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ComplianceManagementPartnerAssignmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ComplianceManagementPartnerAssignmentable)
                }
            }
            m.SetIosEnrollmentAssignments(res)
        }
        return nil
    }
    res["iosOnboarded"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIosOnboarded(val)
        }
        return nil
    }
    res["lastHeartbeatDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastHeartbeatDateTime(val)
        }
        return nil
    }
    res["macOsEnrollmentAssignments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateComplianceManagementPartnerAssignmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ComplianceManagementPartnerAssignmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ComplianceManagementPartnerAssignmentable)
                }
            }
            m.SetMacOsEnrollmentAssignments(res)
        }
        return nil
    }
    res["macOsOnboarded"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMacOsOnboarded(val)
        }
        return nil
    }
    res["partnerState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDeviceManagementPartnerTenantState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPartnerState(val.(*DeviceManagementPartnerTenantState))
        }
        return nil
    }
    return res
}
// GetIosEnrollmentAssignments gets the iosEnrollmentAssignments property value. User groups which enroll ios devices through partner.
// returns a []ComplianceManagementPartnerAssignmentable when successful
func (m *ComplianceManagementPartner) GetIosEnrollmentAssignments()([]ComplianceManagementPartnerAssignmentable) {
    val, err := m.GetBackingStore().Get("iosEnrollmentAssignments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ComplianceManagementPartnerAssignmentable)
    }
    return nil
}
// GetIosOnboarded gets the iosOnboarded property value. Partner onboarded for ios devices.
// returns a *bool when successful
func (m *ComplianceManagementPartner) GetIosOnboarded()(*bool) {
    val, err := m.GetBackingStore().Get("iosOnboarded")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLastHeartbeatDateTime gets the lastHeartbeatDateTime property value. Timestamp of last heartbeat after admin onboarded to the compliance management partner
// returns a *Time when successful
func (m *ComplianceManagementPartner) GetLastHeartbeatDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastHeartbeatDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetMacOsEnrollmentAssignments gets the macOsEnrollmentAssignments property value. User groups which enroll Mac devices through partner.
// returns a []ComplianceManagementPartnerAssignmentable when successful
func (m *ComplianceManagementPartner) GetMacOsEnrollmentAssignments()([]ComplianceManagementPartnerAssignmentable) {
    val, err := m.GetBackingStore().Get("macOsEnrollmentAssignments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ComplianceManagementPartnerAssignmentable)
    }
    return nil
}
// GetMacOsOnboarded gets the macOsOnboarded property value. Partner onboarded for Mac devices.
// returns a *bool when successful
func (m *ComplianceManagementPartner) GetMacOsOnboarded()(*bool) {
    val, err := m.GetBackingStore().Get("macOsOnboarded")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPartnerState gets the partnerState property value. Partner state of this tenant.
// returns a *DeviceManagementPartnerTenantState when successful
func (m *ComplianceManagementPartner) GetPartnerState()(*DeviceManagementPartnerTenantState) {
    val, err := m.GetBackingStore().Get("partnerState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DeviceManagementPartnerTenantState)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ComplianceManagementPartner) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAndroidEnrollmentAssignments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAndroidEnrollmentAssignments()))
        for i, v := range m.GetAndroidEnrollmentAssignments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("androidEnrollmentAssignments", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("androidOnboarded", m.GetAndroidOnboarded())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetIosEnrollmentAssignments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetIosEnrollmentAssignments()))
        for i, v := range m.GetIosEnrollmentAssignments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("iosEnrollmentAssignments", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iosOnboarded", m.GetIosOnboarded())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastHeartbeatDateTime", m.GetLastHeartbeatDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetMacOsEnrollmentAssignments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMacOsEnrollmentAssignments()))
        for i, v := range m.GetMacOsEnrollmentAssignments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("macOsEnrollmentAssignments", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("macOsOnboarded", m.GetMacOsOnboarded())
        if err != nil {
            return err
        }
    }
    if m.GetPartnerState() != nil {
        cast := (*m.GetPartnerState()).String()
        err = writer.WriteStringValue("partnerState", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAndroidEnrollmentAssignments sets the androidEnrollmentAssignments property value. User groups which enroll Android devices through partner.
func (m *ComplianceManagementPartner) SetAndroidEnrollmentAssignments(value []ComplianceManagementPartnerAssignmentable)() {
    err := m.GetBackingStore().Set("androidEnrollmentAssignments", value)
    if err != nil {
        panic(err)
    }
}
// SetAndroidOnboarded sets the androidOnboarded property value. Partner onboarded for Android devices.
func (m *ComplianceManagementPartner) SetAndroidOnboarded(value *bool)() {
    err := m.GetBackingStore().Set("androidOnboarded", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Partner display name
func (m *ComplianceManagementPartner) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetIosEnrollmentAssignments sets the iosEnrollmentAssignments property value. User groups which enroll ios devices through partner.
func (m *ComplianceManagementPartner) SetIosEnrollmentAssignments(value []ComplianceManagementPartnerAssignmentable)() {
    err := m.GetBackingStore().Set("iosEnrollmentAssignments", value)
    if err != nil {
        panic(err)
    }
}
// SetIosOnboarded sets the iosOnboarded property value. Partner onboarded for ios devices.
func (m *ComplianceManagementPartner) SetIosOnboarded(value *bool)() {
    err := m.GetBackingStore().Set("iosOnboarded", value)
    if err != nil {
        panic(err)
    }
}
// SetLastHeartbeatDateTime sets the lastHeartbeatDateTime property value. Timestamp of last heartbeat after admin onboarded to the compliance management partner
func (m *ComplianceManagementPartner) SetLastHeartbeatDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastHeartbeatDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetMacOsEnrollmentAssignments sets the macOsEnrollmentAssignments property value. User groups which enroll Mac devices through partner.
func (m *ComplianceManagementPartner) SetMacOsEnrollmentAssignments(value []ComplianceManagementPartnerAssignmentable)() {
    err := m.GetBackingStore().Set("macOsEnrollmentAssignments", value)
    if err != nil {
        panic(err)
    }
}
// SetMacOsOnboarded sets the macOsOnboarded property value. Partner onboarded for Mac devices.
func (m *ComplianceManagementPartner) SetMacOsOnboarded(value *bool)() {
    err := m.GetBackingStore().Set("macOsOnboarded", value)
    if err != nil {
        panic(err)
    }
}
// SetPartnerState sets the partnerState property value. Partner state of this tenant.
func (m *ComplianceManagementPartner) SetPartnerState(value *DeviceManagementPartnerTenantState)() {
    err := m.GetBackingStore().Set("partnerState", value)
    if err != nil {
        panic(err)
    }
}
type ComplianceManagementPartnerable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAndroidEnrollmentAssignments()([]ComplianceManagementPartnerAssignmentable)
    GetAndroidOnboarded()(*bool)
    GetDisplayName()(*string)
    GetIosEnrollmentAssignments()([]ComplianceManagementPartnerAssignmentable)
    GetIosOnboarded()(*bool)
    GetLastHeartbeatDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetMacOsEnrollmentAssignments()([]ComplianceManagementPartnerAssignmentable)
    GetMacOsOnboarded()(*bool)
    GetPartnerState()(*DeviceManagementPartnerTenantState)
    SetAndroidEnrollmentAssignments(value []ComplianceManagementPartnerAssignmentable)()
    SetAndroidOnboarded(value *bool)()
    SetDisplayName(value *string)()
    SetIosEnrollmentAssignments(value []ComplianceManagementPartnerAssignmentable)()
    SetIosOnboarded(value *bool)()
    SetLastHeartbeatDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetMacOsEnrollmentAssignments(value []ComplianceManagementPartnerAssignmentable)()
    SetMacOsOnboarded(value *bool)()
    SetPartnerState(value *DeviceManagementPartnerTenantState)()
}
