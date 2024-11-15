package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TermsAndConditionsAcceptanceStatusCollectionResponse struct {
    BaseCollectionPaginationCountResponse
}
// NewTermsAndConditionsAcceptanceStatusCollectionResponse instantiates a new TermsAndConditionsAcceptanceStatusCollectionResponse and sets the default values.
func NewTermsAndConditionsAcceptanceStatusCollectionResponse()(*TermsAndConditionsAcceptanceStatusCollectionResponse) {
    m := &TermsAndConditionsAcceptanceStatusCollectionResponse{
        BaseCollectionPaginationCountResponse: *NewBaseCollectionPaginationCountResponse(),
    }
    return m
}
// CreateTermsAndConditionsAcceptanceStatusCollectionResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTermsAndConditionsAcceptanceStatusCollectionResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTermsAndConditionsAcceptanceStatusCollectionResponse(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TermsAndConditionsAcceptanceStatusCollectionResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseCollectionPaginationCountResponse.GetFieldDeserializers()
    res["value"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTermsAndConditionsAcceptanceStatusFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TermsAndConditionsAcceptanceStatusable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TermsAndConditionsAcceptanceStatusable)
                }
            }
            m.SetValue(res)
        }
        return nil
    }
    return res
}
// GetValue gets the value property value. The value property
// returns a []TermsAndConditionsAcceptanceStatusable when successful
func (m *TermsAndConditionsAcceptanceStatusCollectionResponse) GetValue()([]TermsAndConditionsAcceptanceStatusable) {
    val, err := m.GetBackingStore().Get("value")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TermsAndConditionsAcceptanceStatusable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TermsAndConditionsAcceptanceStatusCollectionResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseCollectionPaginationCountResponse.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetValue() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetValue()))
        for i, v := range m.GetValue() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("value", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetValue sets the value property value. The value property
func (m *TermsAndConditionsAcceptanceStatusCollectionResponse) SetValue(value []TermsAndConditionsAcceptanceStatusable)() {
    err := m.GetBackingStore().Set("value", value)
    if err != nil {
        panic(err)
    }
}
type TermsAndConditionsAcceptanceStatusCollectionResponseable interface {
    BaseCollectionPaginationCountResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetValue()([]TermsAndConditionsAcceptanceStatusable)
    SetValue(value []TermsAndConditionsAcceptanceStatusable)()
}
