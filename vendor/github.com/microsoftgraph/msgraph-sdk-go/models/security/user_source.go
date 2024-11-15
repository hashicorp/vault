package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UserSource struct {
    DataSource
}
// NewUserSource instantiates a new UserSource and sets the default values.
func NewUserSource()(*UserSource) {
    m := &UserSource{
        DataSource: *NewDataSource(),
    }
    odataTypeValue := "#microsoft.graph.security.userSource"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateUserSourceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserSourceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserSource(), nil
}
// GetEmail gets the email property value. Email address of the user's mailbox.
// returns a *string when successful
func (m *UserSource) GetEmail()(*string) {
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
func (m *UserSource) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DataSource.GetFieldDeserializers()
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
    res["includedSources"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSourceType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIncludedSources(val.(*SourceType))
        }
        return nil
    }
    res["siteWebUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSiteWebUrl(val)
        }
        return nil
    }
    return res
}
// GetIncludedSources gets the includedSources property value. Specifies which sources are included in this group. Possible values are: mailbox, site.
// returns a *SourceType when successful
func (m *UserSource) GetIncludedSources()(*SourceType) {
    val, err := m.GetBackingStore().Get("includedSources")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SourceType)
    }
    return nil
}
// GetSiteWebUrl gets the siteWebUrl property value. The URL of the user's OneDrive for Business site. Read-only.
// returns a *string when successful
func (m *UserSource) GetSiteWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("siteWebUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserSource) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DataSource.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("email", m.GetEmail())
        if err != nil {
            return err
        }
    }
    if m.GetIncludedSources() != nil {
        cast := (*m.GetIncludedSources()).String()
        err = writer.WriteStringValue("includedSources", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("siteWebUrl", m.GetSiteWebUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetEmail sets the email property value. Email address of the user's mailbox.
func (m *UserSource) SetEmail(value *string)() {
    err := m.GetBackingStore().Set("email", value)
    if err != nil {
        panic(err)
    }
}
// SetIncludedSources sets the includedSources property value. Specifies which sources are included in this group. Possible values are: mailbox, site.
func (m *UserSource) SetIncludedSources(value *SourceType)() {
    err := m.GetBackingStore().Set("includedSources", value)
    if err != nil {
        panic(err)
    }
}
// SetSiteWebUrl sets the siteWebUrl property value. The URL of the user's OneDrive for Business site. Read-only.
func (m *UserSource) SetSiteWebUrl(value *string)() {
    err := m.GetBackingStore().Set("siteWebUrl", value)
    if err != nil {
        panic(err)
    }
}
type UserSourceable interface {
    DataSourceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetEmail()(*string)
    GetIncludedSources()(*SourceType)
    GetSiteWebUrl()(*string)
    SetEmail(value *string)()
    SetIncludedSources(value *SourceType)()
    SetSiteWebUrl(value *string)()
}
