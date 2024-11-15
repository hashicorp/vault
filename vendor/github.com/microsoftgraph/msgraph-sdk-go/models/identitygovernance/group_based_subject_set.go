package identitygovernance

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type GroupBasedSubjectSet struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SubjectSet
}
// NewGroupBasedSubjectSet instantiates a new GroupBasedSubjectSet and sets the default values.
func NewGroupBasedSubjectSet()(*GroupBasedSubjectSet) {
    m := &GroupBasedSubjectSet{
        SubjectSet: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewSubjectSet(),
    }
    odataTypeValue := "#microsoft.graph.identityGovernance.groupBasedSubjectSet"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateGroupBasedSubjectSetFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateGroupBasedSubjectSetFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewGroupBasedSubjectSet(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *GroupBasedSubjectSet) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.SubjectSet.GetFieldDeserializers()
    res["groups"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateGroupFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Groupable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Groupable)
                }
            }
            m.SetGroups(res)
        }
        return nil
    }
    return res
}
// GetGroups gets the groups property value. The groups property
// returns a []Groupable when successful
func (m *GroupBasedSubjectSet) GetGroups()([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Groupable) {
    val, err := m.GetBackingStore().Get("groups")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Groupable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *GroupBasedSubjectSet) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.SubjectSet.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetGroups() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetGroups()))
        for i, v := range m.GetGroups() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("groups", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetGroups sets the groups property value. The groups property
func (m *GroupBasedSubjectSet) SetGroups(value []iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Groupable)() {
    err := m.GetBackingStore().Set("groups", value)
    if err != nil {
        panic(err)
    }
}
type GroupBasedSubjectSetable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SubjectSetable
    GetGroups()([]iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Groupable)
    SetGroups(value []iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Groupable)()
}
