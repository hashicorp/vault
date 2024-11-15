package identitygovernance

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type RuleBasedSubjectSet struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SubjectSet
}
// NewRuleBasedSubjectSet instantiates a new RuleBasedSubjectSet and sets the default values.
func NewRuleBasedSubjectSet()(*RuleBasedSubjectSet) {
    m := &RuleBasedSubjectSet{
        SubjectSet: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewSubjectSet(),
    }
    odataTypeValue := "#microsoft.graph.identityGovernance.ruleBasedSubjectSet"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateRuleBasedSubjectSetFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRuleBasedSubjectSetFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRuleBasedSubjectSet(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *RuleBasedSubjectSet) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.SubjectSet.GetFieldDeserializers()
    res["rule"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRule(val)
        }
        return nil
    }
    return res
}
// GetRule gets the rule property value. The rule for the subject set. Lifecycle Workflows supports a rich set of user properties for configuring the rules using $filter query expressions. For more information, see supported user and query parameters.
// returns a *string when successful
func (m *RuleBasedSubjectSet) GetRule()(*string) {
    val, err := m.GetBackingStore().Get("rule")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *RuleBasedSubjectSet) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.SubjectSet.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("rule", m.GetRule())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetRule sets the rule property value. The rule for the subject set. Lifecycle Workflows supports a rich set of user properties for configuring the rules using $filter query expressions. For more information, see supported user and query parameters.
func (m *RuleBasedSubjectSet) SetRule(value *string)() {
    err := m.GetBackingStore().Set("rule", value)
    if err != nil {
        panic(err)
    }
}
type RuleBasedSubjectSetable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SubjectSetable
    GetRule()(*string)
    SetRule(value *string)()
}
