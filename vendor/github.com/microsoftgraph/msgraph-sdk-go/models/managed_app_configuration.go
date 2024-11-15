package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ManagedAppConfiguration configuration used to deliver a set of custom settings as-is to apps for users to whom the configuration is scoped
type ManagedAppConfiguration struct {
    ManagedAppPolicy
}
// NewManagedAppConfiguration instantiates a new ManagedAppConfiguration and sets the default values.
func NewManagedAppConfiguration()(*ManagedAppConfiguration) {
    m := &ManagedAppConfiguration{
        ManagedAppPolicy: *NewManagedAppPolicy(),
    }
    odataTypeValue := "#microsoft.graph.managedAppConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateManagedAppConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateManagedAppConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.targetedManagedAppConfiguration":
                        return NewTargetedManagedAppConfiguration(), nil
                }
            }
        }
    }
    return NewManagedAppConfiguration(), nil
}
// GetCustomSettings gets the customSettings property value. A set of string key and string value pairs to be sent to apps for users to whom the configuration is scoped, unalterned by this service
// returns a []KeyValuePairable when successful
func (m *ManagedAppConfiguration) GetCustomSettings()([]KeyValuePairable) {
    val, err := m.GetBackingStore().Get("customSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]KeyValuePairable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ManagedAppConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ManagedAppPolicy.GetFieldDeserializers()
    res["customSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateKeyValuePairFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]KeyValuePairable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(KeyValuePairable)
                }
            }
            m.SetCustomSettings(res)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *ManagedAppConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ManagedAppPolicy.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCustomSettings() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCustomSettings()))
        for i, v := range m.GetCustomSettings() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("customSettings", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCustomSettings sets the customSettings property value. A set of string key and string value pairs to be sent to apps for users to whom the configuration is scoped, unalterned by this service
func (m *ManagedAppConfiguration) SetCustomSettings(value []KeyValuePairable)() {
    err := m.GetBackingStore().Set("customSettings", value)
    if err != nil {
        panic(err)
    }
}
type ManagedAppConfigurationable interface {
    ManagedAppPolicyable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCustomSettings()([]KeyValuePairable)
    SetCustomSettings(value []KeyValuePairable)()
}
