package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TermsOfUseContainer struct {
    Entity
}
// NewTermsOfUseContainer instantiates a new TermsOfUseContainer and sets the default values.
func NewTermsOfUseContainer()(*TermsOfUseContainer) {
    m := &TermsOfUseContainer{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTermsOfUseContainerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTermsOfUseContainerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTermsOfUseContainer(), nil
}
// GetAgreementAcceptances gets the agreementAcceptances property value. Represents the current status of a user's response to a company's customizable terms of use agreement.
// returns a []AgreementAcceptanceable when successful
func (m *TermsOfUseContainer) GetAgreementAcceptances()([]AgreementAcceptanceable) {
    val, err := m.GetBackingStore().Get("agreementAcceptances")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AgreementAcceptanceable)
    }
    return nil
}
// GetAgreements gets the agreements property value. Represents a tenant's customizable terms of use agreement that's created and managed with Microsoft Entra ID Governance.
// returns a []Agreementable when successful
func (m *TermsOfUseContainer) GetAgreements()([]Agreementable) {
    val, err := m.GetBackingStore().Get("agreements")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Agreementable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TermsOfUseContainer) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["agreementAcceptances"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAgreementAcceptanceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AgreementAcceptanceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AgreementAcceptanceable)
                }
            }
            m.SetAgreementAcceptances(res)
        }
        return nil
    }
    res["agreements"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAgreementFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Agreementable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Agreementable)
                }
            }
            m.SetAgreements(res)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *TermsOfUseContainer) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAgreementAcceptances() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAgreementAcceptances()))
        for i, v := range m.GetAgreementAcceptances() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("agreementAcceptances", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAgreements() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAgreements()))
        for i, v := range m.GetAgreements() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("agreements", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAgreementAcceptances sets the agreementAcceptances property value. Represents the current status of a user's response to a company's customizable terms of use agreement.
func (m *TermsOfUseContainer) SetAgreementAcceptances(value []AgreementAcceptanceable)() {
    err := m.GetBackingStore().Set("agreementAcceptances", value)
    if err != nil {
        panic(err)
    }
}
// SetAgreements sets the agreements property value. Represents a tenant's customizable terms of use agreement that's created and managed with Microsoft Entra ID Governance.
func (m *TermsOfUseContainer) SetAgreements(value []Agreementable)() {
    err := m.GetBackingStore().Set("agreements", value)
    if err != nil {
        panic(err)
    }
}
type TermsOfUseContainerable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAgreementAcceptances()([]AgreementAcceptanceable)
    GetAgreements()([]Agreementable)
    SetAgreementAcceptances(value []AgreementAcceptanceable)()
    SetAgreements(value []Agreementable)()
}
