package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SamlOrWsFedExternalDomainFederation struct {
    SamlOrWsFedProvider
}
// NewSamlOrWsFedExternalDomainFederation instantiates a new SamlOrWsFedExternalDomainFederation and sets the default values.
func NewSamlOrWsFedExternalDomainFederation()(*SamlOrWsFedExternalDomainFederation) {
    m := &SamlOrWsFedExternalDomainFederation{
        SamlOrWsFedProvider: *NewSamlOrWsFedProvider(),
    }
    odataTypeValue := "#microsoft.graph.samlOrWsFedExternalDomainFederation"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateSamlOrWsFedExternalDomainFederationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSamlOrWsFedExternalDomainFederationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSamlOrWsFedExternalDomainFederation(), nil
}
// GetDomains gets the domains property value. Collection of domain names of the external organizations that the tenant is federating with. Supports $filter (eq).
// returns a []ExternalDomainNameable when successful
func (m *SamlOrWsFedExternalDomainFederation) GetDomains()([]ExternalDomainNameable) {
    val, err := m.GetBackingStore().Get("domains")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ExternalDomainNameable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SamlOrWsFedExternalDomainFederation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.SamlOrWsFedProvider.GetFieldDeserializers()
    res["domains"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateExternalDomainNameFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ExternalDomainNameable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ExternalDomainNameable)
                }
            }
            m.SetDomains(res)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *SamlOrWsFedExternalDomainFederation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.SamlOrWsFedProvider.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetDomains() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDomains()))
        for i, v := range m.GetDomains() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("domains", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDomains sets the domains property value. Collection of domain names of the external organizations that the tenant is federating with. Supports $filter (eq).
func (m *SamlOrWsFedExternalDomainFederation) SetDomains(value []ExternalDomainNameable)() {
    err := m.GetBackingStore().Set("domains", value)
    if err != nil {
        panic(err)
    }
}
type SamlOrWsFedExternalDomainFederationable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    SamlOrWsFedProviderable
    GetDomains()([]ExternalDomainNameable)
    SetDomains(value []ExternalDomainNameable)()
}
