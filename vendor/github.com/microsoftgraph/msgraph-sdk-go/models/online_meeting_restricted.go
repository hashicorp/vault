package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type OnlineMeetingRestricted struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewOnlineMeetingRestricted instantiates a new OnlineMeetingRestricted and sets the default values.
func NewOnlineMeetingRestricted()(*OnlineMeetingRestricted) {
    m := &OnlineMeetingRestricted{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateOnlineMeetingRestrictedFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOnlineMeetingRestrictedFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOnlineMeetingRestricted(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *OnlineMeetingRestricted) GetAdditionalData()(map[string]any) {
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
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *OnlineMeetingRestricted) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetContentSharingDisabled gets the contentSharingDisabled property value. Specifies the reason shared content from this participant is disabled. Possible values are: watermarkProtection, unknownFutureValue.
// returns a *OnlineMeetingContentSharingDisabledReason when successful
func (m *OnlineMeetingRestricted) GetContentSharingDisabled()(*OnlineMeetingContentSharingDisabledReason) {
    val, err := m.GetBackingStore().Get("contentSharingDisabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*OnlineMeetingContentSharingDisabledReason)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OnlineMeetingRestricted) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["contentSharingDisabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseOnlineMeetingContentSharingDisabledReason)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContentSharingDisabled(val.(*OnlineMeetingContentSharingDisabledReason))
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
    res["videoDisabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseOnlineMeetingVideoDisabledReason)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVideoDisabled(val.(*OnlineMeetingVideoDisabledReason))
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *OnlineMeetingRestricted) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetVideoDisabled gets the videoDisabled property value. Specifies the reason video from this participant is disabled. Possible values are: watermarkProtection, unknownFutureValue.
// returns a *OnlineMeetingVideoDisabledReason when successful
func (m *OnlineMeetingRestricted) GetVideoDisabled()(*OnlineMeetingVideoDisabledReason) {
    val, err := m.GetBackingStore().Get("videoDisabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*OnlineMeetingVideoDisabledReason)
    }
    return nil
}
// Serialize serializes information the current object
func (m *OnlineMeetingRestricted) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetContentSharingDisabled() != nil {
        cast := (*m.GetContentSharingDisabled()).String()
        err := writer.WriteStringValue("contentSharingDisabled", &cast)
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
    if m.GetVideoDisabled() != nil {
        cast := (*m.GetVideoDisabled()).String()
        err := writer.WriteStringValue("videoDisabled", &cast)
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
func (m *OnlineMeetingRestricted) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *OnlineMeetingRestricted) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetContentSharingDisabled sets the contentSharingDisabled property value. Specifies the reason shared content from this participant is disabled. Possible values are: watermarkProtection, unknownFutureValue.
func (m *OnlineMeetingRestricted) SetContentSharingDisabled(value *OnlineMeetingContentSharingDisabledReason)() {
    err := m.GetBackingStore().Set("contentSharingDisabled", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *OnlineMeetingRestricted) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetVideoDisabled sets the videoDisabled property value. Specifies the reason video from this participant is disabled. Possible values are: watermarkProtection, unknownFutureValue.
func (m *OnlineMeetingRestricted) SetVideoDisabled(value *OnlineMeetingVideoDisabledReason)() {
    err := m.GetBackingStore().Set("videoDisabled", value)
    if err != nil {
        panic(err)
    }
}
type OnlineMeetingRestrictedable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetContentSharingDisabled()(*OnlineMeetingContentSharingDisabledReason)
    GetOdataType()(*string)
    GetVideoDisabled()(*OnlineMeetingVideoDisabledReason)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetContentSharingDisabled(value *OnlineMeetingContentSharingDisabledReason)()
    SetOdataType(value *string)()
    SetVideoDisabled(value *OnlineMeetingVideoDisabledReason)()
}
