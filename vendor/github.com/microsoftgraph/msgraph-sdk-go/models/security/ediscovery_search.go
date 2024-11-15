package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EdiscoverySearch struct {
    Search
}
// NewEdiscoverySearch instantiates a new EdiscoverySearch and sets the default values.
func NewEdiscoverySearch()(*EdiscoverySearch) {
    m := &EdiscoverySearch{
        Search: *NewSearch(),
    }
    odataTypeValue := "#microsoft.graph.security.ediscoverySearch"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEdiscoverySearchFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEdiscoverySearchFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEdiscoverySearch(), nil
}
// GetAdditionalSources gets the additionalSources property value. Adds an additional source to the eDiscovery search.
// returns a []DataSourceable when successful
func (m *EdiscoverySearch) GetAdditionalSources()([]DataSourceable) {
    val, err := m.GetBackingStore().Get("additionalSources")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DataSourceable)
    }
    return nil
}
// GetAddToReviewSetOperation gets the addToReviewSetOperation property value. Adds the results of the eDiscovery search to the specified reviewSet.
// returns a EdiscoveryAddToReviewSetOperationable when successful
func (m *EdiscoverySearch) GetAddToReviewSetOperation()(EdiscoveryAddToReviewSetOperationable) {
    val, err := m.GetBackingStore().Get("addToReviewSetOperation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EdiscoveryAddToReviewSetOperationable)
    }
    return nil
}
// GetCustodianSources gets the custodianSources property value. Custodian sources that are included in the eDiscovery search.
// returns a []DataSourceable when successful
func (m *EdiscoverySearch) GetCustodianSources()([]DataSourceable) {
    val, err := m.GetBackingStore().Get("custodianSources")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DataSourceable)
    }
    return nil
}
// GetDataSourceScopes gets the dataSourceScopes property value. When specified, the collection will span across a service for an entire workload. Possible values are: none, allTenantMailboxes, allTenantSites, allCaseCustodians, allCaseNoncustodialDataSources.
// returns a *DataSourceScopes when successful
func (m *EdiscoverySearch) GetDataSourceScopes()(*DataSourceScopes) {
    val, err := m.GetBackingStore().Get("dataSourceScopes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DataSourceScopes)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EdiscoverySearch) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Search.GetFieldDeserializers()
    res["additionalSources"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDataSourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DataSourceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DataSourceable)
                }
            }
            m.SetAdditionalSources(res)
        }
        return nil
    }
    res["addToReviewSetOperation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEdiscoveryAddToReviewSetOperationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAddToReviewSetOperation(val.(EdiscoveryAddToReviewSetOperationable))
        }
        return nil
    }
    res["custodianSources"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDataSourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DataSourceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DataSourceable)
                }
            }
            m.SetCustodianSources(res)
        }
        return nil
    }
    res["dataSourceScopes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDataSourceScopes)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDataSourceScopes(val.(*DataSourceScopes))
        }
        return nil
    }
    res["lastEstimateStatisticsOperation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEdiscoveryEstimateOperationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastEstimateStatisticsOperation(val.(EdiscoveryEstimateOperationable))
        }
        return nil
    }
    res["noncustodialSources"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEdiscoveryNoncustodialDataSourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EdiscoveryNoncustodialDataSourceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EdiscoveryNoncustodialDataSourceable)
                }
            }
            m.SetNoncustodialSources(res)
        }
        return nil
    }
    return res
}
// GetLastEstimateStatisticsOperation gets the lastEstimateStatisticsOperation property value. The last estimate operation associated with the eDiscovery search.
// returns a EdiscoveryEstimateOperationable when successful
func (m *EdiscoverySearch) GetLastEstimateStatisticsOperation()(EdiscoveryEstimateOperationable) {
    val, err := m.GetBackingStore().Get("lastEstimateStatisticsOperation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EdiscoveryEstimateOperationable)
    }
    return nil
}
// GetNoncustodialSources gets the noncustodialSources property value. noncustodialDataSource sources that are included in the eDiscovery search
// returns a []EdiscoveryNoncustodialDataSourceable when successful
func (m *EdiscoverySearch) GetNoncustodialSources()([]EdiscoveryNoncustodialDataSourceable) {
    val, err := m.GetBackingStore().Get("noncustodialSources")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EdiscoveryNoncustodialDataSourceable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EdiscoverySearch) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Search.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAdditionalSources() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAdditionalSources()))
        for i, v := range m.GetAdditionalSources() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("additionalSources", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("addToReviewSetOperation", m.GetAddToReviewSetOperation())
        if err != nil {
            return err
        }
    }
    if m.GetCustodianSources() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCustodianSources()))
        for i, v := range m.GetCustodianSources() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("custodianSources", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDataSourceScopes() != nil {
        cast := (*m.GetDataSourceScopes()).String()
        err = writer.WriteStringValue("dataSourceScopes", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("lastEstimateStatisticsOperation", m.GetLastEstimateStatisticsOperation())
        if err != nil {
            return err
        }
    }
    if m.GetNoncustodialSources() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetNoncustodialSources()))
        for i, v := range m.GetNoncustodialSources() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("noncustodialSources", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalSources sets the additionalSources property value. Adds an additional source to the eDiscovery search.
func (m *EdiscoverySearch) SetAdditionalSources(value []DataSourceable)() {
    err := m.GetBackingStore().Set("additionalSources", value)
    if err != nil {
        panic(err)
    }
}
// SetAddToReviewSetOperation sets the addToReviewSetOperation property value. Adds the results of the eDiscovery search to the specified reviewSet.
func (m *EdiscoverySearch) SetAddToReviewSetOperation(value EdiscoveryAddToReviewSetOperationable)() {
    err := m.GetBackingStore().Set("addToReviewSetOperation", value)
    if err != nil {
        panic(err)
    }
}
// SetCustodianSources sets the custodianSources property value. Custodian sources that are included in the eDiscovery search.
func (m *EdiscoverySearch) SetCustodianSources(value []DataSourceable)() {
    err := m.GetBackingStore().Set("custodianSources", value)
    if err != nil {
        panic(err)
    }
}
// SetDataSourceScopes sets the dataSourceScopes property value. When specified, the collection will span across a service for an entire workload. Possible values are: none, allTenantMailboxes, allTenantSites, allCaseCustodians, allCaseNoncustodialDataSources.
func (m *EdiscoverySearch) SetDataSourceScopes(value *DataSourceScopes)() {
    err := m.GetBackingStore().Set("dataSourceScopes", value)
    if err != nil {
        panic(err)
    }
}
// SetLastEstimateStatisticsOperation sets the lastEstimateStatisticsOperation property value. The last estimate operation associated with the eDiscovery search.
func (m *EdiscoverySearch) SetLastEstimateStatisticsOperation(value EdiscoveryEstimateOperationable)() {
    err := m.GetBackingStore().Set("lastEstimateStatisticsOperation", value)
    if err != nil {
        panic(err)
    }
}
// SetNoncustodialSources sets the noncustodialSources property value. noncustodialDataSource sources that are included in the eDiscovery search
func (m *EdiscoverySearch) SetNoncustodialSources(value []EdiscoveryNoncustodialDataSourceable)() {
    err := m.GetBackingStore().Set("noncustodialSources", value)
    if err != nil {
        panic(err)
    }
}
type EdiscoverySearchable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    Searchable
    GetAdditionalSources()([]DataSourceable)
    GetAddToReviewSetOperation()(EdiscoveryAddToReviewSetOperationable)
    GetCustodianSources()([]DataSourceable)
    GetDataSourceScopes()(*DataSourceScopes)
    GetLastEstimateStatisticsOperation()(EdiscoveryEstimateOperationable)
    GetNoncustodialSources()([]EdiscoveryNoncustodialDataSourceable)
    SetAdditionalSources(value []DataSourceable)()
    SetAddToReviewSetOperation(value EdiscoveryAddToReviewSetOperationable)()
    SetCustodianSources(value []DataSourceable)()
    SetDataSourceScopes(value *DataSourceScopes)()
    SetLastEstimateStatisticsOperation(value EdiscoveryEstimateOperationable)()
    SetNoncustodialSources(value []EdiscoveryNoncustodialDataSourceable)()
}
