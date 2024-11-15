package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type EdiscoveryCaseSettings struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
}
// NewEdiscoveryCaseSettings instantiates a new EdiscoveryCaseSettings and sets the default values.
func NewEdiscoveryCaseSettings()(*EdiscoveryCaseSettings) {
    m := &EdiscoveryCaseSettings{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateEdiscoveryCaseSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEdiscoveryCaseSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEdiscoveryCaseSettings(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EdiscoveryCaseSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["ocr"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOcrSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOcr(val.(OcrSettingsable))
        }
        return nil
    }
    res["redundancyDetection"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateRedundancyDetectionSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRedundancyDetection(val.(RedundancyDetectionSettingsable))
        }
        return nil
    }
    res["topicModeling"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTopicModelingSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTopicModeling(val.(TopicModelingSettingsable))
        }
        return nil
    }
    return res
}
// GetOcr gets the ocr property value. The OCR (Optical Character Recognition) settings for the case.
// returns a OcrSettingsable when successful
func (m *EdiscoveryCaseSettings) GetOcr()(OcrSettingsable) {
    val, err := m.GetBackingStore().Get("ocr")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(OcrSettingsable)
    }
    return nil
}
// GetRedundancyDetection gets the redundancyDetection property value. The redundancy (near duplicate and email threading) detection settings for the case.
// returns a RedundancyDetectionSettingsable when successful
func (m *EdiscoveryCaseSettings) GetRedundancyDetection()(RedundancyDetectionSettingsable) {
    val, err := m.GetBackingStore().Get("redundancyDetection")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(RedundancyDetectionSettingsable)
    }
    return nil
}
// GetTopicModeling gets the topicModeling property value. The Topic Modeling (Themes) settings for the case.
// returns a TopicModelingSettingsable when successful
func (m *EdiscoveryCaseSettings) GetTopicModeling()(TopicModelingSettingsable) {
    val, err := m.GetBackingStore().Get("topicModeling")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TopicModelingSettingsable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EdiscoveryCaseSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("ocr", m.GetOcr())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("redundancyDetection", m.GetRedundancyDetection())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("topicModeling", m.GetTopicModeling())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetOcr sets the ocr property value. The OCR (Optical Character Recognition) settings for the case.
func (m *EdiscoveryCaseSettings) SetOcr(value OcrSettingsable)() {
    err := m.GetBackingStore().Set("ocr", value)
    if err != nil {
        panic(err)
    }
}
// SetRedundancyDetection sets the redundancyDetection property value. The redundancy (near duplicate and email threading) detection settings for the case.
func (m *EdiscoveryCaseSettings) SetRedundancyDetection(value RedundancyDetectionSettingsable)() {
    err := m.GetBackingStore().Set("redundancyDetection", value)
    if err != nil {
        panic(err)
    }
}
// SetTopicModeling sets the topicModeling property value. The Topic Modeling (Themes) settings for the case.
func (m *EdiscoveryCaseSettings) SetTopicModeling(value TopicModelingSettingsable)() {
    err := m.GetBackingStore().Set("topicModeling", value)
    if err != nil {
        panic(err)
    }
}
type EdiscoveryCaseSettingsable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetOcr()(OcrSettingsable)
    GetRedundancyDetection()(RedundancyDetectionSettingsable)
    GetTopicModeling()(TopicModelingSettingsable)
    SetOcr(value OcrSettingsable)()
    SetRedundancyDetection(value RedundancyDetectionSettingsable)()
    SetTopicModeling(value TopicModelingSettingsable)()
}
