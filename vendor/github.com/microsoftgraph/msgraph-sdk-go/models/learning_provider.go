package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type LearningProvider struct {
    Entity
}
// NewLearningProvider instantiates a new LearningProvider and sets the default values.
func NewLearningProvider()(*LearningProvider) {
    m := &LearningProvider{
        Entity: *NewEntity(),
    }
    return m
}
// CreateLearningProviderFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateLearningProviderFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewLearningProvider(), nil
}
// GetDisplayName gets the displayName property value. The display name that appears in Viva Learning. Required.
// returns a *string when successful
func (m *LearningProvider) GetDisplayName()(*string) {
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
func (m *LearningProvider) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
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
    res["isCourseActivitySyncEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsCourseActivitySyncEnabled(val)
        }
        return nil
    }
    res["learningContents"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateLearningContentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]LearningContentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(LearningContentable)
                }
            }
            m.SetLearningContents(res)
        }
        return nil
    }
    res["learningCourseActivities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateLearningCourseActivityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]LearningCourseActivityable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(LearningCourseActivityable)
                }
            }
            m.SetLearningCourseActivities(res)
        }
        return nil
    }
    res["loginWebUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLoginWebUrl(val)
        }
        return nil
    }
    res["longLogoWebUrlForDarkTheme"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLongLogoWebUrlForDarkTheme(val)
        }
        return nil
    }
    res["longLogoWebUrlForLightTheme"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLongLogoWebUrlForLightTheme(val)
        }
        return nil
    }
    res["squareLogoWebUrlForDarkTheme"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSquareLogoWebUrlForDarkTheme(val)
        }
        return nil
    }
    res["squareLogoWebUrlForLightTheme"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSquareLogoWebUrlForLightTheme(val)
        }
        return nil
    }
    return res
}
// GetIsCourseActivitySyncEnabled gets the isCourseActivitySyncEnabled property value. Indicates whether a provider can ingest learning course activity records. The default value is false. Set to true to make learningCourseActivities available for this provider.
// returns a *bool when successful
func (m *LearningProvider) GetIsCourseActivitySyncEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isCourseActivitySyncEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLearningContents gets the learningContents property value. Learning catalog items for the provider.
// returns a []LearningContentable when successful
func (m *LearningProvider) GetLearningContents()([]LearningContentable) {
    val, err := m.GetBackingStore().Get("learningContents")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]LearningContentable)
    }
    return nil
}
// GetLearningCourseActivities gets the learningCourseActivities property value. The learningCourseActivities property
// returns a []LearningCourseActivityable when successful
func (m *LearningProvider) GetLearningCourseActivities()([]LearningCourseActivityable) {
    val, err := m.GetBackingStore().Get("learningCourseActivities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]LearningCourseActivityable)
    }
    return nil
}
// GetLoginWebUrl gets the loginWebUrl property value. Authentication URL to access the courses for the provider. Optional.
// returns a *string when successful
func (m *LearningProvider) GetLoginWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("loginWebUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLongLogoWebUrlForDarkTheme gets the longLogoWebUrlForDarkTheme property value. The long logo URL for the dark mode that needs to be a publicly accessible image. This image would be saved to the blob storage of Viva Learning for rendering within the Viva Learning app. Required.
// returns a *string when successful
func (m *LearningProvider) GetLongLogoWebUrlForDarkTheme()(*string) {
    val, err := m.GetBackingStore().Get("longLogoWebUrlForDarkTheme")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLongLogoWebUrlForLightTheme gets the longLogoWebUrlForLightTheme property value. The long logo URL for the light mode that needs to be a publicly accessible image. This image would be saved to the blob storage of Viva Learning for rendering within the Viva Learning app. Required.
// returns a *string when successful
func (m *LearningProvider) GetLongLogoWebUrlForLightTheme()(*string) {
    val, err := m.GetBackingStore().Get("longLogoWebUrlForLightTheme")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSquareLogoWebUrlForDarkTheme gets the squareLogoWebUrlForDarkTheme property value. The square logo URL for the dark mode that needs to be a publicly accessible image. This image would be saved to the blob storage of Viva Learning for rendering within the Viva Learning app. Required.
// returns a *string when successful
func (m *LearningProvider) GetSquareLogoWebUrlForDarkTheme()(*string) {
    val, err := m.GetBackingStore().Get("squareLogoWebUrlForDarkTheme")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSquareLogoWebUrlForLightTheme gets the squareLogoWebUrlForLightTheme property value. The square logo URL for the light mode that needs to be a publicly accessible image. This image would be saved to the blob storage of Viva Learning for rendering within the Viva Learning app. Required.
// returns a *string when successful
func (m *LearningProvider) GetSquareLogoWebUrlForLightTheme()(*string) {
    val, err := m.GetBackingStore().Get("squareLogoWebUrlForLightTheme")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *LearningProvider) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isCourseActivitySyncEnabled", m.GetIsCourseActivitySyncEnabled())
        if err != nil {
            return err
        }
    }
    if m.GetLearningContents() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLearningContents()))
        for i, v := range m.GetLearningContents() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("learningContents", cast)
        if err != nil {
            return err
        }
    }
    if m.GetLearningCourseActivities() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLearningCourseActivities()))
        for i, v := range m.GetLearningCourseActivities() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("learningCourseActivities", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("loginWebUrl", m.GetLoginWebUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("longLogoWebUrlForDarkTheme", m.GetLongLogoWebUrlForDarkTheme())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("longLogoWebUrlForLightTheme", m.GetLongLogoWebUrlForLightTheme())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("squareLogoWebUrlForDarkTheme", m.GetSquareLogoWebUrlForDarkTheme())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("squareLogoWebUrlForLightTheme", m.GetSquareLogoWebUrlForLightTheme())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDisplayName sets the displayName property value. The display name that appears in Viva Learning. Required.
func (m *LearningProvider) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetIsCourseActivitySyncEnabled sets the isCourseActivitySyncEnabled property value. Indicates whether a provider can ingest learning course activity records. The default value is false. Set to true to make learningCourseActivities available for this provider.
func (m *LearningProvider) SetIsCourseActivitySyncEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isCourseActivitySyncEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetLearningContents sets the learningContents property value. Learning catalog items for the provider.
func (m *LearningProvider) SetLearningContents(value []LearningContentable)() {
    err := m.GetBackingStore().Set("learningContents", value)
    if err != nil {
        panic(err)
    }
}
// SetLearningCourseActivities sets the learningCourseActivities property value. The learningCourseActivities property
func (m *LearningProvider) SetLearningCourseActivities(value []LearningCourseActivityable)() {
    err := m.GetBackingStore().Set("learningCourseActivities", value)
    if err != nil {
        panic(err)
    }
}
// SetLoginWebUrl sets the loginWebUrl property value. Authentication URL to access the courses for the provider. Optional.
func (m *LearningProvider) SetLoginWebUrl(value *string)() {
    err := m.GetBackingStore().Set("loginWebUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetLongLogoWebUrlForDarkTheme sets the longLogoWebUrlForDarkTheme property value. The long logo URL for the dark mode that needs to be a publicly accessible image. This image would be saved to the blob storage of Viva Learning for rendering within the Viva Learning app. Required.
func (m *LearningProvider) SetLongLogoWebUrlForDarkTheme(value *string)() {
    err := m.GetBackingStore().Set("longLogoWebUrlForDarkTheme", value)
    if err != nil {
        panic(err)
    }
}
// SetLongLogoWebUrlForLightTheme sets the longLogoWebUrlForLightTheme property value. The long logo URL for the light mode that needs to be a publicly accessible image. This image would be saved to the blob storage of Viva Learning for rendering within the Viva Learning app. Required.
func (m *LearningProvider) SetLongLogoWebUrlForLightTheme(value *string)() {
    err := m.GetBackingStore().Set("longLogoWebUrlForLightTheme", value)
    if err != nil {
        panic(err)
    }
}
// SetSquareLogoWebUrlForDarkTheme sets the squareLogoWebUrlForDarkTheme property value. The square logo URL for the dark mode that needs to be a publicly accessible image. This image would be saved to the blob storage of Viva Learning for rendering within the Viva Learning app. Required.
func (m *LearningProvider) SetSquareLogoWebUrlForDarkTheme(value *string)() {
    err := m.GetBackingStore().Set("squareLogoWebUrlForDarkTheme", value)
    if err != nil {
        panic(err)
    }
}
// SetSquareLogoWebUrlForLightTheme sets the squareLogoWebUrlForLightTheme property value. The square logo URL for the light mode that needs to be a publicly accessible image. This image would be saved to the blob storage of Viva Learning for rendering within the Viva Learning app. Required.
func (m *LearningProvider) SetSquareLogoWebUrlForLightTheme(value *string)() {
    err := m.GetBackingStore().Set("squareLogoWebUrlForLightTheme", value)
    if err != nil {
        panic(err)
    }
}
type LearningProviderable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDisplayName()(*string)
    GetIsCourseActivitySyncEnabled()(*bool)
    GetLearningContents()([]LearningContentable)
    GetLearningCourseActivities()([]LearningCourseActivityable)
    GetLoginWebUrl()(*string)
    GetLongLogoWebUrlForDarkTheme()(*string)
    GetLongLogoWebUrlForLightTheme()(*string)
    GetSquareLogoWebUrlForDarkTheme()(*string)
    GetSquareLogoWebUrlForLightTheme()(*string)
    SetDisplayName(value *string)()
    SetIsCourseActivitySyncEnabled(value *bool)()
    SetLearningContents(value []LearningContentable)()
    SetLearningCourseActivities(value []LearningCourseActivityable)()
    SetLoginWebUrl(value *string)()
    SetLongLogoWebUrlForDarkTheme(value *string)()
    SetLongLogoWebUrlForLightTheme(value *string)()
    SetSquareLogoWebUrlForDarkTheme(value *string)()
    SetSquareLogoWebUrlForLightTheme(value *string)()
}
