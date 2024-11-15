package termstore

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type Set struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
}
// NewSet instantiates a new Set and sets the default values.
func NewSet()(*Set) {
    m := &Set{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateSetFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSetFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSet(), nil
}
// GetChildren gets the children property value. Children terms of set in term [store].
// returns a []Termable when successful
func (m *Set) GetChildren()([]Termable) {
    val, err := m.GetBackingStore().Get("children")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Termable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Date and time of set creation. Read-only.
// returns a *Time when successful
func (m *Set) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDescription gets the description property value. Description that gives details on the term usage.
// returns a *string when successful
func (m *Set) GetDescription()(*string) {
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
func (m *Set) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["children"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTermFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Termable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Termable)
                }
            }
            m.SetChildren(res)
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
    res["localizedNames"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateLocalizedNameFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]LocalizedNameable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(LocalizedNameable)
                }
            }
            m.SetLocalizedNames(res)
        }
        return nil
    }
    res["parentGroup"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateGroupFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParentGroup(val.(Groupable))
        }
        return nil
    }
    res["properties"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateKeyValueFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.KeyValueable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.KeyValueable)
                }
            }
            m.SetProperties(res)
        }
        return nil
    }
    res["relations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRelationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Relationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Relationable)
                }
            }
            m.SetRelations(res)
        }
        return nil
    }
    res["terms"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTermFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Termable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Termable)
                }
            }
            m.SetTerms(res)
        }
        return nil
    }
    return res
}
// GetLocalizedNames gets the localizedNames property value. Name of the set for each languageTag.
// returns a []LocalizedNameable when successful
func (m *Set) GetLocalizedNames()([]LocalizedNameable) {
    val, err := m.GetBackingStore().Get("localizedNames")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]LocalizedNameable)
    }
    return nil
}
// GetParentGroup gets the parentGroup property value. The parentGroup property
// returns a Groupable when successful
func (m *Set) GetParentGroup()(Groupable) {
    val, err := m.GetBackingStore().Get("parentGroup")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Groupable)
    }
    return nil
}
// GetProperties gets the properties property value. Custom properties for the set.
// returns a []KeyValueable when successful
func (m *Set) GetProperties()([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.KeyValueable) {
    val, err := m.GetBackingStore().Get("properties")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.KeyValueable)
    }
    return nil
}
// GetRelations gets the relations property value. Indicates which terms have been pinned or reused directly under the set.
// returns a []Relationable when successful
func (m *Set) GetRelations()([]Relationable) {
    val, err := m.GetBackingStore().Get("relations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Relationable)
    }
    return nil
}
// GetTerms gets the terms property value. All the terms under the set.
// returns a []Termable when successful
func (m *Set) GetTerms()([]Termable) {
    val, err := m.GetBackingStore().Get("terms")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Termable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Set) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetChildren() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetChildren()))
        for i, v := range m.GetChildren() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("children", cast)
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
    if m.GetLocalizedNames() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLocalizedNames()))
        for i, v := range m.GetLocalizedNames() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("localizedNames", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("parentGroup", m.GetParentGroup())
        if err != nil {
            return err
        }
    }
    if m.GetProperties() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetProperties()))
        for i, v := range m.GetProperties() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("properties", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRelations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRelations()))
        for i, v := range m.GetRelations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("relations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTerms() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTerms()))
        for i, v := range m.GetTerms() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("terms", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetChildren sets the children property value. Children terms of set in term [store].
func (m *Set) SetChildren(value []Termable)() {
    err := m.GetBackingStore().Set("children", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Date and time of set creation. Read-only.
func (m *Set) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. Description that gives details on the term usage.
func (m *Set) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetLocalizedNames sets the localizedNames property value. Name of the set for each languageTag.
func (m *Set) SetLocalizedNames(value []LocalizedNameable)() {
    err := m.GetBackingStore().Set("localizedNames", value)
    if err != nil {
        panic(err)
    }
}
// SetParentGroup sets the parentGroup property value. The parentGroup property
func (m *Set) SetParentGroup(value Groupable)() {
    err := m.GetBackingStore().Set("parentGroup", value)
    if err != nil {
        panic(err)
    }
}
// SetProperties sets the properties property value. Custom properties for the set.
func (m *Set) SetProperties(value []iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.KeyValueable)() {
    err := m.GetBackingStore().Set("properties", value)
    if err != nil {
        panic(err)
    }
}
// SetRelations sets the relations property value. Indicates which terms have been pinned or reused directly under the set.
func (m *Set) SetRelations(value []Relationable)() {
    err := m.GetBackingStore().Set("relations", value)
    if err != nil {
        panic(err)
    }
}
// SetTerms sets the terms property value. All the terms under the set.
func (m *Set) SetTerms(value []Termable)() {
    err := m.GetBackingStore().Set("terms", value)
    if err != nil {
        panic(err)
    }
}
type Setable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetChildren()([]Termable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDescription()(*string)
    GetLocalizedNames()([]LocalizedNameable)
    GetParentGroup()(Groupable)
    GetProperties()([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.KeyValueable)
    GetRelations()([]Relationable)
    GetTerms()([]Termable)
    SetChildren(value []Termable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDescription(value *string)()
    SetLocalizedNames(value []LocalizedNameable)()
    SetParentGroup(value Groupable)()
    SetProperties(value []iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.KeyValueable)()
    SetRelations(value []Relationable)()
    SetTerms(value []Termable)()
}
