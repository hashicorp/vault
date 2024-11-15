package identitygovernance

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type CustomTaskExtensionCallbackConfiguration struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionCallbackConfiguration
}
// NewCustomTaskExtensionCallbackConfiguration instantiates a new CustomTaskExtensionCallbackConfiguration and sets the default values.
func NewCustomTaskExtensionCallbackConfiguration()(*CustomTaskExtensionCallbackConfiguration) {
    m := &CustomTaskExtensionCallbackConfiguration{
        CustomExtensionCallbackConfiguration: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewCustomExtensionCallbackConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.identityGovernance.customTaskExtensionCallbackConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateCustomTaskExtensionCallbackConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCustomTaskExtensionCallbackConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCustomTaskExtensionCallbackConfiguration(), nil
}
// GetAuthorizedApps gets the authorizedApps property value. The authorizedApps property
// returns a []Applicationable when successful
func (m *CustomTaskExtensionCallbackConfiguration) GetAuthorizedApps()([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Applicationable) {
    val, err := m.GetBackingStore().Get("authorizedApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Applicationable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CustomTaskExtensionCallbackConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.CustomExtensionCallbackConfiguration.GetFieldDeserializers()
    res["authorizedApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateApplicationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Applicationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Applicationable)
                }
            }
            m.SetAuthorizedApps(res)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *CustomTaskExtensionCallbackConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.CustomExtensionCallbackConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAuthorizedApps() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAuthorizedApps()))
        for i, v := range m.GetAuthorizedApps() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("authorizedApps", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAuthorizedApps sets the authorizedApps property value. The authorizedApps property
func (m *CustomTaskExtensionCallbackConfiguration) SetAuthorizedApps(value []iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Applicationable)() {
    err := m.GetBackingStore().Set("authorizedApps", value)
    if err != nil {
        panic(err)
    }
}
type CustomTaskExtensionCallbackConfigurationable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionCallbackConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAuthorizedApps()([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Applicationable)
    SetAuthorizedApps(value []iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Applicationable)()
}
