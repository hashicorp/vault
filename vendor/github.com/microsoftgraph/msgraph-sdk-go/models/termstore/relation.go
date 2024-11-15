package termstore

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type Relation struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
}
// NewRelation instantiates a new Relation and sets the default values.
func NewRelation()(*Relation) {
    m := &Relation{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateRelationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRelationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRelation(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Relation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["fromTerm"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTermFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFromTerm(val.(Termable))
        }
        return nil
    }
    res["relationship"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseRelationType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRelationship(val.(*RelationType))
        }
        return nil
    }
    res["set"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSet(val.(Setable))
        }
        return nil
    }
    res["toTerm"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTermFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetToTerm(val.(Termable))
        }
        return nil
    }
    return res
}
// GetFromTerm gets the fromTerm property value. The from [term] of the relation. The term from which the relationship is defined. A null value would indicate the relation is directly with the [set].
// returns a Termable when successful
func (m *Relation) GetFromTerm()(Termable) {
    val, err := m.GetBackingStore().Get("fromTerm")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Termable)
    }
    return nil
}
// GetRelationship gets the relationship property value. The type of relation. Possible values are: pin, reuse.
// returns a *RelationType when successful
func (m *Relation) GetRelationship()(*RelationType) {
    val, err := m.GetBackingStore().Get("relationship")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RelationType)
    }
    return nil
}
// GetSet gets the set property value. The [set] in which the relation is relevant.
// returns a Setable when successful
func (m *Relation) GetSet()(Setable) {
    val, err := m.GetBackingStore().Get("set")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Setable)
    }
    return nil
}
// GetToTerm gets the toTerm property value. The to [term] of the relation. The term to which the relationship is defined.
// returns a Termable when successful
func (m *Relation) GetToTerm()(Termable) {
    val, err := m.GetBackingStore().Get("toTerm")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Termable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Relation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("fromTerm", m.GetFromTerm())
        if err != nil {
            return err
        }
    }
    if m.GetRelationship() != nil {
        cast := (*m.GetRelationship()).String()
        err = writer.WriteStringValue("relationship", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("set", m.GetSet())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("toTerm", m.GetToTerm())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFromTerm sets the fromTerm property value. The from [term] of the relation. The term from which the relationship is defined. A null value would indicate the relation is directly with the [set].
func (m *Relation) SetFromTerm(value Termable)() {
    err := m.GetBackingStore().Set("fromTerm", value)
    if err != nil {
        panic(err)
    }
}
// SetRelationship sets the relationship property value. The type of relation. Possible values are: pin, reuse.
func (m *Relation) SetRelationship(value *RelationType)() {
    err := m.GetBackingStore().Set("relationship", value)
    if err != nil {
        panic(err)
    }
}
// SetSet sets the set property value. The [set] in which the relation is relevant.
func (m *Relation) SetSet(value Setable)() {
    err := m.GetBackingStore().Set("set", value)
    if err != nil {
        panic(err)
    }
}
// SetToTerm sets the toTerm property value. The to [term] of the relation. The term to which the relationship is defined.
func (m *Relation) SetToTerm(value Termable)() {
    err := m.GetBackingStore().Set("toTerm", value)
    if err != nil {
        panic(err)
    }
}
type Relationable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetFromTerm()(Termable)
    GetRelationship()(*RelationType)
    GetSet()(Setable)
    GetToTerm()(Termable)
    SetFromTerm(value Termable)()
    SetRelationship(value *RelationType)()
    SetSet(value Setable)()
    SetToTerm(value Termable)()
}
