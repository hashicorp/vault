package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EducationTeamsAppResource struct {
    EducationResource
}
// NewEducationTeamsAppResource instantiates a new EducationTeamsAppResource and sets the default values.
func NewEducationTeamsAppResource()(*EducationTeamsAppResource) {
    m := &EducationTeamsAppResource{
        EducationResource: *NewEducationResource(),
    }
    odataTypeValue := "#microsoft.graph.educationTeamsAppResource"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEducationTeamsAppResourceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEducationTeamsAppResourceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEducationTeamsAppResource(), nil
}
// GetAppIconWebUrl gets the appIconWebUrl property value. URL that points to the icon of the app.
// returns a *string when successful
func (m *EducationTeamsAppResource) GetAppIconWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("appIconWebUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAppId gets the appId property value. Teams app ID of the application.
// returns a *string when successful
func (m *EducationTeamsAppResource) GetAppId()(*string) {
    val, err := m.GetBackingStore().Get("appId")
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
func (m *EducationTeamsAppResource) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EducationResource.GetFieldDeserializers()
    res["appIconWebUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppIconWebUrl(val)
        }
        return nil
    }
    res["appId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppId(val)
        }
        return nil
    }
    res["teamsEmbeddedContentUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTeamsEmbeddedContentUrl(val)
        }
        return nil
    }
    res["webUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWebUrl(val)
        }
        return nil
    }
    return res
}
// GetTeamsEmbeddedContentUrl gets the teamsEmbeddedContentUrl property value. URL for the app resource that will be opened by Teams.
// returns a *string when successful
func (m *EducationTeamsAppResource) GetTeamsEmbeddedContentUrl()(*string) {
    val, err := m.GetBackingStore().Get("teamsEmbeddedContentUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWebUrl gets the webUrl property value. URL for the app resource that can be opened in the browser.
// returns a *string when successful
func (m *EducationTeamsAppResource) GetWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("webUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EducationTeamsAppResource) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EducationResource.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("appIconWebUrl", m.GetAppIconWebUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("appId", m.GetAppId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("teamsEmbeddedContentUrl", m.GetTeamsEmbeddedContentUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("webUrl", m.GetWebUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppIconWebUrl sets the appIconWebUrl property value. URL that points to the icon of the app.
func (m *EducationTeamsAppResource) SetAppIconWebUrl(value *string)() {
    err := m.GetBackingStore().Set("appIconWebUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetAppId sets the appId property value. Teams app ID of the application.
func (m *EducationTeamsAppResource) SetAppId(value *string)() {
    err := m.GetBackingStore().Set("appId", value)
    if err != nil {
        panic(err)
    }
}
// SetTeamsEmbeddedContentUrl sets the teamsEmbeddedContentUrl property value. URL for the app resource that will be opened by Teams.
func (m *EducationTeamsAppResource) SetTeamsEmbeddedContentUrl(value *string)() {
    err := m.GetBackingStore().Set("teamsEmbeddedContentUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetWebUrl sets the webUrl property value. URL for the app resource that can be opened in the browser.
func (m *EducationTeamsAppResource) SetWebUrl(value *string)() {
    err := m.GetBackingStore().Set("webUrl", value)
    if err != nil {
        panic(err)
    }
}
type EducationTeamsAppResourceable interface {
    EducationResourceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppIconWebUrl()(*string)
    GetAppId()(*string)
    GetTeamsEmbeddedContentUrl()(*string)
    GetWebUrl()(*string)
    SetAppIconWebUrl(value *string)()
    SetAppId(value *string)()
    SetTeamsEmbeddedContentUrl(value *string)()
    SetWebUrl(value *string)()
}
