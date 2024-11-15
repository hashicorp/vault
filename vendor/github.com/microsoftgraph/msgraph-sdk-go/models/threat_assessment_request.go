package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ThreatAssessmentRequest struct {
    Entity
}
// NewThreatAssessmentRequest instantiates a new ThreatAssessmentRequest and sets the default values.
func NewThreatAssessmentRequest()(*ThreatAssessmentRequest) {
    m := &ThreatAssessmentRequest{
        Entity: *NewEntity(),
    }
    return m
}
// CreateThreatAssessmentRequestFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateThreatAssessmentRequestFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.emailFileAssessmentRequest":
                        return NewEmailFileAssessmentRequest(), nil
                    case "#microsoft.graph.fileAssessmentRequest":
                        return NewFileAssessmentRequest(), nil
                    case "#microsoft.graph.mailAssessmentRequest":
                        return NewMailAssessmentRequest(), nil
                    case "#microsoft.graph.urlAssessmentRequest":
                        return NewUrlAssessmentRequest(), nil
                }
            }
        }
    }
    return NewThreatAssessmentRequest(), nil
}
// GetCategory gets the category property value. The category property
// returns a *ThreatCategory when successful
func (m *ThreatAssessmentRequest) GetCategory()(*ThreatCategory) {
    val, err := m.GetBackingStore().Get("category")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ThreatCategory)
    }
    return nil
}
// GetContentType gets the contentType property value. The content type of threat assessment. Possible values are: mail, url, file.
// returns a *ThreatAssessmentContentType when successful
func (m *ThreatAssessmentRequest) GetContentType()(*ThreatAssessmentContentType) {
    val, err := m.GetBackingStore().Get("contentType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ThreatAssessmentContentType)
    }
    return nil
}
// GetCreatedBy gets the createdBy property value. The threat assessment request creator.
// returns a IdentitySetable when successful
func (m *ThreatAssessmentRequest) GetCreatedBy()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("createdBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *ThreatAssessmentRequest) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetExpectedAssessment gets the expectedAssessment property value. The expectedAssessment property
// returns a *ThreatExpectedAssessment when successful
func (m *ThreatAssessmentRequest) GetExpectedAssessment()(*ThreatExpectedAssessment) {
    val, err := m.GetBackingStore().Get("expectedAssessment")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ThreatExpectedAssessment)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ThreatAssessmentRequest) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["category"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseThreatCategory)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory(val.(*ThreatCategory))
        }
        return nil
    }
    res["contentType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseThreatAssessmentContentType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContentType(val.(*ThreatAssessmentContentType))
        }
        return nil
    }
    res["createdBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedBy(val.(IdentitySetable))
        }
        return nil
    }
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
        }
        return nil
    }
    res["expectedAssessment"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseThreatExpectedAssessment)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExpectedAssessment(val.(*ThreatExpectedAssessment))
        }
        return nil
    }
    res["requestSource"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseThreatAssessmentRequestSource)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRequestSource(val.(*ThreatAssessmentRequestSource))
        }
        return nil
    }
    res["results"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateThreatAssessmentResultFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ThreatAssessmentResultable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ThreatAssessmentResultable)
                }
            }
            m.SetResults(res)
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseThreatAssessmentStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*ThreatAssessmentStatus))
        }
        return nil
    }
    return res
}
// GetRequestSource gets the requestSource property value. The source of the threat assessment request. Possible values are: administrator.
// returns a *ThreatAssessmentRequestSource when successful
func (m *ThreatAssessmentRequest) GetRequestSource()(*ThreatAssessmentRequestSource) {
    val, err := m.GetBackingStore().Get("requestSource")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ThreatAssessmentRequestSource)
    }
    return nil
}
// GetResults gets the results property value. A collection of threat assessment results. Read-only. By default, a GET /threatAssessmentRequests/{id} does not return this property unless you apply $expand on it.
// returns a []ThreatAssessmentResultable when successful
func (m *ThreatAssessmentRequest) GetResults()([]ThreatAssessmentResultable) {
    val, err := m.GetBackingStore().Get("results")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ThreatAssessmentResultable)
    }
    return nil
}
// GetStatus gets the status property value. The assessment process status. Possible values are: pending, completed.
// returns a *ThreatAssessmentStatus when successful
func (m *ThreatAssessmentRequest) GetStatus()(*ThreatAssessmentStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ThreatAssessmentStatus)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ThreatAssessmentRequest) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCategory() != nil {
        cast := (*m.GetCategory()).String()
        err = writer.WriteStringValue("category", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetContentType() != nil {
        cast := (*m.GetContentType()).String()
        err = writer.WriteStringValue("contentType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("createdBy", m.GetCreatedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetExpectedAssessment() != nil {
        cast := (*m.GetExpectedAssessment()).String()
        err = writer.WriteStringValue("expectedAssessment", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetRequestSource() != nil {
        cast := (*m.GetRequestSource()).String()
        err = writer.WriteStringValue("requestSource", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetResults() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetResults()))
        for i, v := range m.GetResults() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("results", cast)
        if err != nil {
            return err
        }
    }
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err = writer.WriteStringValue("status", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCategory sets the category property value. The category property
func (m *ThreatAssessmentRequest) SetCategory(value *ThreatCategory)() {
    err := m.GetBackingStore().Set("category", value)
    if err != nil {
        panic(err)
    }
}
// SetContentType sets the contentType property value. The content type of threat assessment. Possible values are: mail, url, file.
func (m *ThreatAssessmentRequest) SetContentType(value *ThreatAssessmentContentType)() {
    err := m.GetBackingStore().Set("contentType", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedBy sets the createdBy property value. The threat assessment request creator.
func (m *ThreatAssessmentRequest) SetCreatedBy(value IdentitySetable)() {
    err := m.GetBackingStore().Set("createdBy", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *ThreatAssessmentRequest) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetExpectedAssessment sets the expectedAssessment property value. The expectedAssessment property
func (m *ThreatAssessmentRequest) SetExpectedAssessment(value *ThreatExpectedAssessment)() {
    err := m.GetBackingStore().Set("expectedAssessment", value)
    if err != nil {
        panic(err)
    }
}
// SetRequestSource sets the requestSource property value. The source of the threat assessment request. Possible values are: administrator.
func (m *ThreatAssessmentRequest) SetRequestSource(value *ThreatAssessmentRequestSource)() {
    err := m.GetBackingStore().Set("requestSource", value)
    if err != nil {
        panic(err)
    }
}
// SetResults sets the results property value. A collection of threat assessment results. Read-only. By default, a GET /threatAssessmentRequests/{id} does not return this property unless you apply $expand on it.
func (m *ThreatAssessmentRequest) SetResults(value []ThreatAssessmentResultable)() {
    err := m.GetBackingStore().Set("results", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The assessment process status. Possible values are: pending, completed.
func (m *ThreatAssessmentRequest) SetStatus(value *ThreatAssessmentStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
type ThreatAssessmentRequestable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCategory()(*ThreatCategory)
    GetContentType()(*ThreatAssessmentContentType)
    GetCreatedBy()(IdentitySetable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetExpectedAssessment()(*ThreatExpectedAssessment)
    GetRequestSource()(*ThreatAssessmentRequestSource)
    GetResults()([]ThreatAssessmentResultable)
    GetStatus()(*ThreatAssessmentStatus)
    SetCategory(value *ThreatCategory)()
    SetContentType(value *ThreatAssessmentContentType)()
    SetCreatedBy(value IdentitySetable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetExpectedAssessment(value *ThreatExpectedAssessment)()
    SetRequestSource(value *ThreatAssessmentRequestSource)()
    SetResults(value []ThreatAssessmentResultable)()
    SetStatus(value *ThreatAssessmentStatus)()
}
