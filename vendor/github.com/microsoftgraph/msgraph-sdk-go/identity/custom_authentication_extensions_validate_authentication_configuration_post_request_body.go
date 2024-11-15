package identity

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type CustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBody struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewCustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBody instantiates a new CustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBody and sets the default values.
func NewCustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBody()(*CustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBody) {
    m := &CustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBody{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateCustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBody(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *CustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBody) GetAdditionalData()(map[string]any) {
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
// GetAuthenticationConfiguration gets the authenticationConfiguration property value. The authenticationConfiguration property
// returns a CustomExtensionAuthenticationConfigurationable when successful
func (m *CustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBody) GetAuthenticationConfiguration()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionAuthenticationConfigurationable) {
    val, err := m.GetBackingStore().Get("authenticationConfiguration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionAuthenticationConfigurationable)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *CustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBody) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetEndpointConfiguration gets the endpointConfiguration property value. The endpointConfiguration property
// returns a CustomExtensionEndpointConfigurationable when successful
func (m *CustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBody) GetEndpointConfiguration()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionEndpointConfigurationable) {
    val, err := m.GetBackingStore().Get("endpointConfiguration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionEndpointConfigurationable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["authenticationConfiguration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateCustomExtensionAuthenticationConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAuthenticationConfiguration(val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionAuthenticationConfigurationable))
        }
        return nil
    }
    res["endpointConfiguration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateCustomExtensionEndpointConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEndpointConfiguration(val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionEndpointConfigurationable))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *CustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("authenticationConfiguration", m.GetAuthenticationConfiguration())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("endpointConfiguration", m.GetEndpointConfiguration())
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
func (m *CustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBody) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAuthenticationConfiguration sets the authenticationConfiguration property value. The authenticationConfiguration property
func (m *CustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBody) SetAuthenticationConfiguration(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionAuthenticationConfigurationable)() {
    err := m.GetBackingStore().Set("authenticationConfiguration", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *CustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBody) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetEndpointConfiguration sets the endpointConfiguration property value. The endpointConfiguration property
func (m *CustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBody) SetEndpointConfiguration(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionEndpointConfigurationable)() {
    err := m.GetBackingStore().Set("endpointConfiguration", value)
    if err != nil {
        panic(err)
    }
}
type CustomAuthenticationExtensionsValidateAuthenticationConfigurationPostRequestBodyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAuthenticationConfiguration()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionAuthenticationConfigurationable)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetEndpointConfiguration()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionEndpointConfigurationable)
    SetAuthenticationConfiguration(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionAuthenticationConfigurationable)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetEndpointConfiguration(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionEndpointConfigurationable)()
}
