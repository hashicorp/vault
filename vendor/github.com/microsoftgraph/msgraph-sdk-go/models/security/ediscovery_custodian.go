package security

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EdiscoveryCustodian struct {
    DataSourceContainer
}
// NewEdiscoveryCustodian instantiates a new EdiscoveryCustodian and sets the default values.
func NewEdiscoveryCustodian()(*EdiscoveryCustodian) {
    m := &EdiscoveryCustodian{
        DataSourceContainer: *NewDataSourceContainer(),
    }
    odataTypeValue := "#microsoft.graph.security.ediscoveryCustodian"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEdiscoveryCustodianFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEdiscoveryCustodianFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEdiscoveryCustodian(), nil
}
// GetAcknowledgedDateTime gets the acknowledgedDateTime property value. Date and time the custodian acknowledged a hold notification.
// returns a *Time when successful
func (m *EdiscoveryCustodian) GetAcknowledgedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("acknowledgedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetEmail gets the email property value. Email address of the custodian.
// returns a *string when successful
func (m *EdiscoveryCustodian) GetEmail()(*string) {
    val, err := m.GetBackingStore().Get("email")
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
func (m *EdiscoveryCustodian) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DataSourceContainer.GetFieldDeserializers()
    res["acknowledgedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAcknowledgedDateTime(val)
        }
        return nil
    }
    res["email"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmail(val)
        }
        return nil
    }
    res["lastIndexOperation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEdiscoveryIndexOperationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastIndexOperation(val.(EdiscoveryIndexOperationable))
        }
        return nil
    }
    res["siteSources"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSiteSourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SiteSourceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SiteSourceable)
                }
            }
            m.SetSiteSources(res)
        }
        return nil
    }
    res["unifiedGroupSources"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUnifiedGroupSourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UnifiedGroupSourceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UnifiedGroupSourceable)
                }
            }
            m.SetUnifiedGroupSources(res)
        }
        return nil
    }
    res["userSources"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUserSourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UserSourceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UserSourceable)
                }
            }
            m.SetUserSources(res)
        }
        return nil
    }
    return res
}
// GetLastIndexOperation gets the lastIndexOperation property value. Operation entity that represents the latest indexing for the custodian.
// returns a EdiscoveryIndexOperationable when successful
func (m *EdiscoveryCustodian) GetLastIndexOperation()(EdiscoveryIndexOperationable) {
    val, err := m.GetBackingStore().Get("lastIndexOperation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EdiscoveryIndexOperationable)
    }
    return nil
}
// GetSiteSources gets the siteSources property value. Data source entity for SharePoint sites associated with the custodian.
// returns a []SiteSourceable when successful
func (m *EdiscoveryCustodian) GetSiteSources()([]SiteSourceable) {
    val, err := m.GetBackingStore().Get("siteSources")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SiteSourceable)
    }
    return nil
}
// GetUnifiedGroupSources gets the unifiedGroupSources property value. Data source entity for groups associated with the custodian.
// returns a []UnifiedGroupSourceable when successful
func (m *EdiscoveryCustodian) GetUnifiedGroupSources()([]UnifiedGroupSourceable) {
    val, err := m.GetBackingStore().Get("unifiedGroupSources")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UnifiedGroupSourceable)
    }
    return nil
}
// GetUserSources gets the userSources property value. Data source entity for a the custodian. This is the container for a custodian's mailbox and OneDrive for Business site.
// returns a []UserSourceable when successful
func (m *EdiscoveryCustodian) GetUserSources()([]UserSourceable) {
    val, err := m.GetBackingStore().Get("userSources")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserSourceable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EdiscoveryCustodian) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DataSourceContainer.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("acknowledgedDateTime", m.GetAcknowledgedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("email", m.GetEmail())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("lastIndexOperation", m.GetLastIndexOperation())
        if err != nil {
            return err
        }
    }
    if m.GetSiteSources() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSiteSources()))
        for i, v := range m.GetSiteSources() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("siteSources", cast)
        if err != nil {
            return err
        }
    }
    if m.GetUnifiedGroupSources() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetUnifiedGroupSources()))
        for i, v := range m.GetUnifiedGroupSources() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("unifiedGroupSources", cast)
        if err != nil {
            return err
        }
    }
    if m.GetUserSources() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetUserSources()))
        for i, v := range m.GetUserSources() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("userSources", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAcknowledgedDateTime sets the acknowledgedDateTime property value. Date and time the custodian acknowledged a hold notification.
func (m *EdiscoveryCustodian) SetAcknowledgedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("acknowledgedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetEmail sets the email property value. Email address of the custodian.
func (m *EdiscoveryCustodian) SetEmail(value *string)() {
    err := m.GetBackingStore().Set("email", value)
    if err != nil {
        panic(err)
    }
}
// SetLastIndexOperation sets the lastIndexOperation property value. Operation entity that represents the latest indexing for the custodian.
func (m *EdiscoveryCustodian) SetLastIndexOperation(value EdiscoveryIndexOperationable)() {
    err := m.GetBackingStore().Set("lastIndexOperation", value)
    if err != nil {
        panic(err)
    }
}
// SetSiteSources sets the siteSources property value. Data source entity for SharePoint sites associated with the custodian.
func (m *EdiscoveryCustodian) SetSiteSources(value []SiteSourceable)() {
    err := m.GetBackingStore().Set("siteSources", value)
    if err != nil {
        panic(err)
    }
}
// SetUnifiedGroupSources sets the unifiedGroupSources property value. Data source entity for groups associated with the custodian.
func (m *EdiscoveryCustodian) SetUnifiedGroupSources(value []UnifiedGroupSourceable)() {
    err := m.GetBackingStore().Set("unifiedGroupSources", value)
    if err != nil {
        panic(err)
    }
}
// SetUserSources sets the userSources property value. Data source entity for a the custodian. This is the container for a custodian's mailbox and OneDrive for Business site.
func (m *EdiscoveryCustodian) SetUserSources(value []UserSourceable)() {
    err := m.GetBackingStore().Set("userSources", value)
    if err != nil {
        panic(err)
    }
}
type EdiscoveryCustodianable interface {
    DataSourceContainerable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAcknowledgedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetEmail()(*string)
    GetLastIndexOperation()(EdiscoveryIndexOperationable)
    GetSiteSources()([]SiteSourceable)
    GetUnifiedGroupSources()([]UnifiedGroupSourceable)
    GetUserSources()([]UserSourceable)
    SetAcknowledgedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetEmail(value *string)()
    SetLastIndexOperation(value EdiscoveryIndexOperationable)()
    SetSiteSources(value []SiteSourceable)()
    SetUnifiedGroupSources(value []UnifiedGroupSourceable)()
    SetUserSources(value []UserSourceable)()
}
