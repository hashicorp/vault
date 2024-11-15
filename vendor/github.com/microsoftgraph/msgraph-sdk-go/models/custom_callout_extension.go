package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CustomCalloutExtension struct {
    Entity
}
// NewCustomCalloutExtension instantiates a new CustomCalloutExtension and sets the default values.
func NewCustomCalloutExtension()(*CustomCalloutExtension) {
    m := &CustomCalloutExtension{
        Entity: *NewEntity(),
    }
    return m
}
// CreateCustomCalloutExtensionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCustomCalloutExtensionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.accessPackageAssignmentRequestWorkflowExtension":
                        return NewAccessPackageAssignmentRequestWorkflowExtension(), nil
                    case "#microsoft.graph.accessPackageAssignmentWorkflowExtension":
                        return NewAccessPackageAssignmentWorkflowExtension(), nil
                    case "#microsoft.graph.customAuthenticationExtension":
                        return NewCustomAuthenticationExtension(), nil
                    case "#microsoft.graph.onTokenIssuanceStartCustomExtension":
                        return NewOnTokenIssuanceStartCustomExtension(), nil
                }
            }
        }
    }
    return NewCustomCalloutExtension(), nil
}
// GetAuthenticationConfiguration gets the authenticationConfiguration property value. Configuration for securing the API call to the logic app. For example, using OAuth client credentials flow.
// returns a CustomExtensionAuthenticationConfigurationable when successful
func (m *CustomCalloutExtension) GetAuthenticationConfiguration()(CustomExtensionAuthenticationConfigurationable) {
    val, err := m.GetBackingStore().Get("authenticationConfiguration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CustomExtensionAuthenticationConfigurationable)
    }
    return nil
}
// GetClientConfiguration gets the clientConfiguration property value. HTTP connection settings that define how long Microsoft Entra ID can wait for a connection to a logic app, how many times you can retry a timed-out connection and the exception scenarios when retries are allowed.
// returns a CustomExtensionClientConfigurationable when successful
func (m *CustomCalloutExtension) GetClientConfiguration()(CustomExtensionClientConfigurationable) {
    val, err := m.GetBackingStore().Get("clientConfiguration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CustomExtensionClientConfigurationable)
    }
    return nil
}
// GetDescription gets the description property value. Description for the customCalloutExtension object.
// returns a *string when successful
func (m *CustomCalloutExtension) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Display name for the customCalloutExtension object.
// returns a *string when successful
func (m *CustomCalloutExtension) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEndpointConfiguration gets the endpointConfiguration property value. The type and details for configuring the endpoint to call the logic app's workflow.
// returns a CustomExtensionEndpointConfigurationable when successful
func (m *CustomCalloutExtension) GetEndpointConfiguration()(CustomExtensionEndpointConfigurationable) {
    val, err := m.GetBackingStore().Get("endpointConfiguration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CustomExtensionEndpointConfigurationable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CustomCalloutExtension) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["authenticationConfiguration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCustomExtensionAuthenticationConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAuthenticationConfiguration(val.(CustomExtensionAuthenticationConfigurationable))
        }
        return nil
    }
    res["clientConfiguration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCustomExtensionClientConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClientConfiguration(val.(CustomExtensionClientConfigurationable))
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
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["endpointConfiguration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCustomExtensionEndpointConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEndpointConfiguration(val.(CustomExtensionEndpointConfigurationable))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *CustomCalloutExtension) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("authenticationConfiguration", m.GetAuthenticationConfiguration())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("clientConfiguration", m.GetClientConfiguration())
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
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("endpointConfiguration", m.GetEndpointConfiguration())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAuthenticationConfiguration sets the authenticationConfiguration property value. Configuration for securing the API call to the logic app. For example, using OAuth client credentials flow.
func (m *CustomCalloutExtension) SetAuthenticationConfiguration(value CustomExtensionAuthenticationConfigurationable)() {
    err := m.GetBackingStore().Set("authenticationConfiguration", value)
    if err != nil {
        panic(err)
    }
}
// SetClientConfiguration sets the clientConfiguration property value. HTTP connection settings that define how long Microsoft Entra ID can wait for a connection to a logic app, how many times you can retry a timed-out connection and the exception scenarios when retries are allowed.
func (m *CustomCalloutExtension) SetClientConfiguration(value CustomExtensionClientConfigurationable)() {
    err := m.GetBackingStore().Set("clientConfiguration", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. Description for the customCalloutExtension object.
func (m *CustomCalloutExtension) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Display name for the customCalloutExtension object.
func (m *CustomCalloutExtension) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetEndpointConfiguration sets the endpointConfiguration property value. The type and details for configuring the endpoint to call the logic app's workflow.
func (m *CustomCalloutExtension) SetEndpointConfiguration(value CustomExtensionEndpointConfigurationable)() {
    err := m.GetBackingStore().Set("endpointConfiguration", value)
    if err != nil {
        panic(err)
    }
}
type CustomCalloutExtensionable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAuthenticationConfiguration()(CustomExtensionAuthenticationConfigurationable)
    GetClientConfiguration()(CustomExtensionClientConfigurationable)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetEndpointConfiguration()(CustomExtensionEndpointConfigurationable)
    SetAuthenticationConfiguration(value CustomExtensionAuthenticationConfigurationable)()
    SetClientConfiguration(value CustomExtensionClientConfigurationable)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetEndpointConfiguration(value CustomExtensionEndpointConfigurationable)()
}
