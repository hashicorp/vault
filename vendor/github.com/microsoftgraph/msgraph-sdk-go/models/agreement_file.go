package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AgreementFile struct {
    AgreementFileProperties
}
// NewAgreementFile instantiates a new AgreementFile and sets the default values.
func NewAgreementFile()(*AgreementFile) {
    m := &AgreementFile{
        AgreementFileProperties: *NewAgreementFileProperties(),
    }
    return m
}
// CreateAgreementFileFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAgreementFileFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAgreementFile(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AgreementFile) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AgreementFileProperties.GetFieldDeserializers()
    res["localizations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAgreementFileLocalizationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AgreementFileLocalizationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AgreementFileLocalizationable)
                }
            }
            m.SetLocalizations(res)
        }
        return nil
    }
    return res
}
// GetLocalizations gets the localizations property value. The localized version of the terms of use agreement files attached to the agreement.
// returns a []AgreementFileLocalizationable when successful
func (m *AgreementFile) GetLocalizations()([]AgreementFileLocalizationable) {
    val, err := m.GetBackingStore().Get("localizations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AgreementFileLocalizationable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AgreementFile) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AgreementFileProperties.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetLocalizations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLocalizations()))
        for i, v := range m.GetLocalizations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("localizations", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetLocalizations sets the localizations property value. The localized version of the terms of use agreement files attached to the agreement.
func (m *AgreementFile) SetLocalizations(value []AgreementFileLocalizationable)() {
    err := m.GetBackingStore().Set("localizations", value)
    if err != nil {
        panic(err)
    }
}
type AgreementFileable interface {
    AgreementFilePropertiesable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetLocalizations()([]AgreementFileLocalizationable)
    SetLocalizations(value []AgreementFileLocalizationable)()
}
