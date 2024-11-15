package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type OnTokenIssuanceStartCustomExtension struct {
    CustomAuthenticationExtension
}
// NewOnTokenIssuanceStartCustomExtension instantiates a new OnTokenIssuanceStartCustomExtension and sets the default values.
func NewOnTokenIssuanceStartCustomExtension()(*OnTokenIssuanceStartCustomExtension) {
    m := &OnTokenIssuanceStartCustomExtension{
        CustomAuthenticationExtension: *NewCustomAuthenticationExtension(),
    }
    odataTypeValue := "#microsoft.graph.onTokenIssuanceStartCustomExtension"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateOnTokenIssuanceStartCustomExtensionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOnTokenIssuanceStartCustomExtensionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOnTokenIssuanceStartCustomExtension(), nil
}
// GetClaimsForTokenConfiguration gets the claimsForTokenConfiguration property value. Collection of claims to be returned by the API called by this custom authentication extension. Used to populate claims mapping experience in Microsoft Entra admin center. Optional.
// returns a []OnTokenIssuanceStartReturnClaimable when successful
func (m *OnTokenIssuanceStartCustomExtension) GetClaimsForTokenConfiguration()([]OnTokenIssuanceStartReturnClaimable) {
    val, err := m.GetBackingStore().Get("claimsForTokenConfiguration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]OnTokenIssuanceStartReturnClaimable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OnTokenIssuanceStartCustomExtension) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.CustomAuthenticationExtension.GetFieldDeserializers()
    res["claimsForTokenConfiguration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateOnTokenIssuanceStartReturnClaimFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]OnTokenIssuanceStartReturnClaimable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(OnTokenIssuanceStartReturnClaimable)
                }
            }
            m.SetClaimsForTokenConfiguration(res)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *OnTokenIssuanceStartCustomExtension) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.CustomAuthenticationExtension.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetClaimsForTokenConfiguration() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetClaimsForTokenConfiguration()))
        for i, v := range m.GetClaimsForTokenConfiguration() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("claimsForTokenConfiguration", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetClaimsForTokenConfiguration sets the claimsForTokenConfiguration property value. Collection of claims to be returned by the API called by this custom authentication extension. Used to populate claims mapping experience in Microsoft Entra admin center. Optional.
func (m *OnTokenIssuanceStartCustomExtension) SetClaimsForTokenConfiguration(value []OnTokenIssuanceStartReturnClaimable)() {
    err := m.GetBackingStore().Set("claimsForTokenConfiguration", value)
    if err != nil {
        panic(err)
    }
}
type OnTokenIssuanceStartCustomExtensionable interface {
    CustomAuthenticationExtensionable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetClaimsForTokenConfiguration()([]OnTokenIssuanceStartReturnClaimable)
    SetClaimsForTokenConfiguration(value []OnTokenIssuanceStartReturnClaimable)()
}
