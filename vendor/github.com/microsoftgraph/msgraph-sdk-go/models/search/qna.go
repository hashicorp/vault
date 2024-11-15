package search

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type Qna struct {
    SearchAnswer
}
// NewQna instantiates a new Qna and sets the default values.
func NewQna()(*Qna) {
    m := &Qna{
        SearchAnswer: *NewSearchAnswer(),
    }
    return m
}
// CreateQnaFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateQnaFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewQna(), nil
}
// GetAvailabilityEndDateTime gets the availabilityEndDateTime property value. Date and time when the QnA stops appearing as a search result. Set as null for always available. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *Qna) GetAvailabilityEndDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("availabilityEndDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetAvailabilityStartDateTime gets the availabilityStartDateTime property value. Date and time when the QnA starts to appear as a search result. Set as null for always available. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *Qna) GetAvailabilityStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("availabilityStartDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Qna) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.SearchAnswer.GetFieldDeserializers()
    res["availabilityEndDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAvailabilityEndDateTime(val)
        }
        return nil
    }
    res["availabilityStartDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAvailabilityStartDateTime(val)
        }
        return nil
    }
    res["groupIds"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetGroupIds(res)
        }
        return nil
    }
    res["isSuggested"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSuggested(val)
        }
        return nil
    }
    res["keywords"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAnswerKeywordFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKeywords(val.(AnswerKeywordable))
        }
        return nil
    }
    res["languageTags"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetLanguageTags(res)
        }
        return nil
    }
    res["platforms"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ParseDevicePlatformType)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DevicePlatformType, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DevicePlatformType))
                }
            }
            m.SetPlatforms(res)
        }
        return nil
    }
    res["state"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAnswerState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetState(val.(*AnswerState))
        }
        return nil
    }
    res["targetedVariations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAnswerVariantFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AnswerVariantable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AnswerVariantable)
                }
            }
            m.SetTargetedVariations(res)
        }
        return nil
    }
    return res
}
// GetGroupIds gets the groupIds property value. The list of security groups that are able to view this QnA.
// returns a []string when successful
func (m *Qna) GetGroupIds()([]string) {
    val, err := m.GetBackingStore().Get("groupIds")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetIsSuggested gets the isSuggested property value. True if a user or Microsoft suggested this QnA to the admin. Read-only.
// returns a *bool when successful
func (m *Qna) GetIsSuggested()(*bool) {
    val, err := m.GetBackingStore().Get("isSuggested")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKeywords gets the keywords property value. Keywords that trigger this QnA to appear in search results.
// returns a AnswerKeywordable when successful
func (m *Qna) GetKeywords()(AnswerKeywordable) {
    val, err := m.GetBackingStore().Get("keywords")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AnswerKeywordable)
    }
    return nil
}
// GetLanguageTags gets the languageTags property value. A list of geographically specific language names in which this QnA can be viewed. Each language tag value follows the pattern {language}-{region}. For example, en-us is English as used in the United States. For the list of possible values, see Supported language tags.
// returns a []string when successful
func (m *Qna) GetLanguageTags()([]string) {
    val, err := m.GetBackingStore().Get("languageTags")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetPlatforms gets the platforms property value. List of devices and operating systems that are able to view this QnA. Possible values are: android, androidForWork, ios, macOS, windowsPhone81, windowsPhone81AndLater, windows10AndLater, androidWorkProfile, unknown, androidASOP, androidMobileApplicationManagement, iOSMobileApplicationManagement, unknownFutureValue.
// returns a []DevicePlatformType when successful
func (m *Qna) GetPlatforms()([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DevicePlatformType) {
    val, err := m.GetBackingStore().Get("platforms")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DevicePlatformType)
    }
    return nil
}
// GetState gets the state property value. The state property
// returns a *AnswerState when successful
func (m *Qna) GetState()(*AnswerState) {
    val, err := m.GetBackingStore().Get("state")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AnswerState)
    }
    return nil
}
// GetTargetedVariations gets the targetedVariations property value. Variations of a QnA for different countries or devices. Use when you need to show different content to users based on their device, country/region, or both. The date and group settings apply to all variations.
// returns a []AnswerVariantable when successful
func (m *Qna) GetTargetedVariations()([]AnswerVariantable) {
    val, err := m.GetBackingStore().Get("targetedVariations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AnswerVariantable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Qna) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.SearchAnswer.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("availabilityEndDateTime", m.GetAvailabilityEndDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("availabilityStartDateTime", m.GetAvailabilityStartDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetGroupIds() != nil {
        err = writer.WriteCollectionOfStringValues("groupIds", m.GetGroupIds())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isSuggested", m.GetIsSuggested())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("keywords", m.GetKeywords())
        if err != nil {
            return err
        }
    }
    if m.GetLanguageTags() != nil {
        err = writer.WriteCollectionOfStringValues("languageTags", m.GetLanguageTags())
        if err != nil {
            return err
        }
    }
    if m.GetPlatforms() != nil {
        err = writer.WriteCollectionOfStringValues("platforms", iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SerializeDevicePlatformType(m.GetPlatforms()))
        if err != nil {
            return err
        }
    }
    if m.GetState() != nil {
        cast := (*m.GetState()).String()
        err = writer.WriteStringValue("state", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetTargetedVariations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTargetedVariations()))
        for i, v := range m.GetTargetedVariations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("targetedVariations", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAvailabilityEndDateTime sets the availabilityEndDateTime property value. Date and time when the QnA stops appearing as a search result. Set as null for always available. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *Qna) SetAvailabilityEndDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("availabilityEndDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetAvailabilityStartDateTime sets the availabilityStartDateTime property value. Date and time when the QnA starts to appear as a search result. Set as null for always available. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *Qna) SetAvailabilityStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("availabilityStartDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetGroupIds sets the groupIds property value. The list of security groups that are able to view this QnA.
func (m *Qna) SetGroupIds(value []string)() {
    err := m.GetBackingStore().Set("groupIds", value)
    if err != nil {
        panic(err)
    }
}
// SetIsSuggested sets the isSuggested property value. True if a user or Microsoft suggested this QnA to the admin. Read-only.
func (m *Qna) SetIsSuggested(value *bool)() {
    err := m.GetBackingStore().Set("isSuggested", value)
    if err != nil {
        panic(err)
    }
}
// SetKeywords sets the keywords property value. Keywords that trigger this QnA to appear in search results.
func (m *Qna) SetKeywords(value AnswerKeywordable)() {
    err := m.GetBackingStore().Set("keywords", value)
    if err != nil {
        panic(err)
    }
}
// SetLanguageTags sets the languageTags property value. A list of geographically specific language names in which this QnA can be viewed. Each language tag value follows the pattern {language}-{region}. For example, en-us is English as used in the United States. For the list of possible values, see Supported language tags.
func (m *Qna) SetLanguageTags(value []string)() {
    err := m.GetBackingStore().Set("languageTags", value)
    if err != nil {
        panic(err)
    }
}
// SetPlatforms sets the platforms property value. List of devices and operating systems that are able to view this QnA. Possible values are: android, androidForWork, ios, macOS, windowsPhone81, windowsPhone81AndLater, windows10AndLater, androidWorkProfile, unknown, androidASOP, androidMobileApplicationManagement, iOSMobileApplicationManagement, unknownFutureValue.
func (m *Qna) SetPlatforms(value []iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DevicePlatformType)() {
    err := m.GetBackingStore().Set("platforms", value)
    if err != nil {
        panic(err)
    }
}
// SetState sets the state property value. The state property
func (m *Qna) SetState(value *AnswerState)() {
    err := m.GetBackingStore().Set("state", value)
    if err != nil {
        panic(err)
    }
}
// SetTargetedVariations sets the targetedVariations property value. Variations of a QnA for different countries or devices. Use when you need to show different content to users based on their device, country/region, or both. The date and group settings apply to all variations.
func (m *Qna) SetTargetedVariations(value []AnswerVariantable)() {
    err := m.GetBackingStore().Set("targetedVariations", value)
    if err != nil {
        panic(err)
    }
}
type Qnaable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    SearchAnswerable
    GetAvailabilityEndDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetAvailabilityStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetGroupIds()([]string)
    GetIsSuggested()(*bool)
    GetKeywords()(AnswerKeywordable)
    GetLanguageTags()([]string)
    GetPlatforms()([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DevicePlatformType)
    GetState()(*AnswerState)
    GetTargetedVariations()([]AnswerVariantable)
    SetAvailabilityEndDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetAvailabilityStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetGroupIds(value []string)()
    SetIsSuggested(value *bool)()
    SetKeywords(value AnswerKeywordable)()
    SetLanguageTags(value []string)()
    SetPlatforms(value []iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DevicePlatformType)()
    SetState(value *AnswerState)()
    SetTargetedVariations(value []AnswerVariantable)()
}
