package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Training struct {
    Entity
}
// NewTraining instantiates a new Training and sets the default values.
func NewTraining()(*Training) {
    m := &Training{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTrainingFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTrainingFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTraining(), nil
}
// GetAvailabilityStatus gets the availabilityStatus property value. Training availability status. Possible values are: unknown, notAvailable, available, archive, delete, unknownFutureValue.
// returns a *TrainingAvailabilityStatus when successful
func (m *Training) GetAvailabilityStatus()(*TrainingAvailabilityStatus) {
    val, err := m.GetBackingStore().Get("availabilityStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TrainingAvailabilityStatus)
    }
    return nil
}
// GetCreatedBy gets the createdBy property value. Identity of the user who created the training.
// returns a EmailIdentityable when successful
func (m *Training) GetCreatedBy()(EmailIdentityable) {
    val, err := m.GetBackingStore().Get("createdBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EmailIdentityable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Date and time when the training was created. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *Training) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDescription gets the description property value. The description for the training.
// returns a *string when successful
func (m *Training) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The display name for the training.
// returns a *string when successful
func (m *Training) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDurationInMinutes gets the durationInMinutes property value. Training duration.
// returns a *int32 when successful
func (m *Training) GetDurationInMinutes()(*int32) {
    val, err := m.GetBackingStore().Get("durationInMinutes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Training) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["availabilityStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTrainingAvailabilityStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAvailabilityStatus(val.(*TrainingAvailabilityStatus))
        }
        return nil
    }
    res["createdBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEmailIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedBy(val.(EmailIdentityable))
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
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
        }
        return nil
    }
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["durationInMinutes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDurationInMinutes(val)
        }
        return nil
    }
    res["hasEvaluation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHasEvaluation(val)
        }
        return nil
    }
    res["languageDetails"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTrainingLanguageDetailFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TrainingLanguageDetailable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TrainingLanguageDetailable)
                }
            }
            m.SetLanguageDetails(res)
        }
        return nil
    }
    res["lastModifiedBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEmailIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedBy(val.(EmailIdentityable))
        }
        return nil
    }
    res["lastModifiedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedDateTime(val)
        }
        return nil
    }
    res["source"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSimulationContentSource)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSource(val.(*SimulationContentSource))
        }
        return nil
    }
    res["supportedLocales"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSupportedLocales(res)
        }
        return nil
    }
    res["tags"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetTags(res)
        }
        return nil
    }
    res["type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTrainingType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTypeEscaped(val.(*TrainingType))
        }
        return nil
    }
    return res
}
// GetHasEvaluation gets the hasEvaluation property value. Indicates whether the training has any evaluation.
// returns a *bool when successful
func (m *Training) GetHasEvaluation()(*bool) {
    val, err := m.GetBackingStore().Get("hasEvaluation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLanguageDetails gets the languageDetails property value. Language specific details on a training.
// returns a []TrainingLanguageDetailable when successful
func (m *Training) GetLanguageDetails()([]TrainingLanguageDetailable) {
    val, err := m.GetBackingStore().Get("languageDetails")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TrainingLanguageDetailable)
    }
    return nil
}
// GetLastModifiedBy gets the lastModifiedBy property value. Identity of the user who last modified the training.
// returns a EmailIdentityable when successful
func (m *Training) GetLastModifiedBy()(EmailIdentityable) {
    val, err := m.GetBackingStore().Get("lastModifiedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EmailIdentityable)
    }
    return nil
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. Date and time when the training was last modified. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *Training) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetSource gets the source property value. Training content source. Possible values are: unknown, global, tenant, unknownFutureValue.
// returns a *SimulationContentSource when successful
func (m *Training) GetSource()(*SimulationContentSource) {
    val, err := m.GetBackingStore().Get("source")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SimulationContentSource)
    }
    return nil
}
// GetSupportedLocales gets the supportedLocales property value. Supported locales for content for the associated training.
// returns a []string when successful
func (m *Training) GetSupportedLocales()([]string) {
    val, err := m.GetBackingStore().Get("supportedLocales")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetTags gets the tags property value. Training tags.
// returns a []string when successful
func (m *Training) GetTags()([]string) {
    val, err := m.GetBackingStore().Get("tags")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetTypeEscaped gets the type property value. The type of training. Possible values are: unknown, phishing, unknownFutureValue.
// returns a *TrainingType when successful
func (m *Training) GetTypeEscaped()(*TrainingType) {
    val, err := m.GetBackingStore().Get("typeEscaped")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TrainingType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Training) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAvailabilityStatus() != nil {
        cast := (*m.GetAvailabilityStatus()).String()
        err = writer.WriteStringValue("availabilityStatus", &cast)
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
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("durationInMinutes", m.GetDurationInMinutes())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("hasEvaluation", m.GetHasEvaluation())
        if err != nil {
            return err
        }
    }
    if m.GetLanguageDetails() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLanguageDetails()))
        for i, v := range m.GetLanguageDetails() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("languageDetails", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("lastModifiedBy", m.GetLastModifiedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastModifiedDateTime", m.GetLastModifiedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetSource() != nil {
        cast := (*m.GetSource()).String()
        err = writer.WriteStringValue("source", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetSupportedLocales() != nil {
        err = writer.WriteCollectionOfStringValues("supportedLocales", m.GetSupportedLocales())
        if err != nil {
            return err
        }
    }
    if m.GetTags() != nil {
        err = writer.WriteCollectionOfStringValues("tags", m.GetTags())
        if err != nil {
            return err
        }
    }
    if m.GetTypeEscaped() != nil {
        cast := (*m.GetTypeEscaped()).String()
        err = writer.WriteStringValue("type", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAvailabilityStatus sets the availabilityStatus property value. Training availability status. Possible values are: unknown, notAvailable, available, archive, delete, unknownFutureValue.
func (m *Training) SetAvailabilityStatus(value *TrainingAvailabilityStatus)() {
    err := m.GetBackingStore().Set("availabilityStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedBy sets the createdBy property value. Identity of the user who created the training.
func (m *Training) SetCreatedBy(value EmailIdentityable)() {
    err := m.GetBackingStore().Set("createdBy", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Date and time when the training was created. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *Training) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. The description for the training.
func (m *Training) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The display name for the training.
func (m *Training) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetDurationInMinutes sets the durationInMinutes property value. Training duration.
func (m *Training) SetDurationInMinutes(value *int32)() {
    err := m.GetBackingStore().Set("durationInMinutes", value)
    if err != nil {
        panic(err)
    }
}
// SetHasEvaluation sets the hasEvaluation property value. Indicates whether the training has any evaluation.
func (m *Training) SetHasEvaluation(value *bool)() {
    err := m.GetBackingStore().Set("hasEvaluation", value)
    if err != nil {
        panic(err)
    }
}
// SetLanguageDetails sets the languageDetails property value. Language specific details on a training.
func (m *Training) SetLanguageDetails(value []TrainingLanguageDetailable)() {
    err := m.GetBackingStore().Set("languageDetails", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedBy sets the lastModifiedBy property value. Identity of the user who last modified the training.
func (m *Training) SetLastModifiedBy(value EmailIdentityable)() {
    err := m.GetBackingStore().Set("lastModifiedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. Date and time when the training was last modified. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *Training) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetSource sets the source property value. Training content source. Possible values are: unknown, global, tenant, unknownFutureValue.
func (m *Training) SetSource(value *SimulationContentSource)() {
    err := m.GetBackingStore().Set("source", value)
    if err != nil {
        panic(err)
    }
}
// SetSupportedLocales sets the supportedLocales property value. Supported locales for content for the associated training.
func (m *Training) SetSupportedLocales(value []string)() {
    err := m.GetBackingStore().Set("supportedLocales", value)
    if err != nil {
        panic(err)
    }
}
// SetTags sets the tags property value. Training tags.
func (m *Training) SetTags(value []string)() {
    err := m.GetBackingStore().Set("tags", value)
    if err != nil {
        panic(err)
    }
}
// SetTypeEscaped sets the type property value. The type of training. Possible values are: unknown, phishing, unknownFutureValue.
func (m *Training) SetTypeEscaped(value *TrainingType)() {
    err := m.GetBackingStore().Set("typeEscaped", value)
    if err != nil {
        panic(err)
    }
}
type Trainingable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAvailabilityStatus()(*TrainingAvailabilityStatus)
    GetCreatedBy()(EmailIdentityable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetDurationInMinutes()(*int32)
    GetHasEvaluation()(*bool)
    GetLanguageDetails()([]TrainingLanguageDetailable)
    GetLastModifiedBy()(EmailIdentityable)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetSource()(*SimulationContentSource)
    GetSupportedLocales()([]string)
    GetTags()([]string)
    GetTypeEscaped()(*TrainingType)
    SetAvailabilityStatus(value *TrainingAvailabilityStatus)()
    SetCreatedBy(value EmailIdentityable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetDurationInMinutes(value *int32)()
    SetHasEvaluation(value *bool)()
    SetLanguageDetails(value []TrainingLanguageDetailable)()
    SetLastModifiedBy(value EmailIdentityable)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetSource(value *SimulationContentSource)()
    SetSupportedLocales(value []string)()
    SetTags(value []string)()
    SetTypeEscaped(value *TrainingType)()
}
