package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type BroadcastMeetingSettings struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewBroadcastMeetingSettings instantiates a new BroadcastMeetingSettings and sets the default values.
func NewBroadcastMeetingSettings()(*BroadcastMeetingSettings) {
    m := &BroadcastMeetingSettings{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateBroadcastMeetingSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateBroadcastMeetingSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewBroadcastMeetingSettings(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *BroadcastMeetingSettings) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetAllowedAudience gets the allowedAudience property value. Defines who can join the Teams live event. Possible values are listed in the following table.
// returns a *BroadcastMeetingAudience when successful
func (m *BroadcastMeetingSettings) GetAllowedAudience()(*BroadcastMeetingAudience) {
    val, err := m.GetBackingStore().Get("allowedAudience")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*BroadcastMeetingAudience)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *BroadcastMeetingSettings) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCaptions gets the captions property value. Caption settings of a Teams live event.
// returns a BroadcastMeetingCaptionSettingsable when successful
func (m *BroadcastMeetingSettings) GetCaptions()(BroadcastMeetingCaptionSettingsable) {
    val, err := m.GetBackingStore().Get("captions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(BroadcastMeetingCaptionSettingsable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *BroadcastMeetingSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["allowedAudience"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseBroadcastMeetingAudience)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowedAudience(val.(*BroadcastMeetingAudience))
        }
        return nil
    }
    res["captions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateBroadcastMeetingCaptionSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCaptions(val.(BroadcastMeetingCaptionSettingsable))
        }
        return nil
    }
    res["isAttendeeReportEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsAttendeeReportEnabled(val)
        }
        return nil
    }
    res["isQuestionAndAnswerEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsQuestionAndAnswerEnabled(val)
        }
        return nil
    }
    res["isRecordingEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsRecordingEnabled(val)
        }
        return nil
    }
    res["isVideoOnDemandEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsVideoOnDemandEnabled(val)
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    return res
}
// GetIsAttendeeReportEnabled gets the isAttendeeReportEnabled property value. Indicates whether attendee report is enabled for this Teams live event. Default value is false.
// returns a *bool when successful
func (m *BroadcastMeetingSettings) GetIsAttendeeReportEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isAttendeeReportEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsQuestionAndAnswerEnabled gets the isQuestionAndAnswerEnabled property value. Indicates whether Q&A is enabled for this Teams live event. Default value is false.
// returns a *bool when successful
func (m *BroadcastMeetingSettings) GetIsQuestionAndAnswerEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isQuestionAndAnswerEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsRecordingEnabled gets the isRecordingEnabled property value. Indicates whether recording is enabled for this Teams live event. Default value is false.
// returns a *bool when successful
func (m *BroadcastMeetingSettings) GetIsRecordingEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isRecordingEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsVideoOnDemandEnabled gets the isVideoOnDemandEnabled property value. Indicates whether video on demand is enabled for this Teams live event. Default value is false.
// returns a *bool when successful
func (m *BroadcastMeetingSettings) GetIsVideoOnDemandEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isVideoOnDemandEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *BroadcastMeetingSettings) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *BroadcastMeetingSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAllowedAudience() != nil {
        cast := (*m.GetAllowedAudience()).String()
        err := writer.WriteStringValue("allowedAudience", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("captions", m.GetCaptions())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isAttendeeReportEnabled", m.GetIsAttendeeReportEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isQuestionAndAnswerEnabled", m.GetIsQuestionAndAnswerEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isRecordingEnabled", m.GetIsRecordingEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isVideoOnDemandEnabled", m.GetIsVideoOnDemandEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *BroadcastMeetingSettings) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowedAudience sets the allowedAudience property value. Defines who can join the Teams live event. Possible values are listed in the following table.
func (m *BroadcastMeetingSettings) SetAllowedAudience(value *BroadcastMeetingAudience)() {
    err := m.GetBackingStore().Set("allowedAudience", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *BroadcastMeetingSettings) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCaptions sets the captions property value. Caption settings of a Teams live event.
func (m *BroadcastMeetingSettings) SetCaptions(value BroadcastMeetingCaptionSettingsable)() {
    err := m.GetBackingStore().Set("captions", value)
    if err != nil {
        panic(err)
    }
}
// SetIsAttendeeReportEnabled sets the isAttendeeReportEnabled property value. Indicates whether attendee report is enabled for this Teams live event. Default value is false.
func (m *BroadcastMeetingSettings) SetIsAttendeeReportEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isAttendeeReportEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsQuestionAndAnswerEnabled sets the isQuestionAndAnswerEnabled property value. Indicates whether Q&A is enabled for this Teams live event. Default value is false.
func (m *BroadcastMeetingSettings) SetIsQuestionAndAnswerEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isQuestionAndAnswerEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsRecordingEnabled sets the isRecordingEnabled property value. Indicates whether recording is enabled for this Teams live event. Default value is false.
func (m *BroadcastMeetingSettings) SetIsRecordingEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isRecordingEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsVideoOnDemandEnabled sets the isVideoOnDemandEnabled property value. Indicates whether video on demand is enabled for this Teams live event. Default value is false.
func (m *BroadcastMeetingSettings) SetIsVideoOnDemandEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isVideoOnDemandEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *BroadcastMeetingSettings) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type BroadcastMeetingSettingsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowedAudience()(*BroadcastMeetingAudience)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCaptions()(BroadcastMeetingCaptionSettingsable)
    GetIsAttendeeReportEnabled()(*bool)
    GetIsQuestionAndAnswerEnabled()(*bool)
    GetIsRecordingEnabled()(*bool)
    GetIsVideoOnDemandEnabled()(*bool)
    GetOdataType()(*string)
    SetAllowedAudience(value *BroadcastMeetingAudience)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCaptions(value BroadcastMeetingCaptionSettingsable)()
    SetIsAttendeeReportEnabled(value *bool)()
    SetIsQuestionAndAnswerEnabled(value *bool)()
    SetIsRecordingEnabled(value *bool)()
    SetIsVideoOnDemandEnabled(value *bool)()
    SetOdataType(value *string)()
}
