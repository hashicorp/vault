package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CustomSecurityAttributeDefinition struct {
    Entity
}
// NewCustomSecurityAttributeDefinition instantiates a new CustomSecurityAttributeDefinition and sets the default values.
func NewCustomSecurityAttributeDefinition()(*CustomSecurityAttributeDefinition) {
    m := &CustomSecurityAttributeDefinition{
        Entity: *NewEntity(),
    }
    return m
}
// CreateCustomSecurityAttributeDefinitionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCustomSecurityAttributeDefinitionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCustomSecurityAttributeDefinition(), nil
}
// GetAllowedValues gets the allowedValues property value. Values that are predefined for this custom security attribute. This navigation property is not returned by default and must be specified in an $expand query. For example, /directory/customSecurityAttributeDefinitions?$expand=allowedValues.
// returns a []AllowedValueable when successful
func (m *CustomSecurityAttributeDefinition) GetAllowedValues()([]AllowedValueable) {
    val, err := m.GetBackingStore().Get("allowedValues")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AllowedValueable)
    }
    return nil
}
// GetAttributeSet gets the attributeSet property value. Name of the attribute set. Case insensitive.
// returns a *string when successful
func (m *CustomSecurityAttributeDefinition) GetAttributeSet()(*string) {
    val, err := m.GetBackingStore().Get("attributeSet")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDescription gets the description property value. Description of the custom security attribute. Can be up to 128 characters long and include Unicode characters. Can be changed later.
// returns a *string when successful
func (m *CustomSecurityAttributeDefinition) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
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
func (m *CustomSecurityAttributeDefinition) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["allowedValues"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAllowedValueFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AllowedValueable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AllowedValueable)
                }
            }
            m.SetAllowedValues(res)
        }
        return nil
    }
    res["attributeSet"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAttributeSet(val)
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
    res["isCollection"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsCollection(val)
        }
        return nil
    }
    res["isSearchable"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSearchable(val)
        }
        return nil
    }
    res["name"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetName(val)
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
    res["type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTypeEscaped(val)
        }
        return nil
    }
    res["usePreDefinedValuesOnly"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUsePreDefinedValuesOnly(val)
        }
        return nil
    }
    return res
}
// GetIsCollection gets the isCollection property value. Indicates whether multiple values can be assigned to the custom security attribute. Cannot be changed later. If type is set to Boolean, isCollection cannot be set to true.
// returns a *bool when successful
func (m *CustomSecurityAttributeDefinition) GetIsCollection()(*bool) {
    val, err := m.GetBackingStore().Get("isCollection")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsSearchable gets the isSearchable property value. Indicates whether custom security attribute values are indexed for searching on objects that are assigned attribute values. Cannot be changed later.
// returns a *bool when successful
func (m *CustomSecurityAttributeDefinition) GetIsSearchable()(*bool) {
    val, err := m.GetBackingStore().Get("isSearchable")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetName gets the name property value. Name of the custom security attribute. Must be unique within an attribute set. Can be up to 32 characters long and include Unicode characters. Cannot contain spaces or special characters. Cannot be changed later. Case insensitive.
// returns a *string when successful
func (m *CustomSecurityAttributeDefinition) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStatus gets the status property value. Specifies whether the custom security attribute is active or deactivated. Acceptable values are: Available and Deprecated. Can be changed later.
// returns a *string when successful
func (m *CustomSecurityAttributeDefinition) GetStatus()(*string) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTypeEscaped gets the type property value. Data type for the custom security attribute values. Supported types are: Boolean, Integer, and String. Cannot be changed later.
// returns a *string when successful
func (m *CustomSecurityAttributeDefinition) GetTypeEscaped()(*string) {
    val, err := m.GetBackingStore().Get("typeEscaped")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUsePreDefinedValuesOnly gets the usePreDefinedValuesOnly property value. Indicates whether only predefined values can be assigned to the custom security attribute. If set to false, free-form values are allowed. Can later be changed from true to false, but cannot be changed from false to true. If type is set to Boolean, usePreDefinedValuesOnly cannot be set to true.
// returns a *bool when successful
func (m *CustomSecurityAttributeDefinition) GetUsePreDefinedValuesOnly()(*bool) {
    val, err := m.GetBackingStore().Get("usePreDefinedValuesOnly")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CustomSecurityAttributeDefinition) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAllowedValues() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAllowedValues()))
        for i, v := range m.GetAllowedValues() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("allowedValues", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("attributeSet", m.GetAttributeSet())
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
        err = writer.WriteBoolValue("isCollection", m.GetIsCollection())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isSearchable", m.GetIsSearchable())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("name", m.GetName())
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
        err = writer.WriteStringValue("type", m.GetTypeEscaped())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("usePreDefinedValuesOnly", m.GetUsePreDefinedValuesOnly())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowedValues sets the allowedValues property value. Values that are predefined for this custom security attribute. This navigation property is not returned by default and must be specified in an $expand query. For example, /directory/customSecurityAttributeDefinitions?$expand=allowedValues.
func (m *CustomSecurityAttributeDefinition) SetAllowedValues(value []AllowedValueable)() {
    err := m.GetBackingStore().Set("allowedValues", value)
    if err != nil {
        panic(err)
    }
}
// SetAttributeSet sets the attributeSet property value. Name of the attribute set. Case insensitive.
func (m *CustomSecurityAttributeDefinition) SetAttributeSet(value *string)() {
    err := m.GetBackingStore().Set("attributeSet", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. Description of the custom security attribute. Can be up to 128 characters long and include Unicode characters. Can be changed later.
func (m *CustomSecurityAttributeDefinition) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetIsCollection sets the isCollection property value. Indicates whether multiple values can be assigned to the custom security attribute. Cannot be changed later. If type is set to Boolean, isCollection cannot be set to true.
func (m *CustomSecurityAttributeDefinition) SetIsCollection(value *bool)() {
    err := m.GetBackingStore().Set("isCollection", value)
    if err != nil {
        panic(err)
    }
}
// SetIsSearchable sets the isSearchable property value. Indicates whether custom security attribute values are indexed for searching on objects that are assigned attribute values. Cannot be changed later.
func (m *CustomSecurityAttributeDefinition) SetIsSearchable(value *bool)() {
    err := m.GetBackingStore().Set("isSearchable", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. Name of the custom security attribute. Must be unique within an attribute set. Can be up to 32 characters long and include Unicode characters. Cannot contain spaces or special characters. Cannot be changed later. Case insensitive.
func (m *CustomSecurityAttributeDefinition) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. Specifies whether the custom security attribute is active or deactivated. Acceptable values are: Available and Deprecated. Can be changed later.
func (m *CustomSecurityAttributeDefinition) SetStatus(value *string)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
// SetTypeEscaped sets the type property value. Data type for the custom security attribute values. Supported types are: Boolean, Integer, and String. Cannot be changed later.
func (m *CustomSecurityAttributeDefinition) SetTypeEscaped(value *string)() {
    err := m.GetBackingStore().Set("typeEscaped", value)
    if err != nil {
        panic(err)
    }
}
// SetUsePreDefinedValuesOnly sets the usePreDefinedValuesOnly property value. Indicates whether only predefined values can be assigned to the custom security attribute. If set to false, free-form values are allowed. Can later be changed from true to false, but cannot be changed from false to true. If type is set to Boolean, usePreDefinedValuesOnly cannot be set to true.
func (m *CustomSecurityAttributeDefinition) SetUsePreDefinedValuesOnly(value *bool)() {
    err := m.GetBackingStore().Set("usePreDefinedValuesOnly", value)
    if err != nil {
        panic(err)
    }
}
type CustomSecurityAttributeDefinitionable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowedValues()([]AllowedValueable)
    GetAttributeSet()(*string)
    GetDescription()(*string)
    GetIsCollection()(*bool)
    GetIsSearchable()(*bool)
    GetName()(*string)
    GetStatus()(*string)
    GetTypeEscaped()(*string)
    GetUsePreDefinedValuesOnly()(*bool)
    SetAllowedValues(value []AllowedValueable)()
    SetAttributeSet(value *string)()
    SetDescription(value *string)()
    SetIsCollection(value *bool)()
    SetIsSearchable(value *bool)()
    SetName(value *string)()
    SetStatus(value *string)()
    SetTypeEscaped(value *string)()
    SetUsePreDefinedValuesOnly(value *bool)()
}
