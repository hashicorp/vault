package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type MessageSecurityState struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewMessageSecurityState instantiates a new MessageSecurityState and sets the default values.
func NewMessageSecurityState()(*MessageSecurityState) {
    m := &MessageSecurityState{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateMessageSecurityStateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMessageSecurityStateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMessageSecurityState(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *MessageSecurityState) GetAdditionalData()(map[string]any) {
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
func (m *MessageSecurityState) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetConnectingIP gets the connectingIP property value. The connectingIP property
// returns a *string when successful
func (m *MessageSecurityState) GetConnectingIP()(*string) {
    val, err := m.GetBackingStore().Get("connectingIP")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeliveryAction gets the deliveryAction property value. The deliveryAction property
// returns a *string when successful
func (m *MessageSecurityState) GetDeliveryAction()(*string) {
    val, err := m.GetBackingStore().Get("deliveryAction")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeliveryLocation gets the deliveryLocation property value. The deliveryLocation property
// returns a *string when successful
func (m *MessageSecurityState) GetDeliveryLocation()(*string) {
    val, err := m.GetBackingStore().Get("deliveryLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDirectionality gets the directionality property value. The directionality property
// returns a *string when successful
func (m *MessageSecurityState) GetDirectionality()(*string) {
    val, err := m.GetBackingStore().Get("directionality")
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
func (m *MessageSecurityState) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["connectingIP"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConnectingIP(val)
        }
        return nil
    }
    res["deliveryAction"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeliveryAction(val)
        }
        return nil
    }
    res["deliveryLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeliveryLocation(val)
        }
        return nil
    }
    res["directionality"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDirectionality(val)
        }
        return nil
    }
    res["internetMessageId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInternetMessageId(val)
        }
        return nil
    }
    res["messageFingerprint"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMessageFingerprint(val)
        }
        return nil
    }
    res["messageReceivedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMessageReceivedDateTime(val)
        }
        return nil
    }
    res["messageSubject"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMessageSubject(val)
        }
        return nil
    }
    res["networkMessageId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNetworkMessageId(val)
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
// GetInternetMessageId gets the internetMessageId property value. The internetMessageId property
// returns a *string when successful
func (m *MessageSecurityState) GetInternetMessageId()(*string) {
    val, err := m.GetBackingStore().Get("internetMessageId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMessageFingerprint gets the messageFingerprint property value. The messageFingerprint property
// returns a *string when successful
func (m *MessageSecurityState) GetMessageFingerprint()(*string) {
    val, err := m.GetBackingStore().Get("messageFingerprint")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMessageReceivedDateTime gets the messageReceivedDateTime property value. The messageReceivedDateTime property
// returns a *Time when successful
func (m *MessageSecurityState) GetMessageReceivedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("messageReceivedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetMessageSubject gets the messageSubject property value. The messageSubject property
// returns a *string when successful
func (m *MessageSecurityState) GetMessageSubject()(*string) {
    val, err := m.GetBackingStore().Get("messageSubject")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNetworkMessageId gets the networkMessageId property value. The networkMessageId property
// returns a *string when successful
func (m *MessageSecurityState) GetNetworkMessageId()(*string) {
    val, err := m.GetBackingStore().Get("networkMessageId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *MessageSecurityState) GetOdataType()(*string) {
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
func (m *MessageSecurityState) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("connectingIP", m.GetConnectingIP())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("deliveryAction", m.GetDeliveryAction())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("deliveryLocation", m.GetDeliveryLocation())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("directionality", m.GetDirectionality())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("internetMessageId", m.GetInternetMessageId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("messageFingerprint", m.GetMessageFingerprint())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("messageReceivedDateTime", m.GetMessageReceivedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("messageSubject", m.GetMessageSubject())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("networkMessageId", m.GetNetworkMessageId())
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
func (m *MessageSecurityState) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *MessageSecurityState) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetConnectingIP sets the connectingIP property value. The connectingIP property
func (m *MessageSecurityState) SetConnectingIP(value *string)() {
    err := m.GetBackingStore().Set("connectingIP", value)
    if err != nil {
        panic(err)
    }
}
// SetDeliveryAction sets the deliveryAction property value. The deliveryAction property
func (m *MessageSecurityState) SetDeliveryAction(value *string)() {
    err := m.GetBackingStore().Set("deliveryAction", value)
    if err != nil {
        panic(err)
    }
}
// SetDeliveryLocation sets the deliveryLocation property value. The deliveryLocation property
func (m *MessageSecurityState) SetDeliveryLocation(value *string)() {
    err := m.GetBackingStore().Set("deliveryLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetDirectionality sets the directionality property value. The directionality property
func (m *MessageSecurityState) SetDirectionality(value *string)() {
    err := m.GetBackingStore().Set("directionality", value)
    if err != nil {
        panic(err)
    }
}
// SetInternetMessageId sets the internetMessageId property value. The internetMessageId property
func (m *MessageSecurityState) SetInternetMessageId(value *string)() {
    err := m.GetBackingStore().Set("internetMessageId", value)
    if err != nil {
        panic(err)
    }
}
// SetMessageFingerprint sets the messageFingerprint property value. The messageFingerprint property
func (m *MessageSecurityState) SetMessageFingerprint(value *string)() {
    err := m.GetBackingStore().Set("messageFingerprint", value)
    if err != nil {
        panic(err)
    }
}
// SetMessageReceivedDateTime sets the messageReceivedDateTime property value. The messageReceivedDateTime property
func (m *MessageSecurityState) SetMessageReceivedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("messageReceivedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetMessageSubject sets the messageSubject property value. The messageSubject property
func (m *MessageSecurityState) SetMessageSubject(value *string)() {
    err := m.GetBackingStore().Set("messageSubject", value)
    if err != nil {
        panic(err)
    }
}
// SetNetworkMessageId sets the networkMessageId property value. The networkMessageId property
func (m *MessageSecurityState) SetNetworkMessageId(value *string)() {
    err := m.GetBackingStore().Set("networkMessageId", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *MessageSecurityState) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type MessageSecurityStateable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetConnectingIP()(*string)
    GetDeliveryAction()(*string)
    GetDeliveryLocation()(*string)
    GetDirectionality()(*string)
    GetInternetMessageId()(*string)
    GetMessageFingerprint()(*string)
    GetMessageReceivedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetMessageSubject()(*string)
    GetNetworkMessageId()(*string)
    GetOdataType()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetConnectingIP(value *string)()
    SetDeliveryAction(value *string)()
    SetDeliveryLocation(value *string)()
    SetDirectionality(value *string)()
    SetInternetMessageId(value *string)()
    SetMessageFingerprint(value *string)()
    SetMessageReceivedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetMessageSubject(value *string)()
    SetNetworkMessageId(value *string)()
    SetOdataType(value *string)()
}
