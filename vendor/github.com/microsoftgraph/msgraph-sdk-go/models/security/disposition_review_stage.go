package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type DispositionReviewStage struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
}
// NewDispositionReviewStage instantiates a new DispositionReviewStage and sets the default values.
func NewDispositionReviewStage()(*DispositionReviewStage) {
    m := &DispositionReviewStage{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateDispositionReviewStageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDispositionReviewStageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDispositionReviewStage(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DispositionReviewStage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
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
    res["reviewersEmailAddresses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetReviewersEmailAddresses(res)
        }
        return nil
    }
    res["stageNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStageNumber(val)
        }
        return nil
    }
    return res
}
// GetName gets the name property value. Name representing each stage within a collection.
// returns a *string when successful
func (m *DispositionReviewStage) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetReviewersEmailAddresses gets the reviewersEmailAddresses property value. A collection of reviewers at each stage.
// returns a []string when successful
func (m *DispositionReviewStage) GetReviewersEmailAddresses()([]string) {
    val, err := m.GetBackingStore().Get("reviewersEmailAddresses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetStageNumber gets the stageNumber property value. The unique sequence number for each stage of the disposition review.
// returns a *string when successful
func (m *DispositionReviewStage) GetStageNumber()(*string) {
    val, err := m.GetBackingStore().Get("stageNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DispositionReviewStage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("name", m.GetName())
        if err != nil {
            return err
        }
    }
    if m.GetReviewersEmailAddresses() != nil {
        err = writer.WriteCollectionOfStringValues("reviewersEmailAddresses", m.GetReviewersEmailAddresses())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("stageNumber", m.GetStageNumber())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetName sets the name property value. Name representing each stage within a collection.
func (m *DispositionReviewStage) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetReviewersEmailAddresses sets the reviewersEmailAddresses property value. A collection of reviewers at each stage.
func (m *DispositionReviewStage) SetReviewersEmailAddresses(value []string)() {
    err := m.GetBackingStore().Set("reviewersEmailAddresses", value)
    if err != nil {
        panic(err)
    }
}
// SetStageNumber sets the stageNumber property value. The unique sequence number for each stage of the disposition review.
func (m *DispositionReviewStage) SetStageNumber(value *string)() {
    err := m.GetBackingStore().Set("stageNumber", value)
    if err != nil {
        panic(err)
    }
}
type DispositionReviewStageable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetName()(*string)
    GetReviewersEmailAddresses()([]string)
    GetStageNumber()(*string)
    SetName(value *string)()
    SetReviewersEmailAddresses(value []string)()
    SetStageNumber(value *string)()
}
