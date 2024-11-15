package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type IdentityGovernance struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewIdentityGovernance instantiates a new IdentityGovernance and sets the default values.
func NewIdentityGovernance()(*IdentityGovernance) {
    m := &IdentityGovernance{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateIdentityGovernanceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIdentityGovernanceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIdentityGovernance(), nil
}
// GetAccessReviews gets the accessReviews property value. The accessReviews property
// returns a AccessReviewSetable when successful
func (m *IdentityGovernance) GetAccessReviews()(AccessReviewSetable) {
    val, err := m.GetBackingStore().Get("accessReviews")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessReviewSetable)
    }
    return nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *IdentityGovernance) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetAppConsent gets the appConsent property value. The appConsent property
// returns a AppConsentApprovalRouteable when successful
func (m *IdentityGovernance) GetAppConsent()(AppConsentApprovalRouteable) {
    val, err := m.GetBackingStore().Get("appConsent")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AppConsentApprovalRouteable)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *IdentityGovernance) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetEntitlementManagement gets the entitlementManagement property value. The entitlementManagement property
// returns a EntitlementManagementable when successful
func (m *IdentityGovernance) GetEntitlementManagement()(EntitlementManagementable) {
    val, err := m.GetBackingStore().Get("entitlementManagement")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EntitlementManagementable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *IdentityGovernance) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["accessReviews"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessReviewSetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAccessReviews(val.(AccessReviewSetable))
        }
        return nil
    }
    res["appConsent"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAppConsentApprovalRouteFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppConsent(val.(AppConsentApprovalRouteable))
        }
        return nil
    }
    res["entitlementManagement"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEntitlementManagementFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEntitlementManagement(val.(EntitlementManagementable))
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    res["privilegedAccess"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePrivilegedAccessRootFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrivilegedAccess(val.(PrivilegedAccessRootable))
        }
        return nil
    }
    res["termsOfUse"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTermsOfUseContainerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTermsOfUse(val.(TermsOfUseContainerable))
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *IdentityGovernance) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPrivilegedAccess gets the privilegedAccess property value. The privilegedAccess property
// returns a PrivilegedAccessRootable when successful
func (m *IdentityGovernance) GetPrivilegedAccess()(PrivilegedAccessRootable) {
    val, err := m.GetBackingStore().Get("privilegedAccess")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PrivilegedAccessRootable)
    }
    return nil
}
// GetTermsOfUse gets the termsOfUse property value. The termsOfUse property
// returns a TermsOfUseContainerable when successful
func (m *IdentityGovernance) GetTermsOfUse()(TermsOfUseContainerable) {
    val, err := m.GetBackingStore().Get("termsOfUse")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TermsOfUseContainerable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IdentityGovernance) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("accessReviews", m.GetAccessReviews())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("appConsent", m.GetAppConsent())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("entitlementManagement", m.GetEntitlementManagement())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("privilegedAccess", m.GetPrivilegedAccess())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("termsOfUse", m.GetTermsOfUse())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccessReviews sets the accessReviews property value. The accessReviews property
func (m *IdentityGovernance) SetAccessReviews(value AccessReviewSetable)() {
    err := m.GetBackingStore().Set("accessReviews", value)
    if err != nil {
        panic(err)
    }
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *IdentityGovernance) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAppConsent sets the appConsent property value. The appConsent property
func (m *IdentityGovernance) SetAppConsent(value AppConsentApprovalRouteable)() {
    err := m.GetBackingStore().Set("appConsent", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *IdentityGovernance) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetEntitlementManagement sets the entitlementManagement property value. The entitlementManagement property
func (m *IdentityGovernance) SetEntitlementManagement(value EntitlementManagementable)() {
    err := m.GetBackingStore().Set("entitlementManagement", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *IdentityGovernance) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPrivilegedAccess sets the privilegedAccess property value. The privilegedAccess property
func (m *IdentityGovernance) SetPrivilegedAccess(value PrivilegedAccessRootable)() {
    err := m.GetBackingStore().Set("privilegedAccess", value)
    if err != nil {
        panic(err)
    }
}
// SetTermsOfUse sets the termsOfUse property value. The termsOfUse property
func (m *IdentityGovernance) SetTermsOfUse(value TermsOfUseContainerable)() {
    err := m.GetBackingStore().Set("termsOfUse", value)
    if err != nil {
        panic(err)
    }
}
type IdentityGovernanceable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccessReviews()(AccessReviewSetable)
    GetAppConsent()(AppConsentApprovalRouteable)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetEntitlementManagement()(EntitlementManagementable)
    GetOdataType()(*string)
    GetPrivilegedAccess()(PrivilegedAccessRootable)
    GetTermsOfUse()(TermsOfUseContainerable)
    SetAccessReviews(value AccessReviewSetable)()
    SetAppConsent(value AppConsentApprovalRouteable)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetEntitlementManagement(value EntitlementManagementable)()
    SetOdataType(value *string)()
    SetPrivilegedAccess(value PrivilegedAccessRootable)()
    SetTermsOfUse(value TermsOfUseContainerable)()
}
