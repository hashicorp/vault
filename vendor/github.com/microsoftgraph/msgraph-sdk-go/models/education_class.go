package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EducationClass struct {
    Entity
}
// NewEducationClass instantiates a new EducationClass and sets the default values.
func NewEducationClass()(*EducationClass) {
    m := &EducationClass{
        Entity: *NewEntity(),
    }
    return m
}
// CreateEducationClassFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEducationClassFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEducationClass(), nil
}
// GetAssignmentCategories gets the assignmentCategories property value. All categories associated with this class. Nullable.
// returns a []EducationCategoryable when successful
func (m *EducationClass) GetAssignmentCategories()([]EducationCategoryable) {
    val, err := m.GetBackingStore().Get("assignmentCategories")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EducationCategoryable)
    }
    return nil
}
// GetAssignmentDefaults gets the assignmentDefaults property value. Specifies class-level defaults respected by new assignments created in the class.
// returns a EducationAssignmentDefaultsable when successful
func (m *EducationClass) GetAssignmentDefaults()(EducationAssignmentDefaultsable) {
    val, err := m.GetBackingStore().Get("assignmentDefaults")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EducationAssignmentDefaultsable)
    }
    return nil
}
// GetAssignments gets the assignments property value. All assignments associated with this class. Nullable.
// returns a []EducationAssignmentable when successful
func (m *EducationClass) GetAssignments()([]EducationAssignmentable) {
    val, err := m.GetBackingStore().Get("assignments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EducationAssignmentable)
    }
    return nil
}
// GetAssignmentSettings gets the assignmentSettings property value. Specifies class-level assignments settings.
// returns a EducationAssignmentSettingsable when successful
func (m *EducationClass) GetAssignmentSettings()(EducationAssignmentSettingsable) {
    val, err := m.GetBackingStore().Get("assignmentSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EducationAssignmentSettingsable)
    }
    return nil
}
// GetClassCode gets the classCode property value. Class code used by the school to identify the class.
// returns a *string when successful
func (m *EducationClass) GetClassCode()(*string) {
    val, err := m.GetBackingStore().Get("classCode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCourse gets the course property value. The course property
// returns a EducationCourseable when successful
func (m *EducationClass) GetCourse()(EducationCourseable) {
    val, err := m.GetBackingStore().Get("course")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EducationCourseable)
    }
    return nil
}
// GetCreatedBy gets the createdBy property value. Entity who created the class
// returns a IdentitySetable when successful
func (m *EducationClass) GetCreatedBy()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("createdBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetDescription gets the description property value. Description of the class.
// returns a *string when successful
func (m *EducationClass) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Name of the class.
// returns a *string when successful
func (m *EducationClass) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExternalId gets the externalId property value. ID of the class from the syncing system.
// returns a *string when successful
func (m *EducationClass) GetExternalId()(*string) {
    val, err := m.GetBackingStore().Get("externalId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExternalName gets the externalName property value. Name of the class in the syncing system.
// returns a *string when successful
func (m *EducationClass) GetExternalName()(*string) {
    val, err := m.GetBackingStore().Get("externalName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExternalSource gets the externalSource property value. How this class was created. Possible values are: sis, manual.
// returns a *EducationExternalSource when successful
func (m *EducationClass) GetExternalSource()(*EducationExternalSource) {
    val, err := m.GetBackingStore().Get("externalSource")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*EducationExternalSource)
    }
    return nil
}
// GetExternalSourceDetail gets the externalSourceDetail property value. The name of the external source this resource was generated from.
// returns a *string when successful
func (m *EducationClass) GetExternalSourceDetail()(*string) {
    val, err := m.GetBackingStore().Get("externalSourceDetail")
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
func (m *EducationClass) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["assignmentCategories"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEducationCategoryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EducationCategoryable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EducationCategoryable)
                }
            }
            m.SetAssignmentCategories(res)
        }
        return nil
    }
    res["assignmentDefaults"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEducationAssignmentDefaultsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignmentDefaults(val.(EducationAssignmentDefaultsable))
        }
        return nil
    }
    res["assignments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEducationAssignmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EducationAssignmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EducationAssignmentable)
                }
            }
            m.SetAssignments(res)
        }
        return nil
    }
    res["assignmentSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEducationAssignmentSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignmentSettings(val.(EducationAssignmentSettingsable))
        }
        return nil
    }
    res["classCode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClassCode(val)
        }
        return nil
    }
    res["course"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEducationCourseFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCourse(val.(EducationCourseable))
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
    res["externalName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalName(val)
        }
        return nil
    }
    res["externalSource"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseEducationExternalSource)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalSource(val.(*EducationExternalSource))
        }
        return nil
    }
    res["externalSourceDetail"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalSourceDetail(val)
        }
        return nil
    }
    res["grade"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGrade(val)
        }
        return nil
    }
    res["group"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateGroupFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGroup(val.(Groupable))
        }
        return nil
    }
    res["mailNickname"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMailNickname(val)
        }
        return nil
    }
    res["members"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEducationUserFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EducationUserable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EducationUserable)
                }
            }
            m.SetMembers(res)
        }
        return nil
    }
    res["modules"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEducationModuleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EducationModuleable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EducationModuleable)
                }
            }
            m.SetModules(res)
        }
        return nil
    }
    res["schools"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEducationSchoolFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EducationSchoolable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EducationSchoolable)
                }
            }
            m.SetSchools(res)
        }
        return nil
    }
    res["teachers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEducationUserFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EducationUserable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EducationUserable)
                }
            }
            m.SetTeachers(res)
        }
        return nil
    }
    res["term"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEducationTermFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTerm(val.(EducationTermable))
        }
        return nil
    }
    return res
}
// GetGrade gets the grade property value. Grade level of the class.
// returns a *string when successful
func (m *EducationClass) GetGrade()(*string) {
    val, err := m.GetBackingStore().Get("grade")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetGroup gets the group property value. The underlying Microsoft 365 group object.
// returns a Groupable when successful
func (m *EducationClass) GetGroup()(Groupable) {
    val, err := m.GetBackingStore().Get("group")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Groupable)
    }
    return nil
}
// GetMailNickname gets the mailNickname property value. Mail name for sending email to all members, if this is enabled.
// returns a *string when successful
func (m *EducationClass) GetMailNickname()(*string) {
    val, err := m.GetBackingStore().Get("mailNickname")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMembers gets the members property value. All users in the class. Nullable.
// returns a []EducationUserable when successful
func (m *EducationClass) GetMembers()([]EducationUserable) {
    val, err := m.GetBackingStore().Get("members")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EducationUserable)
    }
    return nil
}
// GetModules gets the modules property value. All modules in the class. Nullable.
// returns a []EducationModuleable when successful
func (m *EducationClass) GetModules()([]EducationModuleable) {
    val, err := m.GetBackingStore().Get("modules")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EducationModuleable)
    }
    return nil
}
// GetSchools gets the schools property value. All schools that this class is associated with. Nullable.
// returns a []EducationSchoolable when successful
func (m *EducationClass) GetSchools()([]EducationSchoolable) {
    val, err := m.GetBackingStore().Get("schools")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EducationSchoolable)
    }
    return nil
}
// GetTeachers gets the teachers property value. All teachers in the class. Nullable.
// returns a []EducationUserable when successful
func (m *EducationClass) GetTeachers()([]EducationUserable) {
    val, err := m.GetBackingStore().Get("teachers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EducationUserable)
    }
    return nil
}
// GetTerm gets the term property value. Term for this class.
// returns a EducationTermable when successful
func (m *EducationClass) GetTerm()(EducationTermable) {
    val, err := m.GetBackingStore().Get("term")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EducationTermable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EducationClass) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAssignmentCategories() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAssignmentCategories()))
        for i, v := range m.GetAssignmentCategories() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("assignmentCategories", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("assignmentDefaults", m.GetAssignmentDefaults())
        if err != nil {
            return err
        }
    }
    if m.GetAssignments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAssignments()))
        for i, v := range m.GetAssignments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("assignments", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("assignmentSettings", m.GetAssignmentSettings())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("classCode", m.GetClassCode())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("course", m.GetCourse())
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
        err = writer.WriteStringValue("externalId", m.GetExternalId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("externalName", m.GetExternalName())
        if err != nil {
            return err
        }
    }
    if m.GetExternalSource() != nil {
        cast := (*m.GetExternalSource()).String()
        err = writer.WriteStringValue("externalSource", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("externalSourceDetail", m.GetExternalSourceDetail())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("grade", m.GetGrade())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("group", m.GetGroup())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("mailNickname", m.GetMailNickname())
        if err != nil {
            return err
        }
    }
    if m.GetMembers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMembers()))
        for i, v := range m.GetMembers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("members", cast)
        if err != nil {
            return err
        }
    }
    if m.GetModules() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetModules()))
        for i, v := range m.GetModules() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("modules", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSchools() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSchools()))
        for i, v := range m.GetSchools() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("schools", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTeachers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTeachers()))
        for i, v := range m.GetTeachers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("teachers", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("term", m.GetTerm())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAssignmentCategories sets the assignmentCategories property value. All categories associated with this class. Nullable.
func (m *EducationClass) SetAssignmentCategories(value []EducationCategoryable)() {
    err := m.GetBackingStore().Set("assignmentCategories", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignmentDefaults sets the assignmentDefaults property value. Specifies class-level defaults respected by new assignments created in the class.
func (m *EducationClass) SetAssignmentDefaults(value EducationAssignmentDefaultsable)() {
    err := m.GetBackingStore().Set("assignmentDefaults", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignments sets the assignments property value. All assignments associated with this class. Nullable.
func (m *EducationClass) SetAssignments(value []EducationAssignmentable)() {
    err := m.GetBackingStore().Set("assignments", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignmentSettings sets the assignmentSettings property value. Specifies class-level assignments settings.
func (m *EducationClass) SetAssignmentSettings(value EducationAssignmentSettingsable)() {
    err := m.GetBackingStore().Set("assignmentSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetClassCode sets the classCode property value. Class code used by the school to identify the class.
func (m *EducationClass) SetClassCode(value *string)() {
    err := m.GetBackingStore().Set("classCode", value)
    if err != nil {
        panic(err)
    }
}
// SetCourse sets the course property value. The course property
func (m *EducationClass) SetCourse(value EducationCourseable)() {
    err := m.GetBackingStore().Set("course", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedBy sets the createdBy property value. Entity who created the class
func (m *EducationClass) SetCreatedBy(value IdentitySetable)() {
    err := m.GetBackingStore().Set("createdBy", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. Description of the class.
func (m *EducationClass) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Name of the class.
func (m *EducationClass) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalId sets the externalId property value. ID of the class from the syncing system.
func (m *EducationClass) SetExternalId(value *string)() {
    err := m.GetBackingStore().Set("externalId", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalName sets the externalName property value. Name of the class in the syncing system.
func (m *EducationClass) SetExternalName(value *string)() {
    err := m.GetBackingStore().Set("externalName", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalSource sets the externalSource property value. How this class was created. Possible values are: sis, manual.
func (m *EducationClass) SetExternalSource(value *EducationExternalSource)() {
    err := m.GetBackingStore().Set("externalSource", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalSourceDetail sets the externalSourceDetail property value. The name of the external source this resource was generated from.
func (m *EducationClass) SetExternalSourceDetail(value *string)() {
    err := m.GetBackingStore().Set("externalSourceDetail", value)
    if err != nil {
        panic(err)
    }
}
// SetGrade sets the grade property value. Grade level of the class.
func (m *EducationClass) SetGrade(value *string)() {
    err := m.GetBackingStore().Set("grade", value)
    if err != nil {
        panic(err)
    }
}
// SetGroup sets the group property value. The underlying Microsoft 365 group object.
func (m *EducationClass) SetGroup(value Groupable)() {
    err := m.GetBackingStore().Set("group", value)
    if err != nil {
        panic(err)
    }
}
// SetMailNickname sets the mailNickname property value. Mail name for sending email to all members, if this is enabled.
func (m *EducationClass) SetMailNickname(value *string)() {
    err := m.GetBackingStore().Set("mailNickname", value)
    if err != nil {
        panic(err)
    }
}
// SetMembers sets the members property value. All users in the class. Nullable.
func (m *EducationClass) SetMembers(value []EducationUserable)() {
    err := m.GetBackingStore().Set("members", value)
    if err != nil {
        panic(err)
    }
}
// SetModules sets the modules property value. All modules in the class. Nullable.
func (m *EducationClass) SetModules(value []EducationModuleable)() {
    err := m.GetBackingStore().Set("modules", value)
    if err != nil {
        panic(err)
    }
}
// SetSchools sets the schools property value. All schools that this class is associated with. Nullable.
func (m *EducationClass) SetSchools(value []EducationSchoolable)() {
    err := m.GetBackingStore().Set("schools", value)
    if err != nil {
        panic(err)
    }
}
// SetTeachers sets the teachers property value. All teachers in the class. Nullable.
func (m *EducationClass) SetTeachers(value []EducationUserable)() {
    err := m.GetBackingStore().Set("teachers", value)
    if err != nil {
        panic(err)
    }
}
// SetTerm sets the term property value. Term for this class.
func (m *EducationClass) SetTerm(value EducationTermable)() {
    err := m.GetBackingStore().Set("term", value)
    if err != nil {
        panic(err)
    }
}
type EducationClassable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAssignmentCategories()([]EducationCategoryable)
    GetAssignmentDefaults()(EducationAssignmentDefaultsable)
    GetAssignments()([]EducationAssignmentable)
    GetAssignmentSettings()(EducationAssignmentSettingsable)
    GetClassCode()(*string)
    GetCourse()(EducationCourseable)
    GetCreatedBy()(IdentitySetable)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetExternalId()(*string)
    GetExternalName()(*string)
    GetExternalSource()(*EducationExternalSource)
    GetExternalSourceDetail()(*string)
    GetGrade()(*string)
    GetGroup()(Groupable)
    GetMailNickname()(*string)
    GetMembers()([]EducationUserable)
    GetModules()([]EducationModuleable)
    GetSchools()([]EducationSchoolable)
    GetTeachers()([]EducationUserable)
    GetTerm()(EducationTermable)
    SetAssignmentCategories(value []EducationCategoryable)()
    SetAssignmentDefaults(value EducationAssignmentDefaultsable)()
    SetAssignments(value []EducationAssignmentable)()
    SetAssignmentSettings(value EducationAssignmentSettingsable)()
    SetClassCode(value *string)()
    SetCourse(value EducationCourseable)()
    SetCreatedBy(value IdentitySetable)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetExternalId(value *string)()
    SetExternalName(value *string)()
    SetExternalSource(value *EducationExternalSource)()
    SetExternalSourceDetail(value *string)()
    SetGrade(value *string)()
    SetGroup(value Groupable)()
    SetMailNickname(value *string)()
    SetMembers(value []EducationUserable)()
    SetModules(value []EducationModuleable)()
    SetSchools(value []EducationSchoolable)()
    SetTeachers(value []EducationUserable)()
    SetTerm(value EducationTermable)()
}
