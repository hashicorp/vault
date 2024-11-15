package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CitationTemplate struct {
    FilePlanDescriptorTemplate
}
// NewCitationTemplate instantiates a new CitationTemplate and sets the default values.
func NewCitationTemplate()(*CitationTemplate) {
    m := &CitationTemplate{
        FilePlanDescriptorTemplate: *NewFilePlanDescriptorTemplate(),
    }
    return m
}
// CreateCitationTemplateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCitationTemplateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCitationTemplate(), nil
}
// GetCitationJurisdiction gets the citationJurisdiction property value. Represents the jurisdiction or agency that published the citation.
// returns a *string when successful
func (m *CitationTemplate) GetCitationJurisdiction()(*string) {
    val, err := m.GetBackingStore().Get("citationJurisdiction")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCitationUrl gets the citationUrl property value. Represents the URL to the published citation.
// returns a *string when successful
func (m *CitationTemplate) GetCitationUrl()(*string) {
    val, err := m.GetBackingStore().Get("citationUrl")
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
func (m *CitationTemplate) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.FilePlanDescriptorTemplate.GetFieldDeserializers()
    res["citationJurisdiction"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCitationJurisdiction(val)
        }
        return nil
    }
    res["citationUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCitationUrl(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *CitationTemplate) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.FilePlanDescriptorTemplate.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("citationJurisdiction", m.GetCitationJurisdiction())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("citationUrl", m.GetCitationUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCitationJurisdiction sets the citationJurisdiction property value. Represents the jurisdiction or agency that published the citation.
func (m *CitationTemplate) SetCitationJurisdiction(value *string)() {
    err := m.GetBackingStore().Set("citationJurisdiction", value)
    if err != nil {
        panic(err)
    }
}
// SetCitationUrl sets the citationUrl property value. Represents the URL to the published citation.
func (m *CitationTemplate) SetCitationUrl(value *string)() {
    err := m.GetBackingStore().Set("citationUrl", value)
    if err != nil {
        panic(err)
    }
}
type CitationTemplateable interface {
    FilePlanDescriptorTemplateable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCitationJurisdiction()(*string)
    GetCitationUrl()(*string)
    SetCitationJurisdiction(value *string)()
    SetCitationUrl(value *string)()
}
