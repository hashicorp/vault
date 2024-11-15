package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CompanySubscription struct {
    Entity
}
// NewCompanySubscription instantiates a new CompanySubscription and sets the default values.
func NewCompanySubscription()(*CompanySubscription) {
    m := &CompanySubscription{
        Entity: *NewEntity(),
    }
    return m
}
// CreateCompanySubscriptionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCompanySubscriptionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCompanySubscription(), nil
}
// GetCommerceSubscriptionId gets the commerceSubscriptionId property value. The ID of this subscription in the commerce system. Alternate key.
// returns a *string when successful
func (m *CompanySubscription) GetCommerceSubscriptionId()(*string) {
    val, err := m.GetBackingStore().Get("commerceSubscriptionId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The date and time when this subscription was created. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *CompanySubscription) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CompanySubscription) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["commerceSubscriptionId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCommerceSubscriptionId(val)
        }
        return nil
    }
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
        }
        return nil
    }
    res["isTrial"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsTrial(val)
        }
        return nil
    }
    res["nextLifecycleDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNextLifecycleDateTime(val)
        }
        return nil
    }
    res["ownerId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOwnerId(val)
        }
        return nil
    }
    res["ownerTenantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOwnerTenantId(val)
        }
        return nil
    }
    res["ownerType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOwnerType(val)
        }
        return nil
    }
    res["serviceStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateServicePlanInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ServicePlanInfoable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ServicePlanInfoable)
                }
            }
            m.SetServiceStatus(res)
        }
        return nil
    }
    res["skuId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSkuId(val)
        }
        return nil
    }
    res["skuPartNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSkuPartNumber(val)
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val)
        }
        return nil
    }
    res["totalLicenses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalLicenses(val)
        }
        return nil
    }
    return res
}
// GetIsTrial gets the isTrial property value. Whether the subscription is a free trial or purchased.
// returns a *bool when successful
func (m *CompanySubscription) GetIsTrial()(*bool) {
    val, err := m.GetBackingStore().Get("isTrial")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetNextLifecycleDateTime gets the nextLifecycleDateTime property value. The date and time when the subscription will move to the next state (as defined by the status property) if not renewed by the tenant. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *CompanySubscription) GetNextLifecycleDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("nextLifecycleDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetOwnerId gets the ownerId property value. The object ID of the account admin.
// returns a *string when successful
func (m *CompanySubscription) GetOwnerId()(*string) {
    val, err := m.GetBackingStore().Get("ownerId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOwnerTenantId gets the ownerTenantId property value. The unique identifier for the Microsoft partner tenant that created the subscription on a customer tenant.
// returns a *string when successful
func (m *CompanySubscription) GetOwnerTenantId()(*string) {
    val, err := m.GetBackingStore().Get("ownerTenantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOwnerType gets the ownerType property value. Indicates the entity that ownerId belongs to, for example, 'User'.
// returns a *string when successful
func (m *CompanySubscription) GetOwnerType()(*string) {
    val, err := m.GetBackingStore().Get("ownerType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetServiceStatus gets the serviceStatus property value. The provisioning status of each service included in this subscription.
// returns a []ServicePlanInfoable when successful
func (m *CompanySubscription) GetServiceStatus()([]ServicePlanInfoable) {
    val, err := m.GetBackingStore().Get("serviceStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ServicePlanInfoable)
    }
    return nil
}
// GetSkuId gets the skuId property value. The object ID of the SKU associated with this subscription.
// returns a *string when successful
func (m *CompanySubscription) GetSkuId()(*string) {
    val, err := m.GetBackingStore().Get("skuId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSkuPartNumber gets the skuPartNumber property value. The SKU associated with this subscription.
// returns a *string when successful
func (m *CompanySubscription) GetSkuPartNumber()(*string) {
    val, err := m.GetBackingStore().Get("skuPartNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStatus gets the status property value. The status of this subscription. Possible values are: Enabled, Deleted, Suspended, Warning, LockedOut.
// returns a *string when successful
func (m *CompanySubscription) GetStatus()(*string) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTotalLicenses gets the totalLicenses property value. The number of licenses included in this subscription.
// returns a *int32 when successful
func (m *CompanySubscription) GetTotalLicenses()(*int32) {
    val, err := m.GetBackingStore().Get("totalLicenses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CompanySubscription) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("commerceSubscriptionId", m.GetCommerceSubscriptionId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isTrial", m.GetIsTrial())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("nextLifecycleDateTime", m.GetNextLifecycleDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("ownerId", m.GetOwnerId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("ownerTenantId", m.GetOwnerTenantId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("ownerType", m.GetOwnerType())
        if err != nil {
            return err
        }
    }
    if m.GetServiceStatus() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetServiceStatus()))
        for i, v := range m.GetServiceStatus() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("serviceStatus", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("skuId", m.GetSkuId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("skuPartNumber", m.GetSkuPartNumber())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("status", m.GetStatus())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("totalLicenses", m.GetTotalLicenses())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCommerceSubscriptionId sets the commerceSubscriptionId property value. The ID of this subscription in the commerce system. Alternate key.
func (m *CompanySubscription) SetCommerceSubscriptionId(value *string)() {
    err := m.GetBackingStore().Set("commerceSubscriptionId", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The date and time when this subscription was created. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *CompanySubscription) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetIsTrial sets the isTrial property value. Whether the subscription is a free trial or purchased.
func (m *CompanySubscription) SetIsTrial(value *bool)() {
    err := m.GetBackingStore().Set("isTrial", value)
    if err != nil {
        panic(err)
    }
}
// SetNextLifecycleDateTime sets the nextLifecycleDateTime property value. The date and time when the subscription will move to the next state (as defined by the status property) if not renewed by the tenant. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *CompanySubscription) SetNextLifecycleDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("nextLifecycleDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetOwnerId sets the ownerId property value. The object ID of the account admin.
func (m *CompanySubscription) SetOwnerId(value *string)() {
    err := m.GetBackingStore().Set("ownerId", value)
    if err != nil {
        panic(err)
    }
}
// SetOwnerTenantId sets the ownerTenantId property value. The unique identifier for the Microsoft partner tenant that created the subscription on a customer tenant.
func (m *CompanySubscription) SetOwnerTenantId(value *string)() {
    err := m.GetBackingStore().Set("ownerTenantId", value)
    if err != nil {
        panic(err)
    }
}
// SetOwnerType sets the ownerType property value. Indicates the entity that ownerId belongs to, for example, 'User'.
func (m *CompanySubscription) SetOwnerType(value *string)() {
    err := m.GetBackingStore().Set("ownerType", value)
    if err != nil {
        panic(err)
    }
}
// SetServiceStatus sets the serviceStatus property value. The provisioning status of each service included in this subscription.
func (m *CompanySubscription) SetServiceStatus(value []ServicePlanInfoable)() {
    err := m.GetBackingStore().Set("serviceStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetSkuId sets the skuId property value. The object ID of the SKU associated with this subscription.
func (m *CompanySubscription) SetSkuId(value *string)() {
    err := m.GetBackingStore().Set("skuId", value)
    if err != nil {
        panic(err)
    }
}
// SetSkuPartNumber sets the skuPartNumber property value. The SKU associated with this subscription.
func (m *CompanySubscription) SetSkuPartNumber(value *string)() {
    err := m.GetBackingStore().Set("skuPartNumber", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The status of this subscription. Possible values are: Enabled, Deleted, Suspended, Warning, LockedOut.
func (m *CompanySubscription) SetStatus(value *string)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalLicenses sets the totalLicenses property value. The number of licenses included in this subscription.
func (m *CompanySubscription) SetTotalLicenses(value *int32)() {
    err := m.GetBackingStore().Set("totalLicenses", value)
    if err != nil {
        panic(err)
    }
}
type CompanySubscriptionable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCommerceSubscriptionId()(*string)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetIsTrial()(*bool)
    GetNextLifecycleDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetOwnerId()(*string)
    GetOwnerTenantId()(*string)
    GetOwnerType()(*string)
    GetServiceStatus()([]ServicePlanInfoable)
    GetSkuId()(*string)
    GetSkuPartNumber()(*string)
    GetStatus()(*string)
    GetTotalLicenses()(*int32)
    SetCommerceSubscriptionId(value *string)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetIsTrial(value *bool)()
    SetNextLifecycleDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetOwnerId(value *string)()
    SetOwnerTenantId(value *string)()
    SetOwnerType(value *string)()
    SetServiceStatus(value []ServicePlanInfoable)()
    SetSkuId(value *string)()
    SetSkuPartNumber(value *string)()
    SetStatus(value *string)()
    SetTotalLicenses(value *int32)()
}
