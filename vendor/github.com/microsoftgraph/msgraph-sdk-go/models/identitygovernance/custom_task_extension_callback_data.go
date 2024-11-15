package identitygovernance

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type CustomTaskExtensionCallbackData struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionData
}
// NewCustomTaskExtensionCallbackData instantiates a new CustomTaskExtensionCallbackData and sets the default values.
func NewCustomTaskExtensionCallbackData()(*CustomTaskExtensionCallbackData) {
    m := &CustomTaskExtensionCallbackData{
        CustomExtensionData: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewCustomExtensionData(),
    }
    odataTypeValue := "#microsoft.graph.identityGovernance.customTaskExtensionCallbackData"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateCustomTaskExtensionCallbackDataFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCustomTaskExtensionCallbackDataFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCustomTaskExtensionCallbackData(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CustomTaskExtensionCallbackData) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.CustomExtensionData.GetFieldDeserializers()
    res["operationStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseCustomTaskExtensionOperationStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOperationStatus(val.(*CustomTaskExtensionOperationStatus))
        }
        return nil
    }
    return res
}
// GetOperationStatus gets the operationStatus property value. Operation status that's provided by the Azure Logic App indicating whenever the Azure Logic App has run successfully or not. Supported values: completed, failed, unknownFutureValue.
// returns a *CustomTaskExtensionOperationStatus when successful
func (m *CustomTaskExtensionCallbackData) GetOperationStatus()(*CustomTaskExtensionOperationStatus) {
    val, err := m.GetBackingStore().Get("operationStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*CustomTaskExtensionOperationStatus)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CustomTaskExtensionCallbackData) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.CustomExtensionData.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetOperationStatus() != nil {
        cast := (*m.GetOperationStatus()).String()
        err = writer.WriteStringValue("operationStatus", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetOperationStatus sets the operationStatus property value. Operation status that's provided by the Azure Logic App indicating whenever the Azure Logic App has run successfully or not. Supported values: completed, failed, unknownFutureValue.
func (m *CustomTaskExtensionCallbackData) SetOperationStatus(value *CustomTaskExtensionOperationStatus)() {
    err := m.GetBackingStore().Set("operationStatus", value)
    if err != nil {
        panic(err)
    }
}
type CustomTaskExtensionCallbackDataable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionDataable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetOperationStatus()(*CustomTaskExtensionOperationStatus)
    SetOperationStatus(value *CustomTaskExtensionOperationStatus)()
}
