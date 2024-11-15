package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Agreement struct {
    Entity
}
// NewAgreement instantiates a new Agreement and sets the default values.
func NewAgreement()(*Agreement) {
    m := &Agreement{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAgreementFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAgreementFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAgreement(), nil
}
// GetAcceptances gets the acceptances property value. Read-only. Information about acceptances of this agreement.
// returns a []AgreementAcceptanceable when successful
func (m *Agreement) GetAcceptances()([]AgreementAcceptanceable) {
    val, err := m.GetBackingStore().Get("acceptances")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AgreementAcceptanceable)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Display name of the agreement. The display name is used for internal tracking of the agreement but isn't shown to end users who view the agreement. Supports $filter (eq).
// returns a *string when successful
func (m *Agreement) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Agreement) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["acceptances"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetAcceptances(res)
        }
        return nil
    }
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["file"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAgreementFileFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFile(val.(AgreementFileable))
        }
        return nil
    }
    res["files"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetFiles(res)
        }
        return nil
    }
    res["isPerDeviceAcceptanceRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsPerDeviceAcceptanceRequired(val)
        }
        return nil
    }
    res["isViewingBeforeAcceptanceRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsViewingBeforeAcceptanceRequired(val)
        }
        return nil
    }
    res["termsExpiration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTermsExpirationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTermsExpiration(val.(TermsExpirationable))
        }
        return nil
    }
    res["userReacceptRequiredFrequency"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserReacceptRequiredFrequency(val)
        }
        return nil
    }
    return res
}
// GetFile gets the file property value. Default PDF linked to this agreement.
// returns a AgreementFileable when successful
func (m *Agreement) GetFile()(AgreementFileable) {
    val, err := m.GetBackingStore().Get("file")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AgreementFileable)
    }
    return nil
}
// GetFiles gets the files property value. PDFs linked to this agreement. This property is in the process of being deprecated. Use the  file property instead. Supports $expand.
// returns a []AgreementFileLocalizationable when successful
func (m *Agreement) GetFiles()([]AgreementFileLocalizationable) {
    val, err := m.GetBackingStore().Get("files")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AgreementFileLocalizationable)
    }
    return nil
}
// GetIsPerDeviceAcceptanceRequired gets the isPerDeviceAcceptanceRequired property value. Indicates whether end users are required to accept this agreement on every device that they access it from. The end user is required to register their device in Microsoft Entra ID, if they haven't already done so. Supports $filter (eq).
// returns a *bool when successful
func (m *Agreement) GetIsPerDeviceAcceptanceRequired()(*bool) {
    val, err := m.GetBackingStore().Get("isPerDeviceAcceptanceRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsViewingBeforeAcceptanceRequired gets the isViewingBeforeAcceptanceRequired property value. Indicates whether the user has to expand the agreement before accepting. Supports $filter (eq).
// returns a *bool when successful
func (m *Agreement) GetIsViewingBeforeAcceptanceRequired()(*bool) {
    val, err := m.GetBackingStore().Get("isViewingBeforeAcceptanceRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetTermsExpiration gets the termsExpiration property value. Expiration schedule and frequency of agreement for all users. Supports $filter (eq).
// returns a TermsExpirationable when successful
func (m *Agreement) GetTermsExpiration()(TermsExpirationable) {
    val, err := m.GetBackingStore().Get("termsExpiration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TermsExpirationable)
    }
    return nil
}
// GetUserReacceptRequiredFrequency gets the userReacceptRequiredFrequency property value. The duration after which the user must reaccept the terms of use. The value is represented in ISO 8601 format for durations. Supports $filter (eq).
// returns a *ISODuration when successful
func (m *Agreement) GetUserReacceptRequiredFrequency()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("userReacceptRequiredFrequency")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Agreement) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAcceptances() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAcceptances()))
        for i, v := range m.GetAcceptances() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("acceptances", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("file", m.GetFile())
        if err != nil {
            return err
        }
    }
    if m.GetFiles() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetFiles()))
        for i, v := range m.GetFiles() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("files", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isPerDeviceAcceptanceRequired", m.GetIsPerDeviceAcceptanceRequired())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isViewingBeforeAcceptanceRequired", m.GetIsViewingBeforeAcceptanceRequired())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("termsExpiration", m.GetTermsExpiration())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteISODurationValue("userReacceptRequiredFrequency", m.GetUserReacceptRequiredFrequency())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAcceptances sets the acceptances property value. Read-only. Information about acceptances of this agreement.
