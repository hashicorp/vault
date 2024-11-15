package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type LearningContent struct {
    Entity
}
// NewLearningContent instantiates a new LearningContent and sets the default values.
func NewLearningContent()(*LearningContent) {
    m := &LearningContent{
        Entity: *NewEntity(),
    }
    return m
}
// CreateLearningContentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateLearningContentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewLearningContent(), nil
}
// GetAdditionalTags gets the additionalTags property value. Keywords, topics, and other tags associated with the learning content. Optional.
// returns a []string when successful
func (m *LearningContent) GetAdditionalTags()([]string) {
    val, err := m.GetBackingStore().Get("additionalTags")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetContentWebUrl gets the contentWebUrl property value. The content web URL for the learning content. Required.
// returns a *string when successful
func (m *LearningContent) GetContentWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("contentWebUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetContributors gets the contributors property value. The authors, creators, or contributors of the learning content. Optional.
// returns a []string when successful
func (m *LearningContent) GetContributors()([]string) {
    val, err := m.GetBackingStore().Get("contributors")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The date and time when the learning content was created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Optional.
// returns a *Time when successful
func (m *LearningContent) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDescription gets the description property value. The description or summary for the learning content. Optional.
// returns a *string when successful
func (m *LearningContent) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDuration gets the duration property value. The duration of the learning content in seconds. The value is represented in ISO 8601 format for durations. Optional.
// returns a *ISODuration when successful
func (m *LearningContent) GetDuration()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("duration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetExternalId gets the externalId property value. Unique external content ID for the learning content. Required.
// returns a *string when successful
func (m *LearningContent) GetExternalId()(*string) {
    val, err := m.GetBackingStore().Get("externalId")
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
func (m *LearningContent) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["additionalTags"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetAdditionalTags(res)
        }
        return nil
    }
    res["contentWebUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContentWebUrl(val)
        }
        return nil
    }
    res["contributors"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetContributors(res)
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
    res["duration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDuration(val)
        }
        return nil
    }
    res["externalId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalId(val)
        }
        return nil
    }
    res["format"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFormat(val)
        }
        return nil
    }
    res["isActive"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsActive(val)
        }
        return nil
    }
    res["isPremium"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsPremium(val)
        }
        return nil
    }
    res["isSearchable"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSearchable(val)
        }
        return nil
    }
    res["languageTag"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLanguageTag(val)
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
    res["level"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseLevel)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLevel(val.(*Level))
        }
        return nil
    }
    res["numberOfPages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNumberOfPages(val)
        }
        return nil
    }
    res["skillTags"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSkillTags(res)
        }
        return nil
    }
    res["sourceName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSourceName(val)
        }
        return nil
    }
    res["thumbnailWebUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetThumbnailWebUrl(val)
        }
        return nil
    }
    res["title"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTitle(val)
        }
        return nil
    }
    return res
}
// GetFormat gets the format property value. The format of the learning content. For example, Course, Video, Book, Book Summary, Audiobook Summary. Optional.
// returns a *string when successful
func (m *LearningContent) GetFormat()(*string) {
    val, err := m.GetBackingStore().Get("format")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIsActive gets the isActive property value. Indicates whether the content is active or not. Inactive content doesn't show up in the UI. The default value is true. Optional.
// returns a *bool when successful
func (m *LearningContent) GetIsActive()(*bool) {
    val, err := m.GetBackingStore().Get("isActive")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsPremium gets the isPremium property value. Indicates whether the learning content requires the user to sign-in on the learning provider platform or not. The default value is false. Optional.
// returns a *bool when successful
func (m *LearningContent) GetIsPremium()(*bool) {
    val, err := m.GetBackingStore().Get("isPremium")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsSearchable gets the isSearchable property value. Indicates whether the learning content is searchable or not. The default value is true. Optional.
// returns a *bool when successful
func (m *LearningContent) GetIsSearchable()(*bool) {
    val, err := m.GetBackingStore().Get("isSearchable")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLanguageTag gets the languageTag property value. The language of the learning content, for example, en-us or fr-fr. Required.
// returns a *string when successful
func (m *LearningContent) GetLanguageTag()(*string) {
    val, err := m.GetBackingStore().Get("languageTag")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. The date and time when the learning content was last modified. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Optional.
// returns a *Time when successful
func (m *LearningContent) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetLevel gets the level property value. The difficulty level of the learning content. Possible values are: Beginner, Intermediate, Advanced, unknownFutureValue. Optional.
// returns a *Level when successful
func (m *LearningContent) GetLevel()(*Level) {
    val, err := m.GetBackingStore().Get("level")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Level)
    }
    return nil
}
// GetNumberOfPages gets the numberOfPages property value. The number of pages of the learning content, for example, 9. Optional.
// returns a *int32 when successful
func (m *LearningContent) GetNumberOfPages()(*int32) {
    val, err := m.GetBackingStore().Get("numberOfPages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSkillTags gets the skillTags property value. The skills tags associated with the learning content. Optional.
// returns a []string when successful
func (m *LearningContent) GetSkillTags()([]string) {
    val, err := m.GetBackingStore().Get("skillTags")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetSourceName gets the sourceName property value. The source name of the learning content, such as LinkedIn Learning or Coursera. Optional.
// returns a *string when successful
func (m *LearningContent) GetSourceName()(*string) {
    val, err := m.GetBackingStore().Get("sourceName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetThumbnailWebUrl gets the thumbnailWebUrl property value. The URL of learning content thumbnail image. Optional.
// returns a *string when successful
func (m *LearningContent) GetThumbnailWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("thumbnailWebUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTitle gets the title property value. The title of the learning content. Required.
// returns a *string when successful
func (m *LearningContent) GetTitle()(*string) {
    val, err := m.GetBackingStore().Get("title")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *LearningContent) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAdditionalTags() != nil {
        err = writer.WriteCollectionOfStringValues("additionalTags", m.GetAdditionalTags())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("contentWebUrl", m.GetContentWebUrl())
        if err != nil {
            return err
        }
    }
    if m.GetContributors() != nil {
        err = writer.WriteCollectionOfStringValues("contributors", m.GetContributors())
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
    {
        err = writer.WriteISODurationValue("duration", m.GetDuration())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("externalId", m.GetExternalId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("format", m.GetFormat())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isActive", m.GetIsActive())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isPremium", m.GetIsPremium())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isSearchable", m.GetIsSearchable())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("languageTag", m.GetLanguageTag())
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
    if m.GetLevel() != nil {
        cast := (*m.GetLevel()).String()
        err = writer.WriteStringValue("level", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("numberOfPages", m.GetNumberOfPages())
        if err != nil {
            return err
        }
    }
    if m.GetSkillTags() != nil {
        err = writer.WriteCollectionOfStringValues("skillTags", m.GetSkillTags())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("sourceName", m.GetSourceName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("thumbnailWebUrl", m.GetThumbnailWebUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("title", m.GetTitle())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalTags sets the additionalTags property value. Keywords, topics, and other tags associated with the learning content. Optional.
func (m *LearningContent) SetAdditionalTags(value []string)() {
    err := m.GetBackingStore().Set("additionalTags", value)
    if err != nil {
        panic(err)
    }
}
// SetContentWebUrl sets the contentWebUrl property value. The content web URL for the learning content. Required.
func (m *LearningContent) SetContentWebUrl(value *string)() {
    err := m.GetBackingStore().Set("contentWebUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetContributors sets the contributors property value. The authors, creators, or contributors of the learning content. Optional.
func (m *LearningContent) SetContributors(value []string)() {
    err := m.GetBackingStore().Set("contributors", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The date and time when the learning content was created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Optional.
func (m *LearningContent) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. The description or summary for the learning content. Optional.
func (m *LearningContent) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDuration sets the duration property value. The duration of the learning content in seconds. The value is represented in ISO 8601 format for durations. Optional.
func (m *LearningContent) SetDuration(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("duration", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalId sets the externalId property value. Unique external content ID for the learning content. Required.
func (m *LearningContent) SetExternalId(value *string)() {
    err := m.GetBackingStore().Set("externalId", value)
    if err != nil {
        panic(err)
    }
}
// SetFormat sets the format property value. The format of the learning content. For example, Course, Video, Book, Book Summary, Audiobook Summary. Optional.
func (m *LearningContent) SetFormat(value *string)() {
    err := m.GetBackingStore().Set("format", value)
    if err != nil {
        panic(err)
    }
}
// SetIsActive sets the isActive property value. Indicates whether the content is active or not. Inactive content doesn't show up in the UI. The default value is true. Optional.
func (m *LearningContent) SetIsActive(value *bool)() {
    err := m.GetBackingStore().Set("isActive", value)
    if err != nil {
        panic(err)
    }
}
// SetIsPremium sets the isPremium property value. Indicates whether the learning content requires the user to sign-in on the learning provider platform or not. The default value is false. Optional.
func (m *LearningContent) SetIsPremium(value *bool)() {
    err := m.GetBackingStore().Set("isPremium", value)
    if err != nil {
        panic(err)
    }
}
// SetIsSearchable sets the isSearchable property value. Indicates whether the learning content is searchable or not. The default value is true. Optional.
func (m *LearningContent) SetIsSearchable(value *bool)() {
    err := m.GetBackingStore().Set("isSearchable", value)
    if err != nil {
        panic(err)
    }
}
// SetLanguageTag sets the languageTag property value. The language of the learning content, for example, en-us or fr-fr. Required.
func (m *LearningContent) SetLanguageTag(value *string)() {
    err := m.GetBackingStore().Set("languageTag", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. The date and time when the learning content was last modified. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Optional.
func (m *LearningContent) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLevel sets the level property value. The difficulty level of the learning content. Possible values are: Beginner, Intermediate, Advanced, unknownFutureValue. Optional.
func (m *LearningContent) SetLevel(value *Level)() {
    err := m.GetBackingStore().Set("level", value)
    if err != nil {
        panic(err)
    }
}
// SetNumberOfPages sets the numberOfPages property value. The number of pages of the learning content, for example, 9. Optional.
func (m *LearningContent) SetNumberOfPages(value *int32)() {
    err := m.GetBackingStore().Set("numberOfPages", value)
    if err != nil {
        panic(err)
    }
}
// SetSkillTags sets the skillTags property value. The skills tags associated with the learning content. Optional.
func (m *LearningContent) SetSkillTags(value []string)() {
    err := m.GetBackingStore().Set("skillTags", value)
    if err != nil {
        panic(err)
    }
}
// SetSourceName sets the sourceName property value. The source name of the learning content, such as LinkedIn Learning or Coursera. Optional.
func (m *LearningContent) SetSourceName(value *string)() {
    err := m.GetBackingStore().Set("sourceName", value)
    if err != nil {
        panic(err)
    }
}
// SetThumbnailWebUrl sets the thumbnailWebUrl property value. The URL of learning content thumbnail image. Optional.
func (m *LearningContent) SetThumbnailWebUrl(value *string)() {
    err := m.GetBackingStore().Set("thumbnailWebUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetTitle sets the title property value. The title of the learning content. Required.
func (m *LearningContent) SetTitle(value *string)() {
    err := m.GetBackingStore().Set("title", value)
    if err != nil {
        panic(err)
    }
}
type LearningContentable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAdditionalTags()([]string)
    GetContentWebUrl()(*string)
    GetContributors()([]string)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDescription()(*string)
    GetDuration()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetExternalId()(*string)
    GetFormat()(*string)
    GetIsActive()(*bool)
    GetIsPremium()(*bool)
    GetIsSearchable()(*bool)
    GetLanguageTag()(*string)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLevel()(*Level)
    GetNumberOfPages()(*int32)
    GetSkillTags()([]string)
    GetSourceName()(*string)
    GetThumbnailWebUrl()(*string)
    GetTitle()(*string)
    SetAdditionalTags(value []string)()
    SetContentWebUrl(value *string)()
    SetContributors(value []string)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDescription(value *string)()
    SetDuration(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetExternalId(value *string)()
    SetFormat(value *string)()
    SetIsActive(value *bool)()
    SetIsPremium(value *bool)()
    SetIsSearchable(value *bool)()
    SetLanguageTag(value *string)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLevel(value *Level)()
    SetNumberOfPages(value *int32)()
    SetSkillTags(value []string)()
    SetSourceName(value *string)()
    SetThumbnailWebUrl(value *string)()
    SetTitle(value *string)()
}
