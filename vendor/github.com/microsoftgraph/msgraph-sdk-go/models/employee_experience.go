package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// EmployeeExperience represents a container that exposes navigation properties for employee experience resources.
type EmployeeExperience struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewEmployeeExperience instantiates a new EmployeeExperience and sets the default values.
func NewEmployeeExperience()(*EmployeeExperience) {
    m := &EmployeeExperience{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateEmployeeExperienceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEmployeeExperienceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEmployeeExperience(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *EmployeeExperience) GetAdditionalData()(map[string]any) {
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
func (m *EmployeeExperience) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCommunities gets the communities property value. A collection of communities in Viva Engage.
// returns a []Communityable when successful
func (m *EmployeeExperience) GetCommunities()([]Communityable) {
    val, err := m.GetBackingStore().Get("communities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Communityable)
    }
    return nil
}
// GetEngagementAsyncOperations gets the engagementAsyncOperations property value. A collection of long-running, asynchronous operations related to Viva Engage.
// returns a []EngagementAsyncOperationable when successful
func (m *EmployeeExperience) GetEngagementAsyncOperations()([]EngagementAsyncOperationable) {
    val, err := m.GetBackingStore().Get("engagementAsyncOperations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EngagementAsyncOperationable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EmployeeExperience) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["communities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCommunityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Communityable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Communityable)
                }
            }
            m.SetCommunities(res)
        }
        return nil
    }
    res["engagementAsyncOperations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEngagementAsyncOperationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EngagementAsyncOperationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EngagementAsyncOperationable)
                }
            }
            m.SetEngagementAsyncOperations(res)
        }
        return nil
    }
    res["learningCourseActivities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateLearningCourseActivityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]LearningCourseActivityable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(LearningCourseActivityable)
                }
            }
            m.SetLearningCourseActivities(res)
        }
        return nil
    }
    res["learningProviders"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateLearningProviderFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]LearningProviderable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(LearningProviderable)
                }
            }
            m.SetLearningProviders(res)
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
    return res
}
// GetLearningCourseActivities gets the learningCourseActivities property value. The learningCourseActivities property
// returns a []LearningCourseActivityable when successful
func (m *EmployeeExperience) GetLearningCourseActivities()([]LearningCourseActivityable) {
    val, err := m.GetBackingStore().Get("learningCourseActivities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]LearningCourseActivityable)
    }
    return nil
}
// GetLearningProviders gets the learningProviders property value. A collection of learning providers.
// returns a []LearningProviderable when successful
func (m *EmployeeExperience) GetLearningProviders()([]LearningProviderable) {
    val, err := m.GetBackingStore().Get("learningProviders")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]LearningProviderable)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *EmployeeExperience) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EmployeeExperience) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetCommunities() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCommunities()))
        for i, v := range m.GetCommunities() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("communities", cast)
        if err != nil {
            return err
        }
    }
    if m.GetEngagementAsyncOperations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetEngagementAsyncOperations()))
        for i, v := range m.GetEngagementAsyncOperations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("engagementAsyncOperations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetLearningCourseActivities() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLearningCourseActivities()))
        for i, v := range m.GetLearningCourseActivities() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("learningCourseActivities", cast)
        if err != nil {
            return err
        }
    }
    if m.GetLearningProviders() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLearningProviders()))
        for i, v := range m.GetLearningProviders() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("learningProviders", cast)
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
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *EmployeeExperience) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *EmployeeExperience) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCommunities sets the communities property value. A collection of communities in Viva Engage.
func (m *EmployeeExperience) SetCommunities(value []Communityable)() {
    err := m.GetBackingStore().Set("communities", value)
    if err != nil {
        panic(err)
    }
}
// SetEngagementAsyncOperations sets the engagementAsyncOperations property value. A collection of long-running, asynchronous operations related to Viva Engage.
func (m *EmployeeExperience) SetEngagementAsyncOperations(value []EngagementAsyncOperationable)() {
    err := m.GetBackingStore().Set("engagementAsyncOperations", value)
    if err != nil {
        panic(err)
    }
}
// SetLearningCourseActivities sets the learningCourseActivities property value. The learningCourseActivities property
func (m *EmployeeExperience) SetLearningCourseActivities(value []LearningCourseActivityable)() {
    err := m.GetBackingStore().Set("learningCourseActivities", value)
    if err != nil {
        panic(err)
    }
}
// SetLearningProviders sets the learningProviders property value. A collection of learning providers.
func (m *EmployeeExperience) SetLearningProviders(value []LearningProviderable)() {
    err := m.GetBackingStore().Set("learningProviders", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *EmployeeExperience) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type EmployeeExperienceable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCommunities()([]Communityable)
    GetEngagementAsyncOperations()([]EngagementAsyncOperationable)
    GetLearningCourseActivities()([]LearningCourseActivityable)
    GetLearningProviders()([]LearningProviderable)
    GetOdataType()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCommunities(value []Communityable)()
    SetEngagementAsyncOperations(value []EngagementAsyncOperationable)()
    SetLearningCourseActivities(value []LearningCourseActivityable)()
    SetLearningProviders(value []LearningProviderable)()
    SetOdataType(value *string)()
}
