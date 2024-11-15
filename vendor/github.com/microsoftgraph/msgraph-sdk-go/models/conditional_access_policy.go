package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ConditionalAccessPolicy struct {
    Entity
}
// NewConditionalAccessPolicy instantiates a new ConditionalAccessPolicy and sets the default values.
func NewConditionalAccessPolicy()(*ConditionalAccessPolicy) {
    m := &ConditionalAccessPolicy{
        Entity: *NewEntity(),
    }
    return m
}
// CreateConditionalAccessPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateConditionalAccessPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConditionalAccessPolicy(), nil
}
// GetConditions gets the conditions property value. The conditions property
// returns a ConditionalAccessConditionSetable when successful
func (m *ConditionalAccessPolicy) GetConditions()(ConditionalAccessConditionSetable) {
    val, err := m.GetBackingStore().Get("conditions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConditionalAccessConditionSetable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Readonly.
// returns a *Time when successful
func (m *ConditionalAccessPolicy) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDescription gets the description property value. The description property
// returns a *string when successful
func (m *ConditionalAccessPolicy) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Specifies a display name for the conditionalAccessPolicy object.
// returns a *string when successful
func (m *ConditionalAccessPolicy) GetDisplayName()(*string) {
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
func (m *ConditionalAccessPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["conditions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConditionalAccessConditionSetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConditions(val.(ConditionalAccessConditionSetable))
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
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
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
    res["grantControls"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConditionalAccessGrantControlsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGrantControls(val.(ConditionalAccessGrantControlsable))
        }
        return nil
    }
    res["modifiedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetModifiedDateTime(val)
        }
        return nil
    }
    res["sessionControls"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConditionalAccessSessionControlsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSessionControls(val.(ConditionalAccessSessionControlsable))
        }
        return nil
    }
    res["state"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseConditionalAccessPolicyState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetState(val.(*ConditionalAccessPolicyState))
        }
        return nil
    }
    res["templateId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTemplateId(val)
        }
        return nil
    }
    return res
}
// GetGrantControls gets the grantControls property value. Specifies the grant controls that must be fulfilled to pass the policy.
// returns a ConditionalAccessGrantControlsable when successful
func (m *ConditionalAccessPolicy) GetGrantControls()(ConditionalAccessGrantControlsable) {
    val, err := m.GetBackingStore().Get("grantControls")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConditionalAccessGrantControlsable)
    }
    return nil
}
// GetModifiedDateTime gets the modifiedDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Readonly.
// returns a *Time when successful
func (m *ConditionalAccessPolicy) GetModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("modifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetSessionControls gets the sessionControls property value. Specifies the session controls that are enforced after sign-in.
// returns a ConditionalAccessSessionControlsable when successful
func (m *ConditionalAccessPolicy) GetSessionControls()(ConditionalAccessSessionControlsable) {
    val, err := m.GetBackingStore().Get("sessionControls")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConditionalAccessSessionControlsable)
    }
    return nil
}
// GetState gets the state property value. The state property
// returns a *ConditionalAccessPolicyState when successful
func (m *ConditionalAccessPolicy) GetState()(*ConditionalAccessPolicyState) {
    val, err := m.GetBackingStore().Get("state")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ConditionalAccessPolicyState)
    }
    return nil
}
// GetTemplateId gets the templateId property value. The templateId property
// returns a *string when successful
func (m *ConditionalAccessPolicy) GetTemplateId()(*string) {
    val, err := m.GetBackingStore().Get("templateId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ConditionalAccessPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("conditions", m.GetConditions())
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
        err = writer.WriteStringValue("description", m.GetDescription())
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
        err = writer.WriteObjectValue("grantControls", m.GetGrantControls())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("modifiedDateTime", m.GetModifiedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("sessionControls", m.GetSessionControls())
        if err != nil {
            return err
        }
    }
    if m.GetState() != nil {
        cast := (*m.GetState()).String()
        err = writer.WriteStringValue("state", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("templateId", m.GetTemplateId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetConditions sets the conditions property value. The conditions property
func (m *ConditionalAccessPolicy) SetConditions(value ConditionalAccessConditionSetable)() {
    err := m.GetBackingStore().Set("conditions", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Readonly.
func (m *ConditionalAccessPolicy) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. The description property
func (m *ConditionalAccessPolicy) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Specifies a display name for the conditionalAccessPolicy object.
func (m *ConditionalAccessPolicy) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetGrantControls sets the grantControls property value. Specifies the grant controls that must be fulfilled to pass the policy.
func (m *ConditionalAccessPolicy) SetGrantControls(value ConditionalAccessGrantControlsable)() {
    err := m.GetBackingStore().Set("grantControls", value)
    if err != nil {
        panic(err)
    }
}
// SetModifiedDateTime sets the modifiedDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Readonly.
func (m *ConditionalAccessPolicy) SetModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("modifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetSessionControls sets the sessionControls property value. Specifies the session controls that are enforced after sign-in.
func (m *ConditionalAccessPolicy) SetSessionControls(value ConditionalAccessSessionControlsable)() {
    err := m.GetBackingStore().Set("sessionControls", value)
    if err != nil {
        panic(err)
    }
}
// SetState sets the state property value. The state property
func (m *ConditionalAccessPolicy) SetState(value *ConditionalAccessPolicyState)() {
    err := m.GetBackingStore().Set("state", value)
    if err != nil {
        panic(err)
    }
}
// SetTemplateId sets the templateId property value. The templateId property
func (m *ConditionalAccessPolicy) SetTemplateId(value *string)() {
    err := m.GetBackingStore().Set("templateId", value)
    if err != nil {
        panic(err)
    }
}
type ConditionalAccessPolicyable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetConditions()(ConditionalAccessConditionSetable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetGrantControls()(ConditionalAccessGrantControlsable)
    GetModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetSessionControls()(ConditionalAccessSessionControlsable)
    GetState()(*ConditionalAccessPolicyState)
    GetTemplateId()(*string)
    SetConditions(value ConditionalAccessConditionSetable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetGrantControls(value ConditionalAccessGrantControlsable)()
    SetModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetSessionControls(value ConditionalAccessSessionControlsable)()
    SetState(value *ConditionalAccessPolicyState)()
    SetTemplateId(value *string)()
}
