package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ManagedAppPolicyDeploymentSummary the ManagedAppEntity is the base entity type for all other entity types under app management workflow.
type ManagedAppPolicyDeploymentSummary struct {
    Entity
}
// NewManagedAppPolicyDeploymentSummary instantiates a new ManagedAppPolicyDeploymentSummary and sets the default values.
func NewManagedAppPolicyDeploymentSummary()(*ManagedAppPolicyDeploymentSummary) {
    m := &ManagedAppPolicyDeploymentSummary{
        Entity: *NewEntity(),
    }
    return m
}
// CreateManagedAppPolicyDeploymentSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateManagedAppPolicyDeploymentSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewManagedAppPolicyDeploymentSummary(), nil
}
// GetConfigurationDeployedUserCount gets the configurationDeployedUserCount property value. Not yet documented
// returns a *int32 when successful
func (m *ManagedAppPolicyDeploymentSummary) GetConfigurationDeployedUserCount()(*int32) {
    val, err := m.GetBackingStore().Get("configurationDeployedUserCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetConfigurationDeploymentSummaryPerApp gets the configurationDeploymentSummaryPerApp property value. Not yet documented
// returns a []ManagedAppPolicyDeploymentSummaryPerAppable when successful
func (m *ManagedAppPolicyDeploymentSummary) GetConfigurationDeploymentSummaryPerApp()([]ManagedAppPolicyDeploymentSummaryPerAppable) {
    val, err := m.GetBackingStore().Get("configurationDeploymentSummaryPerApp")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ManagedAppPolicyDeploymentSummaryPerAppable)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Not yet documented
// returns a *string when successful
func (m *ManagedAppPolicyDeploymentSummary) GetDisplayName()(*string) {
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
func (m *ManagedAppPolicyDeploymentSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["configurationDeployedUserCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConfigurationDeployedUserCount(val)
        }
        return nil
    }
    res["configurationDeploymentSummaryPerApp"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateManagedAppPolicyDeploymentSummaryPerAppFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ManagedAppPolicyDeploymentSummaryPerAppable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ManagedAppPolicyDeploymentSummaryPerAppable)
                }
            }
            m.SetConfigurationDeploymentSummaryPerApp(res)
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
    res["lastRefreshTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastRefreshTime(val)
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
// GetLastRefreshTime gets the lastRefreshTime property value. Not yet documented
// returns a *Time when successful
func (m *ManagedAppPolicyDeploymentSummary) GetLastRefreshTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastRefreshTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetVersion gets the version property value. Version of the entity.
// returns a *string when successful
func (m *ManagedAppPolicyDeploymentSummary) GetVersion()(*string) {
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
func (m *ManagedAppPolicyDeploymentSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("configurationDeployedUserCount", m.GetConfigurationDeployedUserCount())
        if err != nil {
            return err
        }
    }
    if m.GetConfigurationDeploymentSummaryPerApp() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetConfigurationDeploymentSummaryPerApp()))
        for i, v := range m.GetConfigurationDeploymentSummaryPerApp() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("configurationDeploymentSummaryPerApp", cast)
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
    {
        err = writer.WriteTimeValue("lastRefreshTime", m.GetLastRefreshTime())
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
// SetConfigurationDeployedUserCount sets the configurationDeployedUserCount property value. Not yet documented
func (m *ManagedAppPolicyDeploymentSummary) SetConfigurationDeployedUserCount(value *int32)() {
    err := m.GetBackingStore().Set("configurationDeployedUserCount", value)
    if err != nil {
        panic(err)
    }
}
// SetConfigurationDeploymentSummaryPerApp sets the configurationDeploymentSummaryPerApp property value. Not yet documented
func (m *ManagedAppPolicyDeploymentSummary) SetConfigurationDeploymentSummaryPerApp(value []ManagedAppPolicyDeploymentSummaryPerAppable)() {
    err := m.GetBackingStore().Set("configurationDeploymentSummaryPerApp", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Not yet documented
func (m *ManagedAppPolicyDeploymentSummary) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetLastRefreshTime sets the lastRefreshTime property value. Not yet documented
func (m *ManagedAppPolicyDeploymentSummary) SetLastRefreshTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastRefreshTime", value)
    if err != nil {
        panic(err)
    }
}
// SetVersion sets the version property value. Version of the entity.
func (m *ManagedAppPolicyDeploymentSummary) SetVersion(value *string)() {
    err := m.GetBackingStore().Set("version", value)
    if err != nil {
        panic(err)
    }
}
type ManagedAppPolicyDeploymentSummaryable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetConfigurationDeployedUserCount()(*int32)
    GetConfigurationDeploymentSummaryPerApp()([]ManagedAppPolicyDeploymentSummaryPerAppable)
    GetDisplayName()(*string)
    GetLastRefreshTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetVersion()(*string)
    SetConfigurationDeployedUserCount(value *int32)()
    SetConfigurationDeploymentSummaryPerApp(value []ManagedAppPolicyDeploymentSummaryPerAppable)()
    SetDisplayName(value *string)()
    SetLastRefreshTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetVersion(value *string)()
}
