package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AdminConsentRequestPolicy struct {
    Entity
}
// NewAdminConsentRequestPolicy instantiates a new AdminConsentRequestPolicy and sets the default values.
func NewAdminConsentRequestPolicy()(*AdminConsentRequestPolicy) {
    m := &AdminConsentRequestPolicy{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAdminConsentRequestPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAdminConsentRequestPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAdminConsentRequestPolicy(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AdminConsentRequestPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["isEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsEnabled(val)
        }
        return nil
    }
    res["notifyReviewers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotifyReviewers(val)
        }
        return nil
    }
    res["remindersEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRemindersEnabled(val)
        }
        return nil
    }
    res["requestDurationInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRequestDurationInDays(val)
        }
        return nil
    }
    res["reviewers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessReviewReviewerScopeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessReviewReviewerScopeable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessReviewReviewerScopeable)
                }
            }
            m.SetReviewers(res)
        }
        return nil
    }
    res["version"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVersion(val)
        }
        return nil
    }
    return res
}
// GetIsEnabled gets the isEnabled property value. Specifies whether the admin consent request feature is enabled or disabled. Required.
// returns a *bool when successful
func (m *AdminConsentRequestPolicy) GetIsEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetNotifyReviewers gets the notifyReviewers property value. Specifies whether reviewers will receive notifications. Required.
// returns a *bool when successful
func (m *AdminConsentRequestPolicy) GetNotifyReviewers()(*bool) {
    val, err := m.GetBackingStore().Get("notifyReviewers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetRemindersEnabled gets the remindersEnabled property value. Specifies whether reviewers will receive reminder emails. Required.
// returns a *bool when successful
func (m *AdminConsentRequestPolicy) GetRemindersEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("remindersEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetRequestDurationInDays gets the requestDurationInDays property value. Specifies the duration the request is active before it automatically expires if no decision is applied.
// returns a *int32 when successful
func (m *AdminConsentRequestPolicy) GetRequestDurationInDays()(*int32) {
    val, err := m.GetBackingStore().Get("requestDurationInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetReviewers gets the reviewers property value. The list of reviewers for the admin consent. Required.
// returns a []AccessReviewReviewerScopeable when successful
func (m *AdminConsentRequestPolicy) GetReviewers()([]AccessReviewReviewerScopeable) {
    val, err := m.GetBackingStore().Get("reviewers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessReviewReviewerScopeable)
    }
    return nil
}
// GetVersion gets the version property value. Specifies the version of this policy. When the policy is updated, this version is updated. Read-only.
// returns a *int32 when successful
func (m *AdminConsentRequestPolicy) GetVersion()(*int32) {
    val, err := m.GetBackingStore().Get("version")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AdminConsentRequestPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("isEnabled", m.GetIsEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("notifyReviewers", m.GetNotifyReviewers())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("remindersEnabled", m.GetRemindersEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("requestDurationInDays", m.GetRequestDurationInDays())
        if err != nil {
            return err
        }
    }
    if m.GetReviewers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetReviewers()))
        for i, v := range m.GetReviewers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("reviewers", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("version", m.GetVersion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetIsEnabled sets the isEnabled property value. Specifies whether the admin consent request feature is enabled or disabled. Required.
func (m *AdminConsentRequestPolicy) SetIsEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetNotifyReviewers sets the notifyReviewers property value. Specifies whether reviewers will receive notifications. Required.
func (m *AdminConsentRequestPolicy) SetNotifyReviewers(value *bool)() {
    err := m.GetBackingStore().Set("notifyReviewers", value)
    if err != nil {
        panic(err)
    }
}
// SetRemindersEnabled sets the remindersEnabled property value. Specifies whether reviewers will receive reminder emails. Required.
func (m *AdminConsentRequestPolicy) SetRemindersEnabled(value *bool)() {
    err := m.GetBackingStore().Set("remindersEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetRequestDurationInDays sets the requestDurationInDays property value. Specifies the duration the request is active before it automatically expires if no decision is applied.
func (m *AdminConsentRequestPolicy) SetRequestDurationInDays(value *int32)() {
    err := m.GetBackingStore().Set("requestDurationInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetReviewers sets the reviewers property value. The list of reviewers for the admin consent. Required.
func (m *AdminConsentRequestPolicy) SetReviewers(value []AccessReviewReviewerScopeable)() {
    err := m.GetBackingStore().Set("reviewers", value)
    if err != nil {
        panic(err)
    }
}
// SetVersion sets the version property value. Specifies the version of this policy. When the policy is updated, this version is updated. Read-only.
func (m *AdminConsentRequestPolicy) SetVersion(value *int32)() {
    err := m.GetBackingStore().Set("version", value)
    if err != nil {
        panic(err)
    }
}
type AdminConsentRequestPolicyable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetIsEnabled()(*bool)
    GetNotifyReviewers()(*bool)
    GetRemindersEnabled()(*bool)
    GetRequestDurationInDays()(*int32)
    GetReviewers()([]AccessReviewReviewerScopeable)
    GetVersion()(*int32)
    SetIsEnabled(value *bool)()
    SetNotifyReviewers(value *bool)()
    SetRemindersEnabled(value *bool)()
    SetRequestDurationInDays(value *int32)()
    SetReviewers(value []AccessReviewReviewerScopeable)()
    SetVersion(value *int32)()
}
