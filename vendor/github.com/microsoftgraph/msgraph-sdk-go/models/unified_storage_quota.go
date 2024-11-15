package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UnifiedStorageQuota struct {
    Entity
}
// NewUnifiedStorageQuota instantiates a new UnifiedStorageQuota and sets the default values.
func NewUnifiedStorageQuota()(*UnifiedStorageQuota) {
    m := &UnifiedStorageQuota{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUnifiedStorageQuotaFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUnifiedStorageQuotaFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUnifiedStorageQuota(), nil
}
// GetDeleted gets the deleted property value. The deleted property
// returns a *int64 when successful
func (m *UnifiedStorageQuota) GetDeleted()(*int64) {
    val, err := m.GetBackingStore().Get("deleted")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UnifiedStorageQuota) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["deleted"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeleted(val)
        }
        return nil
    }
    res["manageWebUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManageWebUrl(val)
        }
        return nil
    }
    res["remaining"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRemaining(val)
        }
        return nil
    }
    res["services"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateServiceStorageQuotaBreakdownFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ServiceStorageQuotaBreakdownable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ServiceStorageQuotaBreakdownable)
                }
            }
            m.SetServices(res)
        }
        return nil
    }
    res["state"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetState(val)
        }
        return nil
    }
    res["total"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotal(val)
        }
        return nil
    }
    res["used"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUsed(val)
        }
        return nil
    }
    return res
}
// GetManageWebUrl gets the manageWebUrl property value. The manageWebUrl property
// returns a *string when successful
func (m *UnifiedStorageQuota) GetManageWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("manageWebUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRemaining gets the remaining property value. The remaining property
// returns a *int64 when successful
func (m *UnifiedStorageQuota) GetRemaining()(*int64) {
    val, err := m.GetBackingStore().Get("remaining")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetServices gets the services property value. The services property
// returns a []ServiceStorageQuotaBreakdownable when successful
func (m *UnifiedStorageQuota) GetServices()([]ServiceStorageQuotaBreakdownable) {
    val, err := m.GetBackingStore().Get("services")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ServiceStorageQuotaBreakdownable)
    }
    return nil
}
// GetState gets the state property value. The state property
// returns a *string when successful
func (m *UnifiedStorageQuota) GetState()(*string) {
    val, err := m.GetBackingStore().Get("state")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTotal gets the total property value. The total property
// returns a *int64 when successful
func (m *UnifiedStorageQuota) GetTotal()(*int64) {
    val, err := m.GetBackingStore().Get("total")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetUsed gets the used property value. The used property
// returns a *int64 when successful
func (m *UnifiedStorageQuota) GetUsed()(*int64) {
    val, err := m.GetBackingStore().Get("used")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UnifiedStorageQuota) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt64Value("deleted", m.GetDeleted())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("manageWebUrl", m.GetManageWebUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("remaining", m.GetRemaining())
        if err != nil {
            return err
        }
    }
    if m.GetServices() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetServices()))
        for i, v := range m.GetServices() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("services", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("state", m.GetState())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("total", m.GetTotal())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("used", m.GetUsed())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDeleted sets the deleted property value. The deleted property
func (m *UnifiedStorageQuota) SetDeleted(value *int64)() {
    err := m.GetBackingStore().Set("deleted", value)
    if err != nil {
        panic(err)
    }
}
// SetManageWebUrl sets the manageWebUrl property value. The manageWebUrl property
func (m *UnifiedStorageQuota) SetManageWebUrl(value *string)() {
    err := m.GetBackingStore().Set("manageWebUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetRemaining sets the remaining property value. The remaining property
func (m *UnifiedStorageQuota) SetRemaining(value *int64)() {
    err := m.GetBackingStore().Set("remaining", value)
    if err != nil {
        panic(err)
    }
}
// SetServices sets the services property value. The services property
func (m *UnifiedStorageQuota) SetServices(value []ServiceStorageQuotaBreakdownable)() {
    err := m.GetBackingStore().Set("services", value)
    if err != nil {
        panic(err)
    }
}
// SetState sets the state property value. The state property
func (m *UnifiedStorageQuota) SetState(value *string)() {
    err := m.GetBackingStore().Set("state", value)
    if err != nil {
        panic(err)
    }
}
// SetTotal sets the total property value. The total property
func (m *UnifiedStorageQuota) SetTotal(value *int64)() {
    err := m.GetBackingStore().Set("total", value)
    if err != nil {
        panic(err)
    }
}
// SetUsed sets the used property value. The used property
func (m *UnifiedStorageQuota) SetUsed(value *int64)() {
    err := m.GetBackingStore().Set("used", value)
    if err != nil {
        panic(err)
    }
}
type UnifiedStorageQuotaable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDeleted()(*int64)
    GetManageWebUrl()(*string)
    GetRemaining()(*int64)
    GetServices()([]ServiceStorageQuotaBreakdownable)
    GetState()(*string)
    GetTotal()(*int64)
    GetUsed()(*int64)
    SetDeleted(value *int64)()
    SetManageWebUrl(value *string)()
    SetRemaining(value *int64)()
    SetServices(value []ServiceStorageQuotaBreakdownable)()
    SetState(value *string)()
    SetTotal(value *int64)()
    SetUsed(value *int64)()
}
