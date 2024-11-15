package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AgreementFileLocalization struct {
    AgreementFileProperties
}
// NewAgreementFileLocalization instantiates a new AgreementFileLocalization and sets the default values.
func NewAgreementFileLocalization()(*AgreementFileLocalization) {
    m := &AgreementFileLocalization{
        AgreementFileProperties: *NewAgreementFileProperties(),
    }
    return m
}
// CreateAgreementFileLocalizationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAgreementFileLocalizationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAgreementFileLocalization(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AgreementFileLocalization) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AgreementFileProperties.GetFieldDeserializers()
    res["versions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAgreementFileVersionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AgreementFileVersionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AgreementFileVersionable)
                }
            }
            m.SetVersions(res)
        }
        return nil
    }
    return res
}
// GetVersions gets the versions property value. Read-only. Customized versions of the terms of use agreement in the Microsoft Entra tenant.
// returns a []AgreementFileVersionable when successful
func (m *AgreementFileLocalization) GetVersions()([]AgreementFileVersionable) {
    val, err := m.GetBackingStore().Get("versions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AgreementFileVersionable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AgreementFileLocalization) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AgreementFileProperties.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetVersions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetVersions()))
        for i, v := range m.GetVersions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("versions", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetVersions sets the versions property value. Read-only. Customized versions of the terms of use agreement in the Microsoft Entra tenant.
func (m *AgreementFileLocalization) SetVersions(value []AgreementFileVersionable)() {
    err := m.GetBackingStore().Set("versions", value)
    if err != nil {
        panic(err)
    }
}
type AgreementFileLocalizationable interface {
    AgreementFilePropertiesable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetVersions()([]AgreementFileVersionable)
    SetVersions(value []AgreementFileVersionable)()
}
