package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type MultiTenantOrganizationMemberTransitionDetails struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewMultiTenantOrganizationMemberTransitionDetails instantiates a new MultiTenantOrganizationMemberTransitionDetails and sets the default values.
func NewMultiTenantOrganizationMemberTransitionDetails()(*MultiTenantOrganizationMemberTransitionDetails) {
    m := &MultiTenantOrganizationMemberTransitionDetails{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateMultiTenantOrganizationMemberTransitionDetailsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMultiTenantOrganizationMemberTransitionDetailsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMultiTenantOrganizationMemberTransitionDetails(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *MultiTenantOrganizationMemberTransitionDetails) GetAdditionalData()(map[string]any) {
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
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *MultiTenantOrganizationMemberTransitionDetails) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDesiredRole gets the desiredRole property value. Role of the tenant in the multitenant organization. The possible values are: owner, member, unknownFutureValue.
// returns a *MultiTenantOrganizationMemberRole when successful
func (m *MultiTenantOrganizationMemberTransitionDetails) GetDesiredRole()(*MultiTenantOrganizationMemberRole) {
    val, err := m.GetBackingStore().Get("desiredRole")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MultiTenantOrganizationMemberRole)
    }
    return nil
}
// GetDesiredState gets the desiredState property value. State of the tenant in the multitenant organization currently being processed. The possible values are: pending, active, removed, unknownFutureValue. Read-only.
// returns a *MultiTenantOrganizationMemberState when successful
func (m *MultiTenantOrganizationMemberTransitionDetails) GetDesiredState()(*MultiTenantOrganizationMemberState) {
    val, err := m.GetBackingStore().Get("desiredState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MultiTenantOrganizationMemberState)
    }
    return nil
}
// GetDetails gets the details property value. Details that explain the processing status if any. Read-only.
// returns a *string when successful
func (m *MultiTenantOrganizationMemberTransitionDetails) GetDetails()(*string) {
    val, err := m.GetBackingStore().Get("details")
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
func (m *MultiTenantOrganizationMemberTransitionDetails) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["desiredRole"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMultiTenantOrganizationMemberRole)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDesiredRole(val.(*MultiTenantOrganizationMemberRole))
        }
        return nil
    }
    res["desiredState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMultiTenantOrganizationMemberState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDesiredState(val.(*MultiTenantOrganizationMemberState))
        }
        return nil
    }
    res["details"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDetails(val)
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
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMultiTenantOrganizationMemberProcessingStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*MultiTenantOrganizationMemberProcessingStatus))
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *MultiTenantOrganizationMemberTransitionDetails) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStatus gets the status property value. Processing state of the asynchronous job. The possible values are: notStarted, running, succeeded, failed, unknownFutureValue. Read-only.
// returns a *MultiTenantOrganizationMemberProcessingStatus when successful
func (m *MultiTenantOrganizationMemberTransitionDetails) GetStatus()(*MultiTenantOrganizationMemberProcessingStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MultiTenantOrganizationMemberProcessingStatus)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MultiTenantOrganizationMemberTransitionDetails) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetDesiredRole() != nil {
        cast := (*m.GetDesiredRole()).String()
        err := writer.WriteStringValue("desiredRole", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetDesiredState() != nil {
        cast := (*m.GetDesiredState()).String()
        err := writer.WriteStringValue("desiredState", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("details", m.GetDetails())
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
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err := writer.WriteStringValue("status", &cast)
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
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *MultiTenantOrganizationMemberTransitionDetails) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *MultiTenantOrganizationMemberTransitionDetails) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDesiredRole sets the desiredRole property value. Role of the tenant in the multitenant organization. The possible values are: owner, member, unknownFutureValue.
func (m *MultiTenantOrganizationMemberTransitionDetails) SetDesiredRole(value *MultiTenantOrganizationMemberRole)() {
    err := m.GetBackingStore().Set("desiredRole", value)
    if err != nil {
        panic(err)
    }
}
// SetDesiredState sets the desiredState property value. State of the tenant in the multitenant organization currently being processed. The possible values are: pending, active, removed, unknownFutureValue. Read-only.
func (m *MultiTenantOrganizationMemberTransitionDetails) SetDesiredState(value *MultiTenantOrganizationMemberState)() {
    err := m.GetBackingStore().Set("desiredState", value)
    if err != nil {
        panic(err)
    }
}
// SetDetails sets the details property value. Details that explain the processing status if any. Read-only.
func (m *MultiTenantOrganizationMemberTransitionDetails) SetDetails(value *string)() {
    err := m.GetBackingStore().Set("details", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *MultiTenantOrganizationMemberTransitionDetails) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. Processing state of the asynchronous job. The possible values are: notStarted, running, succeeded, failed, unknownFutureValue. Read-only.
func (m *MultiTenantOrganizationMemberTransitionDetails) SetStatus(value *MultiTenantOrganizationMemberProcessingStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
type MultiTenantOrganizationMemberTransitionDetailsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDesiredRole()(*MultiTenantOrganizationMemberRole)
    GetDesiredState()(*MultiTenantOrganizationMemberState)
    GetDetails()(*string)
    GetOdataType()(*string)
    GetStatus()(*MultiTenantOrganizationMemberProcessingStatus)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDesiredRole(value *MultiTenantOrganizationMemberRole)()
    SetDesiredState(value *MultiTenantOrganizationMemberState)()
    SetDetails(value *string)()
    SetOdataType(value *string)()
    SetStatus(value *MultiTenantOrganizationMemberProcessingStatus)()
}
