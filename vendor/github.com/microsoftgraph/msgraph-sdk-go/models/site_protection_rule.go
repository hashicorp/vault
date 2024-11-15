package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SiteProtectionRule struct {
    ProtectionRuleBase
}
// NewSiteProtectionRule instantiates a new SiteProtectionRule and sets the default values.
func NewSiteProtectionRule()(*SiteProtectionRule) {
    m := &SiteProtectionRule{
        ProtectionRuleBase: *NewProtectionRuleBase(),
    }
    odataTypeValue := "#microsoft.graph.siteProtectionRule"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateSiteProtectionRuleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSiteProtectionRuleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSiteProtectionRule(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SiteProtectionRule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ProtectionRuleBase.GetFieldDeserializers()
    res["siteExpression"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSiteExpression(val)
        }
        return nil
    }
    return res
}
// GetSiteExpression gets the siteExpression property value. Contains a site expression. For examples, see siteExpression example.
// returns a *string when successful
func (m *SiteProtectionRule) GetSiteExpression()(*string) {
    val, err := m.GetBackingStore().Get("siteExpression")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SiteProtectionRule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ProtectionRuleBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("siteExpression", m.GetSiteExpression())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetSiteExpression sets the siteExpression property value. Contains a site expression. For examples, see siteExpression example.
func (m *SiteProtectionRule) SetSiteExpression(value *string)() {
    err := m.GetBackingStore().Set("siteExpression", value)
    if err != nil {
        panic(err)
    }
}
type SiteProtectionRuleable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    ProtectionRuleBaseable
    GetSiteExpression()(*string)
    SetSiteExpression(value *string)()
}
