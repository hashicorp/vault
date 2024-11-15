package drives

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type ItemItemsItemAssignSensitivityLabelPostRequestBody struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewItemItemsItemAssignSensitivityLabelPostRequestBody instantiates a new ItemItemsItemAssignSensitivityLabelPostRequestBody and sets the default values.
func NewItemItemsItemAssignSensitivityLabelPostRequestBody()(*ItemItemsItemAssignSensitivityLabelPostRequestBody) {
    m := &ItemItemsItemAssignSensitivityLabelPostRequestBody{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateItemItemsItemAssignSensitivityLabelPostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemItemsItemAssignSensitivityLabelPostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemItemsItemAssignSensitivityLabelPostRequestBody(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *ItemItemsItemAssignSensitivityLabelPostRequestBody) GetAdditionalData()(map[string]any) {
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
// GetAssignmentMethod gets the assignmentMethod property value. The assignmentMethod property
// returns a *SensitivityLabelAssignmentMethod when successful
func (m *ItemItemsItemAssignSensitivityLabelPostRequestBody) GetAssignmentMethod()(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SensitivityLabelAssignmentMethod) {
    val, err := m.GetBackingStore().Get("assignmentMethod")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SensitivityLabelAssignmentMethod)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *ItemItemsItemAssignSensitivityLabelPostRequestBody) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ItemItemsItemAssignSensitivityLabelPostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["assignmentMethod"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ParseSensitivityLabelAssignmentMethod)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignmentMethod(val.(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SensitivityLabelAssignmentMethod))
        }
        return nil
    }
    res["justificationText"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetJustificationText(val)
        }
        return nil
    }
    res["sensitivityLabelId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSensitivityLabelId(val)
        }
        return nil
    }
    return res
}
// GetJustificationText gets the justificationText property value. The justificationText property
// returns a *string when successful
func (m *ItemItemsItemAssignSensitivityLabelPostRequestBody) GetJustificationText()(*string) {
    val, err := m.GetBackingStore().Get("justificationText")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSensitivityLabelId gets the sensitivityLabelId property value. The sensitivityLabelId property
// returns a *string when successful
func (m *ItemItemsItemAssignSensitivityLabelPostRequestBody) GetSensitivityLabelId()(*string) {
    val, err := m.GetBackingStore().Get("sensitivityLabelId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ItemItemsItemAssignSensitivityLabelPostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAssignmentMethod() != nil {
        cast := (*m.GetAssignmentMethod()).String()
        err := writer.WriteStringValue("assignmentMethod", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("justificationText", m.GetJustificationText())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("sensitivityLabelId", m.GetSensitivityLabelId())
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
func (m *ItemItemsItemAssignSensitivityLabelPostRequestBody) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignmentMethod sets the assignmentMethod property value. The assignmentMethod property
func (m *ItemItemsItemAssignSensitivityLabelPostRequestBody) SetAssignmentMethod(value *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SensitivityLabelAssignmentMethod)() {
    err := m.GetBackingStore().Set("assignmentMethod", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *ItemItemsItemAssignSensitivityLabelPostRequestBody) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetJustificationText sets the justificationText property value. The justificationText property
func (m *ItemItemsItemAssignSensitivityLabelPostRequestBody) SetJustificationText(value *string)() {
    err := m.GetBackingStore().Set("justificationText", value)
    if err != nil {
        panic(err)
    }
}
// SetSensitivityLabelId sets the sensitivityLabelId property value. The sensitivityLabelId property
func (m *ItemItemsItemAssignSensitivityLabelPostRequestBody) SetSensitivityLabelId(value *string)() {
    err := m.GetBackingStore().Set("sensitivityLabelId", value)
    if err != nil {
        panic(err)
    }
}
type ItemItemsItemAssignSensitivityLabelPostRequestBodyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAssignmentMethod()(*iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SensitivityLabelAssignmentMethod)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetJustificationText()(*string)
    GetSensitivityLabelId()(*string)
    SetAssignmentMethod(value *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SensitivityLabelAssignmentMethod)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetJustificationText(value *string)()
    SetSensitivityLabelId(value *string)()
}
