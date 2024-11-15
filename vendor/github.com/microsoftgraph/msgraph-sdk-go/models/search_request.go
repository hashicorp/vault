package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type SearchRequest struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewSearchRequest instantiates a new SearchRequest and sets the default values.
func NewSearchRequest()(*SearchRequest) {
    m := &SearchRequest{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateSearchRequestFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSearchRequestFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSearchRequest(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *SearchRequest) GetAdditionalData()(map[string]any) {
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
// GetAggregationFilters gets the aggregationFilters property value. Contains one or more filters to obtain search results aggregated and filtered to a specific value of a field. Optional.Build this filter based on a prior search that aggregates by the same field. From the response of the prior search, identify the searchBucket that filters results to the specific value of the field, use the string in its aggregationFilterToken property, and build an aggregation filter string in the format '{field}:/'{aggregationFilterToken}/''. If multiple values for the same field need to be provided, use the strings in its aggregationFilterToken property and build an aggregation filter string in the format '{field}:or(/'{aggregationFilterToken1}/',/'{aggregationFilterToken2}/')'. For example, searching and aggregating drive items by file type returns a searchBucket for the file type docx in the response. You can conveniently use the aggregationFilterToken returned for this searchBucket in a subsequent search query and filter matches down to drive items of the docx file type. Example 1 and example 2 show the actual requests and responses.
// returns a []string when successful
func (m *SearchRequest) GetAggregationFilters()([]string) {
    val, err := m.GetBackingStore().Get("aggregationFilters")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetAggregations gets the aggregations property value. Specifies aggregations (also known as refiners) to be returned alongside search results. Optional.
// returns a []AggregationOptionable when successful
func (m *SearchRequest) GetAggregations()([]AggregationOptionable) {
    val, err := m.GetBackingStore().Get("aggregations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AggregationOptionable)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *SearchRequest) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCollapseProperties gets the collapseProperties property value. Contains the ordered collection of fields and limit to collapse results. Optional.
// returns a []CollapsePropertyable when successful
func (m *SearchRequest) GetCollapseProperties()([]CollapsePropertyable) {
    val, err := m.GetBackingStore().Get("collapseProperties")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CollapsePropertyable)
    }
    return nil
}
// GetContentSources gets the contentSources property value. Contains the connection to be targeted.
// returns a []string when successful
func (m *SearchRequest) GetContentSources()([]string) {
    val, err := m.GetBackingStore().Get("contentSources")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetEnableTopResults gets the enableTopResults property value. This triggers hybrid sort for messages : the first 3 messages are the most relevant. This property is only applicable to entityType=message. Optional.
// returns a *bool when successful
func (m *SearchRequest) GetEnableTopResults()(*bool) {
    val, err := m.GetBackingStore().Get("enableTopResults")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEntityTypes gets the entityTypes property value. One or more types of resources expected in the response. Possible values are: event, message, driveItem, externalItem, site, list, listItem, drive, chatMessage, person, acronym, bookmark.  Note that you must use the Prefer: include-unknown-enum-members request header to get the following value(s) in this evolvable enum: chatMessage, person, acronym, bookmark. See known limitations for those combinations of two or more entity types that are supported in the same search request. Required.
// returns a []EntityType when successful
func (m *SearchRequest) GetEntityTypes()([]EntityType) {
    val, err := m.GetBackingStore().Get("entityTypes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EntityType)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SearchRequest) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["aggregationFilters"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetAggregationFilters(res)
        }
        return nil
    }
    res["aggregations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAggregationOptionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AggregationOptionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AggregationOptionable)
                }
            }
            m.SetAggregations(res)
        }
        return nil
    }
    res["collapseProperties"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCollapsePropertyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CollapsePropertyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CollapsePropertyable)
                }
            }
            m.SetCollapseProperties(res)
        }
        return nil
    }
    res["contentSources"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetContentSources(res)
        }
        return nil
    }
    res["enableTopResults"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnableTopResults(val)
        }
        return nil
    }
    res["entityTypes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseEntityType)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EntityType, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*EntityType))
                }
            }
            m.SetEntityTypes(res)
        }
        return nil
    }
    res["fields"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetFields(res)
        }
        return nil
    }
    res["from"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFrom(val)
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
    res["query"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSearchQueryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQuery(val.(SearchQueryable))
        }
        return nil
    }
    res["queryAlterationOptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSearchAlterationOptionsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQueryAlterationOptions(val.(SearchAlterationOptionsable))
        }
        return nil
    }
    res["region"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRegion(val)
        }
        return nil
    }
    res["resultTemplateOptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateResultTemplateOptionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResultTemplateOptions(val.(ResultTemplateOptionable))
        }
        return nil
    }
    res["sharePointOneDriveOptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSharePointOneDriveOptionsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSharePointOneDriveOptions(val.(SharePointOneDriveOptionsable))
        }
        return nil
    }
    res["size"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSize(val)
        }
        return nil
    }
    res["sortProperties"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSortPropertyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SortPropertyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SortPropertyable)
                }
            }
            m.SetSortProperties(res)
        }
        return nil
    }
    return res
}
// GetFields gets the fields property value. Contains the fields to be returned for each resource object specified in entityTypes, allowing customization of the fields returned by default; otherwise, including additional fields such as custom managed properties from SharePoint and OneDrive, or custom fields in externalItem from the content that Microsoft Graph connectors bring in. The fields property can use the semantic labels applied to properties. For example, if a property is labeled as title, you can retrieve it using the following syntax: label_title. Optional.
// returns a []string when successful
func (m *SearchRequest) GetFields()([]string) {
    val, err := m.GetBackingStore().Get("fields")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetFrom gets the from property value. Specifies the offset for the search results. Offset 0 returns the very first result. Optional.
// returns a *int32 when successful
func (m *SearchRequest) GetFrom()(*int32) {
    val, err := m.GetBackingStore().Get("from")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *SearchRequest) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetQuery gets the query property value. The query property
// returns a SearchQueryable when successful
func (m *SearchRequest) GetQuery()(SearchQueryable) {
    val, err := m.GetBackingStore().Get("query")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SearchQueryable)
    }
    return nil
}
// GetQueryAlterationOptions gets the queryAlterationOptions property value. Query alteration options formatted in a JSON blob that contains two optional flags related to spelling correction. Optional.
// returns a SearchAlterationOptionsable when successful
func (m *SearchRequest) GetQueryAlterationOptions()(SearchAlterationOptionsable) {
    val, err := m.GetBackingStore().Get("queryAlterationOptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SearchAlterationOptionsable)
    }
    return nil
}
// GetRegion gets the region property value. The geographic location for the search. Required for searches that use application permissions. For details, see Get the region value.
// returns a *string when successful
func (m *SearchRequest) GetRegion()(*string) {
    val, err := m.GetBackingStore().Get("region")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetResultTemplateOptions gets the resultTemplateOptions property value. Provides the search result template options to render search results from connectors.
// returns a ResultTemplateOptionable when successful
func (m *SearchRequest) GetResultTemplateOptions()(ResultTemplateOptionable) {
    val, err := m.GetBackingStore().Get("resultTemplateOptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ResultTemplateOptionable)
    }
    return nil
}
// GetSharePointOneDriveOptions gets the sharePointOneDriveOptions property value. Indicates the kind of contents to be searched when a search is performed using application permissions. Optional.
// returns a SharePointOneDriveOptionsable when successful
func (m *SearchRequest) GetSharePointOneDriveOptions()(SharePointOneDriveOptionsable) {
    val, err := m.GetBackingStore().Get("sharePointOneDriveOptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SharePointOneDriveOptionsable)
    }
    return nil
}
// GetSize gets the size property value. The size of the page to be retrieved. The maximum value is 500. Optional.
// returns a *int32 when successful
func (m *SearchRequest) GetSize()(*int32) {
    val, err := m.GetBackingStore().Get("size")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSortProperties gets the sortProperties property value. Contains the ordered collection of fields and direction to sort results. There can be at most 5 sort properties in the collection. Optional.
// returns a []SortPropertyable when successful
func (m *SearchRequest) GetSortProperties()([]SortPropertyable) {
    val, err := m.GetBackingStore().Get("sortProperties")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SortPropertyable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SearchRequest) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAggregationFilters() != nil {
        err := writer.WriteCollectionOfStringValues("aggregationFilters", m.GetAggregationFilters())
        if err != nil {
            return err
        }
    }
    if m.GetAggregations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAggregations()))
        for i, v := range m.GetAggregations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("aggregations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetCollapseProperties() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCollapseProperties()))
        for i, v := range m.GetCollapseProperties() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("collapseProperties", cast)
        if err != nil {
            return err
        }
    }
    if m.GetContentSources() != nil {
        err := writer.WriteCollectionOfStringValues("contentSources", m.GetContentSources())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("enableTopResults", m.GetEnableTopResults())
        if err != nil {
            return err
        }
    }
    if m.GetEntityTypes() != nil {
        err := writer.WriteCollectionOfStringValues("entityTypes", SerializeEntityType(m.GetEntityTypes()))
        if err != nil {
            return err
        }
    }
    if m.GetFields() != nil {
        err := writer.WriteCollectionOfStringValues("fields", m.GetFields())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("from", m.GetFrom())
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
        err := writer.WriteObjectValue("query", m.GetQuery())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("queryAlterationOptions", m.GetQueryAlterationOptions())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("region", m.GetRegion())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("resultTemplateOptions", m.GetResultTemplateOptions())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("sharePointOneDriveOptions", m.GetSharePointOneDriveOptions())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("size", m.GetSize())
        if err != nil {
            return err
        }
    }
    if m.GetSortProperties() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSortProperties()))
        for i, v := range m.GetSortProperties() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("sortProperties", cast)
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
func (m *SearchRequest) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAggregationFilters sets the aggregationFilters property value. Contains one or more filters to obtain search results aggregated and filtered to a specific value of a field. Optional.Build this filter based on a prior search that aggregates by the same field. From the response of the prior search, identify the searchBucket that filters results to the specific value of the field, use the string in its aggregationFilterToken property, and build an aggregation filter string in the format '{field}:/'{aggregationFilterToken}/''. If multiple values for the same field need to be provided, use the strings in its aggregationFilterToken property and build an aggregation filter string in the format '{field}:or(/'{aggregationFilterToken1}/',/'{aggregationFilterToken2}/')'. For example, searching and aggregating drive items by file type returns a searchBucket for the file type docx in the response. You can conveniently use the aggregationFilterToken returned for this searchBucket in a subsequent search query and filter matches down to drive items of the docx file type. Example 1 and example 2 show the actual requests and responses.
func (m *SearchRequest) SetAggregationFilters(value []string)() {
    err := m.GetBackingStore().Set("aggregationFilters", value)
    if err != nil {
        panic(err)
    }
}
// SetAggregations sets the aggregations property value. Specifies aggregations (also known as refiners) to be returned alongside search results. Optional.
func (m *SearchRequest) SetAggregations(value []AggregationOptionable)() {
    err := m.GetBackingStore().Set("aggregations", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *SearchRequest) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCollapseProperties sets the collapseProperties property value. Contains the ordered collection of fields and limit to collapse results. Optional.
func (m *SearchRequest) SetCollapseProperties(value []CollapsePropertyable)() {
    err := m.GetBackingStore().Set("collapseProperties", value)
    if err != nil {
        panic(err)
    }
}
// SetContentSources sets the contentSources property value. Contains the connection to be targeted.
func (m *SearchRequest) SetContentSources(value []string)() {
    err := m.GetBackingStore().Set("contentSources", value)
    if err != nil {
        panic(err)
    }
}
// SetEnableTopResults sets the enableTopResults property value. This triggers hybrid sort for messages : the first 3 messages are the most relevant. This property is only applicable to entityType=message. Optional.
func (m *SearchRequest) SetEnableTopResults(value *bool)() {
    err := m.GetBackingStore().Set("enableTopResults", value)
    if err != nil {
        panic(err)
    }
}
// SetEntityTypes sets the entityTypes property value. One or more types of resources expected in the response. Possible values are: event, message, driveItem, externalItem, site, list, listItem, drive, chatMessage, person, acronym, bookmark.  Note that you must use the Prefer: include-unknown-enum-members request header to get the following value(s) in this evolvable enum: chatMessage, person, acronym, bookmark. See known limitations for those combinations of two or more entity types that are supported in the same search request. Required.
func (m *SearchRequest) SetEntityTypes(value []EntityType)() {
    err := m.GetBackingStore().Set("entityTypes", value)
    if err != nil {
        panic(err)
    }
}
// SetFields sets the fields property value. Contains the fields to be returned for each resource object specified in entityTypes, allowing customization of the fields returned by default; otherwise, including additional fields such as custom managed properties from SharePoint and OneDrive, or custom fields in externalItem from the content that Microsoft Graph connectors bring in. The fields property can use the semantic labels applied to properties. For example, if a property is labeled as title, you can retrieve it using the following syntax: label_title. Optional.
func (m *SearchRequest) SetFields(value []string)() {
    err := m.GetBackingStore().Set("fields", value)
    if err != nil {
        panic(err)
    }
}
// SetFrom sets the from property value. Specifies the offset for the search results. Offset 0 returns the very first result. Optional.
func (m *SearchRequest) SetFrom(value *int32)() {
    err := m.GetBackingStore().Set("from", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SearchRequest) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetQuery sets the query property value. The query property
func (m *SearchRequest) SetQuery(value SearchQueryable)() {
    err := m.GetBackingStore().Set("query", value)
    if err != nil {
        panic(err)
    }
}
// SetQueryAlterationOptions sets the queryAlterationOptions property value. Query alteration options formatted in a JSON blob that contains two optional flags related to spelling correction. Optional.
func (m *SearchRequest) SetQueryAlterationOptions(value SearchAlterationOptionsable)() {
    err := m.GetBackingStore().Set("queryAlterationOptions", value)
    if err != nil {
        panic(err)
    }
}
// SetRegion sets the region property value. The geographic location for the search. Required for searches that use application permissions. For details, see Get the region value.
func (m *SearchRequest) SetRegion(value *string)() {
    err := m.GetBackingStore().Set("region", value)
    if err != nil {
        panic(err)
    }
}
// SetResultTemplateOptions sets the resultTemplateOptions property value. Provides the search result template options to render search results from connectors.
func (m *SearchRequest) SetResultTemplateOptions(value ResultTemplateOptionable)() {
    err := m.GetBackingStore().Set("resultTemplateOptions", value)
    if err != nil {
        panic(err)
    }
}
// SetSharePointOneDriveOptions sets the sharePointOneDriveOptions property value. Indicates the kind of contents to be searched when a search is performed using application permissions. Optional.
func (m *SearchRequest) SetSharePointOneDriveOptions(value SharePointOneDriveOptionsable)() {
    err := m.GetBackingStore().Set("sharePointOneDriveOptions", value)
    if err != nil {
        panic(err)
    }
}
// SetSize sets the size property value. The size of the page to be retrieved. The maximum value is 500. Optional.
func (m *SearchRequest) SetSize(value *int32)() {
    err := m.GetBackingStore().Set("size", value)
    if err != nil {
        panic(err)
    }
}
// SetSortProperties sets the sortProperties property value. Contains the ordered collection of fields and direction to sort results. There can be at most 5 sort properties in the collection. Optional.
func (m *SearchRequest) SetSortProperties(value []SortPropertyable)() {
    err := m.GetBackingStore().Set("sortProperties", value)
    if err != nil {
        panic(err)
    }
}
type SearchRequestable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAggregationFilters()([]string)
    GetAggregations()([]AggregationOptionable)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCollapseProperties()([]CollapsePropertyable)
    GetContentSources()([]string)
    GetEnableTopResults()(*bool)
    GetEntityTypes()([]EntityType)
    GetFields()([]string)
    GetFrom()(*int32)
    GetOdataType()(*string)
    GetQuery()(SearchQueryable)
    GetQueryAlterationOptions()(SearchAlterationOptionsable)
    GetRegion()(*string)
    GetResultTemplateOptions()(ResultTemplateOptionable)
    GetSharePointOneDriveOptions()(SharePointOneDriveOptionsable)
    GetSize()(*int32)
    GetSortProperties()([]SortPropertyable)
    SetAggregationFilters(value []string)()
    SetAggregations(value []AggregationOptionable)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCollapseProperties(value []CollapsePropertyable)()
    SetContentSources(value []string)()
    SetEnableTopResults(value *bool)()
    SetEntityTypes(value []EntityType)()
    SetFields(value []string)()
    SetFrom(value *int32)()
    SetOdataType(value *string)()
    SetQuery(value SearchQueryable)()
    SetQueryAlterationOptions(value SearchAlterationOptionsable)()
    SetRegion(value *string)()
    SetResultTemplateOptions(value ResultTemplateOptionable)()
    SetSharePointOneDriveOptions(value SharePointOneDriveOptionsable)()
    SetSize(value *int32)()
    SetSortProperties(value []SortPropertyable)()
}
