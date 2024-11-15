package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CloudAppSecuritySessionControl struct {
    ConditionalAccessSessionControl
}
// NewCloudAppSecuritySessionControl instantiates a new CloudAppSecuritySessionControl and sets the default values.
func NewCloudAppSecuritySessionControl()(*CloudAppSecuritySessionControl) {
    m := &CloudAppSecuritySessionControl{
        ConditionalAccessSessionControl: *NewConditionalAccessSessionControl(),
    }
    odataTypeValue := "#microsoft.graph.cloudAppSecuritySessionControl"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateCloudAppSecuritySessionControlFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCloudAppSecuritySessionControlFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCloudAppSecuritySessionControl(), nil
}
// GetCloudAppSecurityType gets the cloudAppSecurityType property value. Possible values are: mcasConfigured, monitorOnly, blockDownloads, unknownFutureValue. For more information, see Deploy Conditional Access App Control for featured apps.
// returns a *CloudAppSecuritySessionControlType when successful
func (m *CloudAppSecuritySessionControl) GetCloudAppSecurityType()(*CloudAppSecuritySessionControlType) {
    val, err := m.GetBackingStore().Get("cloudAppSecurityType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*CloudAppSecuritySessionControlType)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CloudAppSecuritySessionControl) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ConditionalAccessSessionControl.GetFieldDeserializers()
    res["cloudAppSecurityType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseCloudAppSecuritySessionControlType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCloudAppSecurityType(val.(*CloudAppSecuritySessionControlType))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *CloudAppSecuritySessionControl) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ConditionalAccessSessionControl.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCloudAppSecurityType() != nil {
        cast := (*m.GetCloudAppSecurityType()).String()
        err = writer.WriteStringValue("cloudAppSecurityType", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCloudAppSecurityType sets the cloudAppSecurityType property value. Possible values are: mcasConfigured, monitorOnly, blockDownloads, unknownFutureValue. For more information, see Deploy Conditional Access App Control for featured apps.
func (m *CloudAppSecuritySessionControl) SetCloudAppSecurityType(value *CloudAppSecuritySessionControlType)() {
    err := m.GetBackingStore().Set("cloudAppSecurityType", value)
    if err != nil {
        panic(err)
    }
}
type CloudAppSecuritySessionControlable interface {
    ConditionalAccessSessionControlable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCloudAppSecurityType()(*CloudAppSecuritySessionControlType)
    SetCloudAppSecurityType(value *CloudAppSecuritySessionControlType)()
}
