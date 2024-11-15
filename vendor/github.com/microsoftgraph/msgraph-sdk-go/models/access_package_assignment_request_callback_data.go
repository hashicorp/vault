package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessPackageAssignmentRequestCallbackData struct {
    CustomExtensionData
}
// NewAccessPackageAssignmentRequestCallbackData instantiates a new AccessPackageAssignmentRequestCallbackData and sets the default values.
func NewAccessPackageAssignmentRequestCallbackData()(*AccessPackageAssignmentRequestCallbackData) {
    m := &AccessPackageAssignmentRequestCallbackData{
        CustomExtensionData: *NewCustomExtensionData(),
    }
    odataTypeValue := "#microsoft.graph.accessPackageAssignmentRequestCallbackData"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAccessPackageAssignmentRequestCallbackDataFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessPackageAssignmentRequestCallbackDataFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessPackageAssignmentRequestCallbackData(), nil
}
// GetCustomExtensionStageInstanceDetail gets the customExtensionStageInstanceDetail property value. Details for the callback.
// returns a *string when successful
func (m *AccessPackageAssignmentRequestCallbackData) GetCustomExtensionStageInstanceDetail()(*string) {
    val, err := m.GetBackingStore().Get("customExtensionStageInstanceDetail")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCustomExtensionStageInstanceId gets the customExtensionStageInstanceId property value. Unique identifier of the callout to the custom extension.
// returns a *string when successful
func (m *AccessPackageAssignmentRequestCallbackData) GetCustomExtensionStageInstanceId()(*string) {
    val, err := m.GetBackingStore().Get("customExtensionStageInstanceId")
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
func (m *AccessPackageAssignmentRequestCallbackData) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.CustomExtensionData.GetFieldDeserializers()
    res["customExtensionStageInstanceDetail"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCustomExtensionStageInstanceDetail(val)
        }
        return nil
    }
    res["customExtensionStageInstanceId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCustomExtensionStageInstanceId(val)
        }
        return nil
    }
    res["stage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAccessPackageCustomExtensionStage)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStage(val.(*AccessPackageCustomExtensionStage))
        }
        return nil
    }
    res["state"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetState(val)
        }
        return nil
    }
    return res
}
// GetStage gets the stage property value. Indicates the stage at which the custom callout extension is executed. The possible values are: assignmentRequestCreated, assignmentRequestApproved, assignmentRequestGranted, assignmentRequestRemoved, assignmentFourteenDaysBeforeExpiration, assignmentOneDayBeforeExpiration, unknownFutureValue.
// returns a *AccessPackageCustomExtensionStage when successful
func (m *AccessPackageAssignmentRequestCallbackData) GetStage()(*AccessPackageCustomExtensionStage) {
    val, err := m.GetBackingStore().Get("stage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AccessPackageCustomExtensionStage)
    }
    return nil
}
// GetState gets the state property value. Allow the extension to be able to deny or cancel the request submitted by the requestor. The supported values are Denied and Canceled. This property can only be set for an assignmentRequestCreated stage.
// returns a *string when successful
func (m *AccessPackageAssignmentRequestCallbackData) GetState()(*string) {
    val, err := m.GetBackingStore().Get("state")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessPackageAssignmentRequestCallbackData) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.CustomExtensionData.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("customExtensionStageInstanceDetail", m.GetCustomExtensionStageInstanceDetail())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("customExtensionStageInstanceId", m.GetCustomExtensionStageInstanceId())
        if err != nil {
            return err
        }
    }
    if m.GetStage() != nil {
        cast := (*m.GetStage()).String()
        err = writer.WriteStringValue("stage", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("state", m.GetState())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCustomExtensionStageInstanceDetail sets the customExtensionStageInstanceDetail property value. Details for the callback.
func (m *AccessPackageAssignmentRequestCallbackData) SetCustomExtensionStageInstanceDetail(value *string)() {
    err := m.GetBackingStore().Set("customExtensionStageInstanceDetail", value)
    if err != nil {
        panic(err)
    }
}
// SetCustomExtensionStageInstanceId sets the customExtensionStageInstanceId property value. Unique identifier of the callout to the custom extension.
func (m *AccessPackageAssignmentRequestCallbackData) SetCustomExtensionStageInstanceId(value *string)() {
    err := m.GetBackingStore().Set("customExtensionStageInstanceId", value)
    if err != nil {
        panic(err)
    }
}
// SetStage sets the stage property value. Indicates the stage at which the custom callout extension is executed. The possible values are: assignmentRequestCreated, assignmentRequestApproved, assignmentRequestGranted, assignmentRequestRemoved, assignmentFourteenDaysBeforeExpiration, assignmentOneDayBeforeExpiration, unknownFutureValue.
func (m *AccessPackageAssignmentRequestCallbackData) SetStage(value *AccessPackageCustomExtensionStage)() {
    err := m.GetBackingStore().Set("stage", value)
    if err != nil {
        panic(err)
    }
}
// SetState sets the state property value. Allow the extension to be able to deny or cancel the request submitted by the requestor. The supported values are Denied and Canceled. This property can only be set for an assignmentRequestCreated stage.
func (m *AccessPackageAssignmentRequestCallbackData) SetState(value *string)() {
    err := m.GetBackingStore().Set("state", value)
    if err != nil {
        panic(err)
    }
}
type AccessPackageAssignmentRequestCallbackDataable interface {
    CustomExtensionDataable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCustomExtensionStageInstanceDetail()(*string)
    GetCustomExtensionStageInstanceId()(*string)
    GetStage()(*AccessPackageCustomExtensionStage)
    GetState()(*string)
    SetCustomExtensionStageInstanceDetail(value *string)()
    SetCustomExtensionStageInstanceId(value *string)()
    SetStage(value *AccessPackageCustomExtensionStage)()
    SetState(value *string)()
}
