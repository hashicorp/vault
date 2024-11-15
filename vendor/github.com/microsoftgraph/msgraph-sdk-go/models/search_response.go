package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type SearchResponse struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewSearchResponse instantiates a new SearchResponse and sets the default values.
func NewSearchResponse()(*SearchResponse) {
    m := &SearchResponse{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateSearchResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSearchResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSearchResponse(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *SearchResponse) GetAdditionalData()(map[string]any) {
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
func (m *SearchResponse) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SearchResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["hitsContainers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSearchHitsContainerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SearchHitsContainerable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SearchHitsContainerable)
                }
            }
            m.SetHitsContainers(res)
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
    res["queryAlterationResponse"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAlterationResponseFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQueryAlterationResponse(val.(AlterationResponseable))
        }
        return nil
    }
    res["resultTemplates"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateResultTemplateDictionaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResultTemplates(val.(ResultTemplateDictionaryable))
        }
        return nil
    }
    res["searchTerms"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSearchTerms(res)
        }
        return nil
    }
    return res
}
// GetHitsContainers gets the hitsContainers property value. A collection of search results.
// returns a []SearchHitsContainerable when successful
func (m *SearchResponse) GetHitsContainers()([]SearchHitsContainerable) {
    val, err := m.GetBackingStore().Get("hitsContainers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SearchHitsContainerable)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *SearchResponse) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetQueryAlterationResponse gets the queryAlterationResponse property value. Provides information related to spelling corrections in the alteration response.
// returns a AlterationResponseable when successful
func (m *SearchResponse) GetQueryAlterationResponse()(AlterationResponseable) {
    val, err := m.GetBackingStore().Get("queryAlterationResponse")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AlterationResponseable)
    }
    return nil
}
// GetResultTemplates gets the resultTemplates property value. A dictionary of resultTemplateIds and associated values, which include the name and JSON schema of the result templates.
// returns a ResultTemplateDictionaryable when successful
func (m *SearchResponse) GetResultTemplates()(ResultTemplateDictionaryable) {
    val, err := m.GetBackingStore().Get("resultTemplates")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ResultTemplateDictionaryable)
    }
    return nil
}
// GetSearchTerms gets the searchTerms property value. Contains the search terms sent in the initial search query.
// returns a []string when successful
func (m *SearchResponse) GetSearchTerms()([]string) {
    val, err := m.GetBackingStore().Get("searchTerms")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SearchResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetHitsContainers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetHitsContainers()))
        for i, v := range m.GetHitsContainers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("hitsContainers", cast)
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
    {
        err := writer.WriteObjectValue("queryAlterationResponse", m.GetQueryAlterationResponse())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("resultTemplates", m.GetResultTemplates())
        if err != nil {
            return err
        }
    }
    if m.GetSearchTerms() != nil {
        err := writer.WriteCollectionOfStringValues("searchTerms", m.GetSearchTerms())
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
func (m *SearchResponse) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *SearchResponse) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetHitsContainers sets the hitsContainers property value. A collection of search results.
func (m *SearchResponse) SetHitsContainers(value []SearchHitsContainerable)() {
    err := m.GetBackingStore().Set("hitsContainers", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SearchResponse) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetQueryAlterationResponse sets the queryAlterationResponse property value. Provides information related to spelling corrections in the alteration response.
func (m *SearchResponse) SetQueryAlterationResponse(value AlterationResponseable)() {
    err := m.GetBackingStore().Set("queryAlterationResponse", value)
    if err != nil {
        panic(err)
    }
}
// SetResultTemplates sets the resultTemplates property value. A dictionary of resultTemplateIds and associated values, which include the name and JSON schema of the result templates.
func (m *SearchResponse) SetResultTemplates(value ResultTemplateDictionaryable)() {
    err := m.GetBackingStore().Set("resultTemplates", value)
    if err != nil {
        panic(err)
    }
}
// SetSearchTerms sets the searchTerms property value. Contains the search terms sent in the initial search query.
func (m *SearchResponse) SetSearchTerms(value []string)() {
    err := m.GetBackingStore().Set("searchTerms", value)
    if err != nil {
        panic(err)
    }
}
type SearchResponseable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetHitsContainers()([]SearchHitsContainerable)
    GetOdataType()(*string)
    GetQueryAlterationResponse()(AlterationResponseable)
    GetResultTemplates()(ResultTemplateDictionaryable)
    GetSearchTerms()([]string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetHitsContainers(value []SearchHitsContainerable)()
    SetOdataType(value *string)()
    SetQueryAlterationResponse(value AlterationResponseable)()
    SetResultTemplates(value ResultTemplateDictionaryable)()
    SetSearchTerms(value []string)()
}
