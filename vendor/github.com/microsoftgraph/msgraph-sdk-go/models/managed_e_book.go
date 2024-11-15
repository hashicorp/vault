package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ManagedEBook an abstract class containing the base properties for Managed eBook.
type ManagedEBook struct {
    Entity
}
// NewManagedEBook instantiates a new ManagedEBook and sets the default values.
func NewManagedEBook()(*ManagedEBook) {
    m := &ManagedEBook{
        Entity: *NewEntity(),
    }
    return m
}
// CreateManagedEBookFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateManagedEBookFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.iosVppEBook":
                        return NewIosVppEBook(), nil
                }
            }
        }
    }
    return NewManagedEBook(), nil
}
// GetAssignments gets the assignments property value. The list of assignments for this eBook.
// returns a []ManagedEBookAssignmentable when successful
func (m *ManagedEBook) GetAssignments()([]ManagedEBookAssignmentable) {
    val, err := m.GetBackingStore().Get("assignments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ManagedEBookAssignmentable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The date and time when the eBook file was created.
// returns a *Time when successful
func (m *ManagedEBook) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDescription gets the description property value. Description.
// returns a *string when successful
func (m *ManagedEBook) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceStates gets the deviceStates property value. The list of installation states for this eBook.
// returns a []DeviceInstallStateable when successful
func (m *ManagedEBook) GetDeviceStates()([]DeviceInstallStateable) {
    val, err := m.GetBackingStore().Get("deviceStates")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceInstallStateable)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Name of the eBook.
// returns a *string when successful
func (m *ManagedEBook) GetDisplayName()(*string) {
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
func (m *ManagedEBook) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["assignments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateManagedEBookAssignmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ManagedEBookAssignmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ManagedEBookAssignmentable)
                }
            }
            m.SetAssignments(res)
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
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
        }
        return nil
    }
    res["deviceStates"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceInstallStateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceInstallStateable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DeviceInstallStateable)
                }
            }
            m.SetDeviceStates(res)
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
    res["informationUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInformationUrl(val)
        }
        return nil
    }
    res["installSummary"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEBookInstallSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInstallSummary(val.(EBookInstallSummaryable))
        }
        return nil
    }
    res["largeCover"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMimeContentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLargeCover(val.(MimeContentable))
        }
        return nil
    }
    res["lastModifiedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedDateTime(val)
        }
        return nil
    }
    res["privacyInformationUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrivacyInformationUrl(val)
        }
        return nil
    }
    res["publishedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPublishedDateTime(val)
        }
        return nil
    }
    res["publisher"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPublisher(val)
        }
        return nil
    }
    res["userStateSummary"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUserInstallStateSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UserInstallStateSummaryable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UserInstallStateSummaryable)
                }
            }
            m.SetUserStateSummary(res)
        }
        return nil
    }
    return res
}
// GetInformationUrl gets the informationUrl property value. The more information Url.
// returns a *string when successful
func (m *ManagedEBook) GetInformationUrl()(*string) {
    val, err := m.GetBackingStore().Get("informationUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetInstallSummary gets the installSummary property value. Mobile App Install Summary.
// returns a EBookInstallSummaryable when successful
func (m *ManagedEBook) GetInstallSummary()(EBookInstallSummaryable) {
    val, err := m.GetBackingStore().Get("installSummary")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EBookInstallSummaryable)
    }
    return nil
}
// GetLargeCover gets the largeCover property value. Book cover.
// returns a MimeContentable when successful
func (m *ManagedEBook) GetLargeCover()(MimeContentable) {
    val, err := m.GetBackingStore().Get("largeCover")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MimeContentable)
    }
    return nil
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. The date and time when the eBook was last modified.
// returns a *Time when successful
func (m *ManagedEBook) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetPrivacyInformationUrl gets the privacyInformationUrl property value. The privacy statement Url.
// returns a *string when successful
func (m *ManagedEBook) GetPrivacyInformationUrl()(*string) {
    val, err := m.GetBackingStore().Get("privacyInformationUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPublishedDateTime gets the publishedDateTime property value. The date and time when the eBook was published.
// returns a *Time when successful
func (m *ManagedEBook) GetPublishedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("publishedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetPublisher gets the publisher property value. Publisher.
// returns a *string when successful
func (m *ManagedEBook) GetPublisher()(*string) {
    val, err := m.GetBackingStore().Get("publisher")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserStateSummary gets the userStateSummary property value. The list of installation states for this eBook.
// returns a []UserInstallStateSummaryable when successful
func (m *ManagedEBook) GetUserStateSummary()([]UserInstallStateSummaryable) {
    val, err := m.GetBackingStore().Get("userStateSummary")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserInstallStateSummaryable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ManagedEBook) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAssignments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAssignments()))
        for i, v := range m.GetAssignments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("assignments", cast)
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
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    if m.GetDeviceStates() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDeviceStates()))
        for i, v := range m.GetDeviceStates() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("deviceStates", cast)
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
        err = writer.WriteStringValue("informationUrl", m.GetInformationUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("installSummary", m.GetInstallSummary())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("largeCover", m.GetLargeCover())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastModifiedDateTime", m.GetLastModifiedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("privacyInformationUrl", m.GetPrivacyInformationUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("publishedDateTime", m.GetPublishedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("publisher", m.GetPublisher())
        if err != nil {
            return err
        }
    }
    if m.GetUserStateSummary() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetUserStateSummary()))
        for i, v := range m.GetUserStateSummary() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("userStateSummary", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAssignments sets the assignments property value. The list of assignments for this eBook.
func (m *ManagedEBook) SetAssignments(value []ManagedEBookAssignmentable)() {
    err := m.GetBackingStore().Set("assignments", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The date and time when the eBook file was created.
func (m *ManagedEBook) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. Description.
func (m *ManagedEBook) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceStates sets the deviceStates property value. The list of installation states for this eBook.
func (m *ManagedEBook) SetDeviceStates(value []DeviceInstallStateable)() {
    err := m.GetBackingStore().Set("deviceStates", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Name of the eBook.
func (m *ManagedEBook) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetInformationUrl sets the informationUrl property value. The more information Url.
func (m *ManagedEBook) SetInformationUrl(value *string)() {
    err := m.GetBackingStore().Set("informationUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetInstallSummary sets the installSummary property value. Mobile App Install Summary.
func (m *ManagedEBook) SetInstallSummary(value EBookInstallSummaryable)() {
    err := m.GetBackingStore().Set("installSummary", value)
    if err != nil {
        panic(err)
    }
}
// SetLargeCover sets the largeCover property value. Book cover.
func (m *ManagedEBook) SetLargeCover(value MimeContentable)() {
    err := m.GetBackingStore().Set("largeCover", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. The date and time when the eBook was last modified.
func (m *ManagedEBook) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetPrivacyInformationUrl sets the privacyInformationUrl property value. The privacy statement Url.
func (m *ManagedEBook) SetPrivacyInformationUrl(value *string)() {
    err := m.GetBackingStore().Set("privacyInformationUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetPublishedDateTime sets the publishedDateTime property value. The date and time when the eBook was published.
func (m *ManagedEBook) SetPublishedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("publishedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetPublisher sets the publisher property value. Publisher.
func (m *ManagedEBook) SetPublisher(value *string)() {
    err := m.GetBackingStore().Set("publisher", value)
    if err != nil {
        panic(err)
    }
}
// SetUserStateSummary sets the userStateSummary property value. The list of installation states for this eBook.
func (m *ManagedEBook) SetUserStateSummary(value []UserInstallStateSummaryable)() {
    err := m.GetBackingStore().Set("userStateSummary", value)
    if err != nil {
        panic(err)
    }
}
type ManagedEBookable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAssignments()([]ManagedEBookAssignmentable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDescription()(*string)
    GetDeviceStates()([]DeviceInstallStateable)
    GetDisplayName()(*string)
    GetInformationUrl()(*string)
    GetInstallSummary()(EBookInstallSummaryable)
    GetLargeCover()(MimeContentable)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetPrivacyInformationUrl()(*string)
    GetPublishedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetPublisher()(*string)
    GetUserStateSummary()([]UserInstallStateSummaryable)
    SetAssignments(value []ManagedEBookAssignmentable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDescription(value *string)()
    SetDeviceStates(value []DeviceInstallStateable)()
    SetDisplayName(value *string)()
    SetInformationUrl(value *string)()
    SetInstallSummary(value EBookInstallSummaryable)()
    SetLargeCover(value MimeContentable)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetPrivacyInformationUrl(value *string)()
    SetPublishedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetPublisher(value *string)()
    SetUserStateSummary(value []UserInstallStateSummaryable)()
}
