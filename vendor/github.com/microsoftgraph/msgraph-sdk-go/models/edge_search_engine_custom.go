package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// EdgeSearchEngineCustom allows IT admins to set a custom default search engine for MDM-Controlled devices.
type EdgeSearchEngineCustom struct {
    EdgeSearchEngineBase
}
// NewEdgeSearchEngineCustom instantiates a new EdgeSearchEngineCustom and sets the default values.
func NewEdgeSearchEngineCustom()(*EdgeSearchEngineCustom) {
    m := &EdgeSearchEngineCustom{
        EdgeSearchEngineBase: *NewEdgeSearchEngineBase(),
    }
    odataTypeValue := "#microsoft.graph.edgeSearchEngineCustom"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEdgeSearchEngineCustomFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEdgeSearchEngineCustomFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEdgeSearchEngineCustom(), nil
}
// GetEdgeSearchEngineOpenSearchXmlUrl gets the edgeSearchEngineOpenSearchXmlUrl property value. Points to a https link containing the OpenSearch xml file that contains, at minimum, the short name and the URL to the search Engine.
// returns a *string when successful
func (m *EdgeSearchEngineCustom) GetEdgeSearchEngineOpenSearchXmlUrl()(*string) {
    val, err := m.GetBackingStore().Get("edgeSearchEngineOpenSearchXmlUrl")
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
func (m *EdgeSearchEngineCustom) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EdgeSearchEngineBase.GetFieldDeserializers()
    res["edgeSearchEngineOpenSearchXmlUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeSearchEngineOpenSearchXmlUrl(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *EdgeSearchEngineCustom) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EdgeSearchEngineBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("edgeSearchEngineOpenSearchXmlUrl", m.GetEdgeSearchEngineOpenSearchXmlUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetEdgeSearchEngineOpenSearchXmlUrl sets the edgeSearchEngineOpenSearchXmlUrl property value. Points to a https link containing the OpenSearch xml file that contains, at minimum, the short name and the URL to the search Engine.
func (m *EdgeSearchEngineCustom) SetEdgeSearchEngineOpenSearchXmlUrl(value *string)() {
    err := m.GetBackingStore().Set("edgeSearchEngineOpenSearchXmlUrl", value)
    if err != nil {
        panic(err)
    }
}
type EdgeSearchEngineCustomable interface {
    EdgeSearchEngineBaseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetEdgeSearchEngineOpenSearchXmlUrl()(*string)
    SetEdgeSearchEngineOpenSearchXmlUrl(value *string)()
}
