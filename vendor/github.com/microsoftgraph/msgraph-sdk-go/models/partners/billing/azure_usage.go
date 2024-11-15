package billing

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type AzureUsage struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
}
// NewAzureUsage instantiates a new AzureUsage and sets the default values.
func NewAzureUsage()(*AzureUsage) {
    m := &AzureUsage{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateAzureUsageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAzureUsageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAzureUsage(), nil
}
// GetBilled gets the billed property value. The billed property
// returns a BilledUsageable when successful
func (m *AzureUsage) GetBilled()(BilledUsageable) {
    val, err := m.GetBackingStore().Get("billed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(BilledUsageable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AzureUsage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["billed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateBilledUsageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBilled(val.(BilledUsageable))
        }
        return nil
    }
    res["unbilled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUnbilledUsageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUnbilled(val.(UnbilledUsageable))
        }
        return nil
    }
    return res
}
// GetUnbilled gets the unbilled property value. The unbilled property
// returns a UnbilledUsageable when successful
func (m *AzureUsage) GetUnbilled()(UnbilledUsageable) {
    val, err := m.GetBackingStore().Get("unbilled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UnbilledUsageable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AzureUsage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("billed", m.GetBilled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("unbilled", m.GetUnbilled())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBilled sets the billed property value. The billed property
func (m *AzureUsage) SetBilled(value BilledUsageable)() {
    err := m.GetBackingStore().Set("billed", value)
    if err != nil {
        panic(err)
    }
}
// SetUnbilled sets the unbilled property value. The unbilled property
func (m *AzureUsage) SetUnbilled(value UnbilledUsageable)() {
    err := m.GetBackingStore().Set("unbilled", value)
    if err != nil {
        panic(err)
    }
}
type AzureUsageable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBilled()(BilledUsageable)
    GetUnbilled()(UnbilledUsageable)
    SetBilled(value BilledUsageable)()
    SetUnbilled(value UnbilledUsageable)()
}
