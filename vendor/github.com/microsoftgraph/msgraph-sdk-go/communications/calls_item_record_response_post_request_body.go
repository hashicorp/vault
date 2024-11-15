package communications

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type CallsItemRecordResponsePostRequestBody struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewCallsItemRecordResponsePostRequestBody instantiates a new CallsItemRecordResponsePostRequestBody and sets the default values.
func NewCallsItemRecordResponsePostRequestBody()(*CallsItemRecordResponsePostRequestBody) {
    m := &CallsItemRecordResponsePostRequestBody{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateCallsItemRecordResponsePostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCallsItemRecordResponsePostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCallsItemRecordResponsePostRequestBody(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *CallsItemRecordResponsePostRequestBody) GetAdditionalData()(map[string]any) {
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
func (m *CallsItemRecordResponsePostRequestBody) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetBargeInAllowed gets the bargeInAllowed property value. The bargeInAllowed property
// returns a *bool when successful
func (m *CallsItemRecordResponsePostRequestBody) GetBargeInAllowed()(*bool) {
    val, err := m.GetBackingStore().Get("bargeInAllowed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetClientContext gets the clientContext property value. The clientContext property
// returns a *string when successful
func (m *CallsItemRecordResponsePostRequestBody) GetClientContext()(*string) {
    val, err := m.GetBackingStore().Get("clientContext")
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
func (m *CallsItemRecordResponsePostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["bargeInAllowed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBargeInAllowed(val)
        }
        return nil
    }
    res["clientContext"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClientContext(val)
        }
        return nil
    }
    res["initialSilenceTimeoutInSeconds"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInitialSilenceTimeoutInSeconds(val)
        }
        return nil
    }
    res["maxRecordDurationInSeconds"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMaxRecordDurationInSeconds(val)
        }
        return nil
    }
    res["maxSilenceTimeoutInSeconds"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMaxSilenceTimeoutInSeconds(val)
        }
        return nil
    }
    res["playBeep"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPlayBeep(val)
        }
        return nil
    }
    res["prompts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePromptFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Promptable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Promptable)
                }
            }
            m.SetPrompts(res)
        }
        return nil
    }
    res["stopTones"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetStopTones(res)
        }
        return nil
    }
    return res
}
// GetInitialSilenceTimeoutInSeconds gets the initialSilenceTimeoutInSeconds property value. The initialSilenceTimeoutInSeconds property
// returns a *int32 when successful
func (m *CallsItemRecordResponsePostRequestBody) GetInitialSilenceTimeoutInSeconds()(*int32) {
    val, err := m.GetBackingStore().Get("initialSilenceTimeoutInSeconds")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetMaxRecordDurationInSeconds gets the maxRecordDurationInSeconds property value. The maxRecordDurationInSeconds property
// returns a *int32 when successful
func (m *CallsItemRecordResponsePostRequestBody) GetMaxRecordDurationInSeconds()(*int32) {
    val, err := m.GetBackingStore().Get("maxRecordDurationInSeconds")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetMaxSilenceTimeoutInSeconds gets the maxSilenceTimeoutInSeconds property value. The maxSilenceTimeoutInSeconds property
// returns a *int32 when successful
func (m *CallsItemRecordResponsePostRequestBody) GetMaxSilenceTimeoutInSeconds()(*int32) {
    val, err := m.GetBackingStore().Get("maxSilenceTimeoutInSeconds")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPlayBeep gets the playBeep property value. The playBeep property
// returns a *bool when successful
func (m *CallsItemRecordResponsePostRequestBody) GetPlayBeep()(*bool) {
    val, err := m.GetBackingStore().Get("playBeep")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPrompts gets the prompts property value. The prompts property
// returns a []Promptable when successful
func (m *CallsItemRecordResponsePostRequestBody) GetPrompts()([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Promptable) {
    val, err := m.GetBackingStore().Get("prompts")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Promptable)
    }
    return nil
}
// GetStopTones gets the stopTones property value. The stopTones property
// returns a []string when successful
func (m *CallsItemRecordResponsePostRequestBody) GetStopTones()([]string) {
    val, err := m.GetBackingStore().Get("stopTones")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CallsItemRecordResponsePostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("bargeInAllowed", m.GetBargeInAllowed())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("clientContext", m.GetClientContext())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("initialSilenceTimeoutInSeconds", m.GetInitialSilenceTimeoutInSeconds())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("maxRecordDurationInSeconds", m.GetMaxRecordDurationInSeconds())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("maxSilenceTimeoutInSeconds", m.GetMaxSilenceTimeoutInSeconds())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("playBeep", m.GetPlayBeep())
        if err != nil {
            return err
        }
    }
    if m.GetPrompts() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPrompts()))
        for i, v := range m.GetPrompts() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("prompts", cast)
        if err != nil {
            return err
        }
    }
    if m.GetStopTones() != nil {
        err := writer.WriteCollectionOfStringValues("stopTones", m.GetStopTones())
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
func (m *CallsItemRecordResponsePostRequestBody) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *CallsItemRecordResponsePostRequestBody) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetBargeInAllowed sets the bargeInAllowed property value. The bargeInAllowed property
func (m *CallsItemRecordResponsePostRequestBody) SetBargeInAllowed(value *bool)() {
    err := m.GetBackingStore().Set("bargeInAllowed", value)
    if err != nil {
        panic(err)
    }
}
// SetClientContext sets the clientContext property value. The clientContext property
func (m *CallsItemRecordResponsePostRequestBody) SetClientContext(value *string)() {
    err := m.GetBackingStore().Set("clientContext", value)
    if err != nil {
        panic(err)
    }
}
// SetInitialSilenceTimeoutInSeconds sets the initialSilenceTimeoutInSeconds property value. The initialSilenceTimeoutInSeconds property
func (m *CallsItemRecordResponsePostRequestBody) SetInitialSilenceTimeoutInSeconds(value *int32)() {
    err := m.GetBackingStore().Set("initialSilenceTimeoutInSeconds", value)
    if err != nil {
        panic(err)
    }
}
// SetMaxRecordDurationInSeconds sets the maxRecordDurationInSeconds property value. The maxRecordDurationInSeconds property
func (m *CallsItemRecordResponsePostRequestBody) SetMaxRecordDurationInSeconds(value *int32)() {
    err := m.GetBackingStore().Set("maxRecordDurationInSeconds", value)
    if err != nil {
        panic(err)
    }
}
// SetMaxSilenceTimeoutInSeconds sets the maxSilenceTimeoutInSeconds property value. The maxSilenceTimeoutInSeconds property
func (m *CallsItemRecordResponsePostRequestBody) SetMaxSilenceTimeoutInSeconds(value *int32)() {
    err := m.GetBackingStore().Set("maxSilenceTimeoutInSeconds", value)
    if err != nil {
        panic(err)
    }
}
// SetPlayBeep sets the playBeep property value. The playBeep property
func (m *CallsItemRecordResponsePostRequestBody) SetPlayBeep(value *bool)() {
    err := m.GetBackingStore().Set("playBeep", value)
    if err != nil {
        panic(err)
    }
}
// SetPrompts sets the prompts property value. The prompts property
func (m *CallsItemRecordResponsePostRequestBody) SetPrompts(value []iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Promptable)() {
    err := m.GetBackingStore().Set("prompts", value)
    if err != nil {
        panic(err)
    }
}
// SetStopTones sets the stopTones property value. The stopTones property
func (m *CallsItemRecordResponsePostRequestBody) SetStopTones(value []string)() {
    err := m.GetBackingStore().Set("stopTones", value)
    if err != nil {
        panic(err)
    }
}
type CallsItemRecordResponsePostRequestBodyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetBargeInAllowed()(*bool)
    GetClientContext()(*string)
    GetInitialSilenceTimeoutInSeconds()(*int32)
    GetMaxRecordDurationInSeconds()(*int32)
    GetMaxSilenceTimeoutInSeconds()(*int32)
    GetPlayBeep()(*bool)
    GetPrompts()([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Promptable)
    GetStopTones()([]string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetBargeInAllowed(value *bool)()
    SetClientContext(value *string)()
    SetInitialSilenceTimeoutInSeconds(value *int32)()
    SetMaxRecordDurationInSeconds(value *int32)()
    SetMaxSilenceTimeoutInSeconds(value *int32)()
    SetPlayBeep(value *bool)()
    SetPrompts(value []iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Promptable)()
    SetStopTones(value []string)()
}
