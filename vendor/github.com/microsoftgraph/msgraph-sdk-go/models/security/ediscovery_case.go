package security

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type EdiscoveryCase struct {
    CaseEscaped
}
// NewEdiscoveryCase instantiates a new EdiscoveryCase and sets the default values.
func NewEdiscoveryCase()(*EdiscoveryCase) {
    m := &EdiscoveryCase{
        CaseEscaped: *NewCaseEscaped(),
    }
    odataTypeValue := "#microsoft.graph.security.ediscoveryCase"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEdiscoveryCaseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEdiscoveryCaseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEdiscoveryCase(), nil
}
// GetClosedBy gets the closedBy property value. The user who closed the case.
// returns a IdentitySetable when successful
func (m *EdiscoveryCase) GetClosedBy()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable) {
    val, err := m.GetBackingStore().Get("closedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable)
    }
    return nil
}
// GetClosedDateTime gets the closedDateTime property value. The date and time when the case was closed. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
// returns a *Time when successful
func (m *EdiscoveryCase) GetClosedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("closedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetCustodians gets the custodians property value. Returns a list of case ediscoveryCustodian objects for this case.
// returns a []EdiscoveryCustodianable when successful
func (m *EdiscoveryCase) GetCustodians()([]EdiscoveryCustodianable) {
    val, err := m.GetBackingStore().Get("custodians")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EdiscoveryCustodianable)
    }
    return nil
}
// GetExternalId gets the externalId property value. The external case number for customer reference.
// returns a *string when successful
func (m *EdiscoveryCase) GetExternalId()(*string) {
    val, err := m.GetBackingStore().Get("externalId")
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
func (m *EdiscoveryCase) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.CaseEscaped.GetFieldDeserializers()
    res["closedBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClosedBy(val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable))
        }
        return nil
    }
    res["closedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClosedDateTime(val)
        }
        return nil
    }
    res["custodians"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEdiscoveryCustodianFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EdiscoveryCustodianable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EdiscoveryCustodianable)
                }
            }
            m.SetCustodians(res)
        }
        return nil
    }
    res["externalId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalId(val)
        }
        return nil
    }
    res["noncustodialDataSources"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetNoncustodialDataSources(res)
        }
        return nil
    }
    res["operations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCaseOperationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CaseOperationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CaseOperationable)
                }
            }
            m.SetOperations(res)
        }
        return nil
    }
    res["reviewSets"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEdiscoveryReviewSetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EdiscoveryReviewSetable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EdiscoveryReviewSetable)
                }
            }
            m.SetReviewSets(res)
        }
        return nil
    }
    res["searches"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEdiscoverySearchFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EdiscoverySearchable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EdiscoverySearchable)
                }
            }
            m.SetSearches(res)
        }
        return nil
    }
    res["settings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEdiscoveryCaseSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettings(val.(EdiscoveryCaseSettingsable))
        }
        return nil
    }
    res["tags"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEdiscoveryReviewTagFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EdiscoveryReviewTagable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EdiscoveryReviewTagable)
                }
            }
            m.SetTags(res)
        }
        return nil
    }
    return res
}
// GetNoncustodialDataSources gets the noncustodialDataSources property value. Returns a list of case ediscoveryNoncustodialDataSource objects for this case.
// returns a []EdiscoveryNoncustodialDataSourceable when successful
func (m *EdiscoveryCase) GetNoncustodialDataSources()([]EdiscoveryNoncustodialDataSourceable) {
    val, err := m.GetBackingStore().Get("noncustodialDataSources")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EdiscoveryNoncustodialDataSourceable)
    }
    return nil
}
// GetOperations gets the operations property value. Returns a list of case caseOperation objects for this case.
// returns a []CaseOperationable when successful
func (m *EdiscoveryCase) GetOperations()([]CaseOperationable) {
    val, err := m.GetBackingStore().Get("operations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CaseOperationable)
    }
    return nil
}
// GetReviewSets gets the reviewSets property value. Returns a list of eDiscoveryReviewSet objects in the case.
// returns a []EdiscoveryReviewSetable when successful
func (m *EdiscoveryCase) GetReviewSets()([]EdiscoveryReviewSetable) {
    val, err := m.GetBackingStore().Get("reviewSets")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EdiscoveryReviewSetable)
    }
    return nil
}
// GetSearches gets the searches property value. Returns a list of eDiscoverySearch objects associated with this case.
// returns a []EdiscoverySearchable when successful
func (m *EdiscoveryCase) GetSearches()([]EdiscoverySearchable) {
    val, err := m.GetBackingStore().Get("searches")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EdiscoverySearchable)
    }
    return nil
}
// GetSettings gets the settings property value. Returns a list of eDIscoverySettings objects in the case.
// returns a EdiscoveryCaseSettingsable when successful
func (m *EdiscoveryCase) GetSettings()(EdiscoveryCaseSettingsable) {
    val, err := m.GetBackingStore().Get("settings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EdiscoveryCaseSettingsable)
    }
    return nil
}
// GetTags gets the tags property value. Returns a list of ediscoveryReviewTag objects associated to this case.
// returns a []EdiscoveryReviewTagable when successful
func (m *EdiscoveryCase) GetTags()([]EdiscoveryReviewTagable) {
    val, err := m.GetBackingStore().Get("tags")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EdiscoveryReviewTagable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EdiscoveryCase) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.CaseEscaped.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("closedBy", m.GetClosedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("closedDateTime", m.GetClosedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetCustodians() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCustodians()))
        for i, v := range m.GetCustodians() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("custodians", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("externalId", m.GetExternalId())
        if err != nil {
            return err
        }
    }
    if m.GetNoncustodialDataSources() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetNoncustodialDataSources()))
        for i, v := range m.GetNoncustodialDataSources() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("noncustodialDataSources", cast)
        if err != nil {
            return err
        }
    }
    if m.GetOperations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOperations()))
        for i, v := range m.GetOperations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("operations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetReviewSets() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetReviewSets()))
        for i, v := range m.GetReviewSets() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("reviewSets", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSearches() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSearches()))
        for i, v := range m.GetSearches() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("searches", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("settings", m.GetSettings())
        if err != nil {
            return err
        }
    }
    if m.GetTags() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTags()))
        for i, v := range m.GetTags() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("tags", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetClosedBy sets the closedBy property value. The user who closed the case.
func (m *EdiscoveryCase) SetClosedBy(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable)() {
    err := m.GetBackingStore().Set("closedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetClosedDateTime sets the closedDateTime property value. The date and time when the case was closed. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *EdiscoveryCase) SetClosedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("closedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetCustodians sets the custodians property value. Returns a list of case ediscoveryCustodian objects for this case.
func (m *EdiscoveryCase) SetCustodians(value []EdiscoveryCustodianable)() {
    err := m.GetBackingStore().Set("custodians", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalId sets the externalId property value. The external case number for customer reference.
func (m *EdiscoveryCase) SetExternalId(value *string)() {
    err := m.GetBackingStore().Set("externalId", value)
    if err != nil {
        panic(err)
    }
}
// SetNoncustodialDataSources sets the noncustodialDataSources property value. Returns a list of case ediscoveryNoncustodialDataSource objects for this case.
func (m *EdiscoveryCase) SetNoncustodialDataSources(value []EdiscoveryNoncustodialDataSourceable)() {
    err := m.GetBackingStore().Set("noncustodialDataSources", value)
    if err != nil {
        panic(err)
    }
}
// SetOperations sets the operations property value. Returns a list of case caseOperation objects for this case.
func (m *EdiscoveryCase) SetOperations(value []CaseOperationable)() {
    err := m.GetBackingStore().Set("operations", value)
    if err != nil {
        panic(err)
    }
}
// SetReviewSets sets the reviewSets property value. Returns a list of eDiscoveryReviewSet objects in the case.
func (m *EdiscoveryCase) SetReviewSets(value []EdiscoveryReviewSetable)() {
    err := m.GetBackingStore().Set("reviewSets", value)
    if err != nil {
        panic(err)
    }
}
// SetSearches sets the searches property value. Returns a list of eDiscoverySearch objects associated with this case.
func (m *EdiscoveryCase) SetSearches(value []EdiscoverySearchable)() {
    err := m.GetBackingStore().Set("searches", value)
    if err != nil {
        panic(err)
    }
}
// SetSettings sets the settings property value. Returns a list of eDIscoverySettings objects in the case.
func (m *EdiscoveryCase) SetSettings(value EdiscoveryCaseSettingsable)() {
    err := m.GetBackingStore().Set("settings", value)
    if err != nil {
        panic(err)
    }
}
// SetTags sets the tags property value. Returns a list of ediscoveryReviewTag objects associated to this case.
func (m *EdiscoveryCase) SetTags(value []EdiscoveryReviewTagable)() {
    err := m.GetBackingStore().Set("tags", value)
    if err != nil {
        panic(err)
    }
}
type EdiscoveryCaseable interface {
    CaseEscapedable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetClosedBy()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable)
    GetClosedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetCustodians()([]EdiscoveryCustodianable)
    GetExternalId()(*string)
    GetNoncustodialDataSources()([]EdiscoveryNoncustodialDataSourceable)
    GetOperations()([]CaseOperationable)
    GetReviewSets()([]EdiscoveryReviewSetable)
    GetSearches()([]EdiscoverySearchable)
    GetSettings()(EdiscoveryCaseSettingsable)
    GetTags()([]EdiscoveryReviewTagable)
    SetClosedBy(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable)()
    SetClosedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetCustodians(value []EdiscoveryCustodianable)()
    SetExternalId(value *string)()
    SetNoncustodialDataSources(value []EdiscoveryNoncustodialDataSourceable)()
    SetOperations(value []CaseOperationable)()
    SetReviewSets(value []EdiscoveryReviewSetable)()
    SetSearches(value []EdiscoverySearchable)()
    SetSettings(value EdiscoveryCaseSettingsable)()
    SetTags(value []EdiscoveryReviewTagable)()
}