func (m *Agreement) SetAcceptances(value []AgreementAcceptanceable)() {
    err := m.GetBackingStore().Set("acceptances", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Display name of the agreement. The display name is used for internal tracking of the agreement but isn't shown to end users who view the agreement. Supports $filter (eq).
func (m *Agreement) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetFile sets the file property value. Default PDF linked to this agreement.
func (m *Agreement) SetFile(value AgreementFileable)() {
    err := m.GetBackingStore().Set("file", value)
    if err != nil {
        panic(err)
    }
}
// SetFiles sets the files property value. PDFs linked to this agreement. This property is in the process of being deprecated. Use the  file property instead. Supports $expand.
func (m *Agreement) SetFiles(value []AgreementFileLocalizationable)() {
    err := m.GetBackingStore().Set("files", value)
    if err != nil {
        panic(err)
    }
}
// SetIsPerDeviceAcceptanceRequired sets the isPerDeviceAcceptanceRequired property value. Indicates whether end users are required to accept this agreement on every device that they access it from. The end user is required to register their device in Microsoft Entra ID, if they haven't already done so. Supports $filter (eq).
func (m *Agreement) SetIsPerDeviceAcceptanceRequired(value *bool)() {
    err := m.GetBackingStore().Set("isPerDeviceAcceptanceRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetIsViewingBeforeAcceptanceRequired sets the isViewingBeforeAcceptanceRequired property value. Indicates whether the user has to expand the agreement before accepting. Supports $filter (eq).
func (m *Agreement) SetIsViewingBeforeAcceptanceRequired(value *bool)() {
    err := m.GetBackingStore().Set("isViewingBeforeAcceptanceRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetTermsExpiration sets the termsExpiration property value. Expiration schedule and frequency of agreement for all users. Supports $filter (eq).
func (m *Agreement) SetTermsExpiration(value TermsExpirationable)() {
    err := m.GetBackingStore().Set("termsExpiration", value)
    if err != nil {
        panic(err)
    }
}
// SetUserReacceptRequiredFrequency sets the userReacceptRequiredFrequency property value. The duration after which the user must reaccept the terms of use. The value is represented in ISO 8601 format for durations. Supports $filter (eq).
func (m *Agreement) SetUserReacceptRequiredFrequency(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("userReacceptRequiredFrequency", value)
    if err != nil {
        panic(err)
    }
}
type Agreementable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAcceptances()([]AgreementAcceptanceable)
    GetDisplayName()(*string)
    GetFile()(AgreementFileable)
    GetFiles()([]AgreementFileLocalizationable)
    GetIsPerDeviceAcceptanceRequired()(*bool)
    GetIsViewingBeforeAcceptanceRequired()(*bool)
    GetTermsExpiration()(TermsExpirationable)
    GetUserReacceptRequiredFrequency()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    SetAcceptances(value []AgreementAcceptanceable)()
    SetDisplayName(value *string)()
    SetFile(value AgreementFileable)()
    SetFiles(value []AgreementFileLocalizationable)()
    SetIsPerDeviceAcceptanceRequired(value *bool)()
    SetIsViewingBeforeAcceptanceRequired(value *bool)()
    SetTermsExpiration(value TermsExpirationable)()
    SetUserReacceptRequiredFrequency(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
}
