package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type UnifiedRoleManagementPolicyRuleTarget struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewUnifiedRoleManagementPolicyRuleTarget instantiates a new UnifiedRoleManagementPolicyRuleTarget and sets the default values.
func NewUnifiedRoleManagementPolicyRuleTarget()(*UnifiedRoleManagementPolicyRuleTarget) {
    m := &UnifiedRoleManagementPolicyRuleTarget{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateUnifiedRoleManagementPolicyRuleTargetFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUnifiedRoleManagementPolicyRuleTargetFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUnifiedRoleManagementPolicyRuleTarget(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *UnifiedRoleManagementPolicyRuleTarget) GetAdditionalData()(map[string]any) {
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
func (m *UnifiedRoleManagementPolicyRuleTarget) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCaller gets the caller property value. The type of caller that's the target of the policy rule. Allowed values are: None, Admin, EndUser.
// returns a *string when successful
func (m *UnifiedRoleManagementPolicyRuleTarget) GetCaller()(*string) {
    val, err := m.GetBackingStore().Get("caller")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEnforcedSettings gets the enforcedSettings property value. The list of role settings that are enforced and cannot be overridden by child scopes. Use All for all settings.
// returns a []string when successful
func (m *UnifiedRoleManagementPolicyRuleTarget) GetEnforcedSettings()([]string) {
    val, err := m.GetBackingStore().Get("enforcedSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UnifiedRoleManagementPolicyRuleTarget) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["caller"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCaller(val)
        }
        return nil
    }
    res["enforcedSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetEnforcedSettings(res)
        }
        return nil
    }
    res["inheritableSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetInheritableSettings(res)
        }
        return nil
    }
    res["level"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLevel(val)
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
    res["operations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseUnifiedRoleManagementPolicyRuleTargetOperations)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UnifiedRoleManagementPolicyRuleTargetOperations, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*UnifiedRoleManagementPolicyRuleTargetOperations))
                }
            }
            m.SetOperations(res)
        }
        return nil
    }
    res["targetObjects"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DirectoryObjectable)
                }
            }
            m.SetTargetObjects(res)
        }
        return nil
    }
    return res
}
// GetInheritableSettings gets the inheritableSettings property value. The list of role settings that can be inherited by child scopes. Use All for all settings.
// returns a []string when successful
func (m *UnifiedRoleManagementPolicyRuleTarget) GetInheritableSettings()([]string) {
    val, err := m.GetBackingStore().Get("inheritableSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetLevel gets the level property value. The role assignment type that's the target of policy rule. Allowed values are: Eligibility, Assignment.
// returns a *string when successful
func (m *UnifiedRoleManagementPolicyRuleTarget) GetLevel()(*string) {
    val, err := m.GetBackingStore().Get("level")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *UnifiedRoleManagementPolicyRuleTarget) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOperations gets the operations property value. The role management operations that are the target of the policy rule. Allowed values are: All, Activate, Deactivate, Assign, Update, Remove, Extend, Renew.
// returns a []UnifiedRoleManagementPolicyRuleTargetOperations when successful
func (m *UnifiedRoleManagementPolicyRuleTarget) GetOperations()([]UnifiedRoleManagementPolicyRuleTargetOperations) {
    val, err := m.GetBackingStore().Get("operations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UnifiedRoleManagementPolicyRuleTargetOperations)
    }
    return nil
}
// GetTargetObjects gets the targetObjects property value. The targetObjects property
// returns a []DirectoryObjectable when successful
func (m *UnifiedRoleManagementPolicyRuleTarget) GetTargetObjects()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("targetObjects")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UnifiedRoleManagementPolicyRuleTarget) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("caller", m.GetCaller())
        if err != nil {
            return err
        }
    }
    if m.GetEnforcedSettings() != nil {
        err := writer.WriteCollectionOfStringValues("enforcedSettings", m.GetEnforcedSettings())
        if err != nil {
            return err
        }
    }
    if m.GetInheritableSettings() != nil {
        err := writer.WriteCollectionOfStringValues("inheritableSettings", m.GetInheritableSettings())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("level", m.GetLevel())
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
    if m.GetOperations() != nil {
        err := writer.WriteCollectionOfStringValues("operations", SerializeUnifiedRoleManagementPolicyRuleTargetOperations(m.GetOperations()))
        if err != nil {
            return err
        }
    }
    if m.GetTargetObjects() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTargetObjects()))
        for i, v := range m.GetTargetObjects() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("targetObjects", cast)
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
func (m *UnifiedRoleManagementPolicyRuleTarget) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *UnifiedRoleManagementPolicyRuleTarget) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCaller sets the caller property value. The type of caller that's the target of the policy rule. Allowed values are: None, Admin, EndUser.
func (m *UnifiedRoleManagementPolicyRuleTarget) SetCaller(value *string)() {
    err := m.GetBackingStore().Set("caller", value)
    if err != nil {
        panic(err)
    }
}
// SetEnforcedSettings sets the enforcedSettings property value. The list of role settings that are enforced and cannot be overridden by child scopes. Use All for all settings.
func (m *UnifiedRoleManagementPolicyRuleTarget) SetEnforcedSettings(value []string)() {
    err := m.GetBackingStore().Set("enforcedSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetInheritableSettings sets the inheritableSettings property value. The list of role settings that can be inherited by child scopes. Use All for all settings.
func (m *UnifiedRoleManagementPolicyRuleTarget) SetInheritableSettings(value []string)() {
    err := m.GetBackingStore().Set("inheritableSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetLevel sets the level property value. The role assignment type that's the target of policy rule. Allowed values are: Eligibility, Assignment.
func (m *UnifiedRoleManagementPolicyRuleTarget) SetLevel(value *string)() {
    err := m.GetBackingStore().Set("level", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *UnifiedRoleManagementPolicyRuleTarget) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetOperations sets the operations property value. The role management operations that are the target of the policy rule. Allowed values are: All, Activate, Deactivate, Assign, Update, Remove, Extend, Renew.
func (m *UnifiedRoleManagementPolicyRuleTarget) SetOperations(value []UnifiedRoleManagementPolicyRuleTargetOperations)() {
    err := m.GetBackingStore().Set("operations", value)
    if err != nil {
        panic(err)
    }
}
// SetTargetObjects sets the targetObjects property value. The targetObjects property
func (m *UnifiedRoleManagementPolicyRuleTarget) SetTargetObjects(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("targetObjects", value)
    if err != nil {
        panic(err)
    }
}
type UnifiedRoleManagementPolicyRuleTargetable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCaller()(*string)
    GetEnforcedSettings()([]string)
    GetInheritableSettings()([]string)
    GetLevel()(*string)
    GetOdataType()(*string)
    GetOperations()([]UnifiedRoleManagementPolicyRuleTargetOperations)
    GetTargetObjects()([]DirectoryObjectable)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCaller(value *string)()
    SetEnforcedSettings(value []string)()
    SetInheritableSettings(value []string)()
    SetLevel(value *string)()
    SetOdataType(value *string)()
    SetOperations(value []UnifiedRoleManagementPolicyRuleTargetOperations)()
    SetTargetObjects(value []DirectoryObjectable)()
}
