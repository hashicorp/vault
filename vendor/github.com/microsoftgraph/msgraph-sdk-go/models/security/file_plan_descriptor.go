package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type FilePlanDescriptor struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
}
// NewFilePlanDescriptor instantiates a new FilePlanDescriptor and sets the default values.
func NewFilePlanDescriptor()(*FilePlanDescriptor) {
    m := &FilePlanDescriptor{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateFilePlanDescriptorFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateFilePlanDescriptorFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFilePlanDescriptor(), nil
}
// GetAuthority gets the authority property value. Represents the file plan descriptor of type authority applied to a particular retention label.
// returns a FilePlanAuthorityable when successful
func (m *FilePlanDescriptor) GetAuthority()(FilePlanAuthorityable) {
    val, err := m.GetBackingStore().Get("authority")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(FilePlanAuthorityable)
    }
    return nil
}
// GetAuthorityTemplate gets the authorityTemplate property value. Specifies the underlying authority that describes the type of content to be retained and its retention schedule.
// returns a AuthorityTemplateable when successful
func (m *FilePlanDescriptor) GetAuthorityTemplate()(AuthorityTemplateable) {
    val, err := m.GetBackingStore().Get("authorityTemplate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AuthorityTemplateable)
    }
    return nil
}
// GetCategory gets the category property value. The category property
// returns a FilePlanAppliedCategoryable when successful
func (m *FilePlanDescriptor) GetCategory()(FilePlanAppliedCategoryable) {
    val, err := m.GetBackingStore().Get("category")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(FilePlanAppliedCategoryable)
    }
    return nil
}
// GetCategoryTemplate gets the categoryTemplate property value. Specifies a group of similar types of content in a particular department.
// returns a CategoryTemplateable when successful
func (m *FilePlanDescriptor) GetCategoryTemplate()(CategoryTemplateable) {
    val, err := m.GetBackingStore().Get("categoryTemplate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CategoryTemplateable)
    }
    return nil
}
// GetCitation gets the citation property value. Represents the file plan descriptor of type citation applied to a particular retention label.
// returns a FilePlanCitationable when successful
func (m *FilePlanDescriptor) GetCitation()(FilePlanCitationable) {
    val, err := m.GetBackingStore().Get("citation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(FilePlanCitationable)
    }
    return nil
}
// GetCitationTemplate gets the citationTemplate property value. The specific rule or regulation created by a jurisdiction used to determine whether certain labels and content should be retained or deleted.
// returns a CitationTemplateable when successful
func (m *FilePlanDescriptor) GetCitationTemplate()(CitationTemplateable) {
    val, err := m.GetBackingStore().Get("citationTemplate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CitationTemplateable)
    }
    return nil
}
// GetDepartment gets the department property value. Represents the file plan descriptor of type department applied to a particular retention label.
// returns a FilePlanDepartmentable when successful
func (m *FilePlanDescriptor) GetDepartment()(FilePlanDepartmentable) {
    val, err := m.GetBackingStore().Get("department")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(FilePlanDepartmentable)
    }
    return nil
}
// GetDepartmentTemplate gets the departmentTemplate property value. Specifies the  department or business unit of an organization to which a label belongs.
// returns a DepartmentTemplateable when successful
func (m *FilePlanDescriptor) GetDepartmentTemplate()(DepartmentTemplateable) {
    val, err := m.GetBackingStore().Get("departmentTemplate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DepartmentTemplateable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *FilePlanDescriptor) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["authority"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFilePlanAuthorityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAuthority(val.(FilePlanAuthorityable))
        }
        return nil
    }
    res["authorityTemplate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAuthorityTemplateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAuthorityTemplate(val.(AuthorityTemplateable))
        }
        return nil
    }
    res["category"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFilePlanAppliedCategoryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory(val.(FilePlanAppliedCategoryable))
        }
        return nil
    }
    res["categoryTemplate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCategoryTemplateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategoryTemplate(val.(CategoryTemplateable))
        }
        return nil
    }
    res["citation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFilePlanCitationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCitation(val.(FilePlanCitationable))
        }
        return nil
    }
    res["citationTemplate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCitationTemplateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCitationTemplate(val.(CitationTemplateable))
        }
        return nil
    }
    res["department"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFilePlanDepartmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDepartment(val.(FilePlanDepartmentable))
        }
        return nil
    }
    res["departmentTemplate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDepartmentTemplateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDepartmentTemplate(val.(DepartmentTemplateable))
        }
        return nil
    }
    res["filePlanReference"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFilePlanReferenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFilePlanReference(val.(FilePlanReferenceable))
        }
        return nil
    }
    res["filePlanReferenceTemplate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFilePlanReferenceTemplateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFilePlanReferenceTemplate(val.(FilePlanReferenceTemplateable))
        }
        return nil
    }
    return res
}
// GetFilePlanReference gets the filePlanReference property value. Represents the file plan descriptor of type filePlanReference applied to a particular retention label.
// returns a FilePlanReferenceable when successful
func (m *FilePlanDescriptor) GetFilePlanReference()(FilePlanReferenceable) {
    val, err := m.GetBackingStore().Get("filePlanReference")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(FilePlanReferenceable)
    }
    return nil
}
// GetFilePlanReferenceTemplate gets the filePlanReferenceTemplate property value. Specifies a unique alpha-numeric identifier for an organization’s retention schedule.
// returns a FilePlanReferenceTemplateable when successful
func (m *FilePlanDescriptor) GetFilePlanReferenceTemplate()(FilePlanReferenceTemplateable) {
    val, err := m.GetBackingStore().Get("filePlanReferenceTemplate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(FilePlanReferenceTemplateable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *FilePlanDescriptor) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("authority", m.GetAuthority())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("authorityTemplate", m.GetAuthorityTemplate())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("category", m.GetCategory())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("categoryTemplate", m.GetCategoryTemplate())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("citation", m.GetCitation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("citationTemplate", m.GetCitationTemplate())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("department", m.GetDepartment())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("departmentTemplate", m.GetDepartmentTemplate())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("filePlanReference", m.GetFilePlanReference())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("filePlanReferenceTemplate", m.GetFilePlanReferenceTemplate())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAuthority sets the authority property value. Represents the file plan descriptor of type authority applied to a particular retention label.
func (m *FilePlanDescriptor) SetAuthority(value FilePlanAuthorityable)() {
    err := m.GetBackingStore().Set("authority", value)
    if err != nil {
        panic(err)
    }
}
// SetAuthorityTemplate sets the authorityTemplate property value. Specifies the underlying authority that describes the type of content to be retained and its retention schedule.
func (m *FilePlanDescriptor) SetAuthorityTemplate(value AuthorityTemplateable)() {
    err := m.GetBackingStore().Set("authorityTemplate", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory sets the category property value. The category property
func (m *FilePlanDescriptor) SetCategory(value FilePlanAppliedCategoryable)() {
    err := m.GetBackingStore().Set("category", value)
    if err != nil {
        panic(err)
    }
}
// SetCategoryTemplate sets the categoryTemplate property value. Specifies a group of similar types of content in a particular department.
func (m *FilePlanDescriptor) SetCategoryTemplate(value CategoryTemplateable)() {
    err := m.GetBackingStore().Set("categoryTemplate", value)
    if err != nil {
        panic(err)
    }
}
// SetCitation sets the citation property value. Represents the file plan descriptor of type citation applied to a particular retention label.
func (m *FilePlanDescriptor) SetCitation(value FilePlanCitationable)() {
    err := m.GetBackingStore().Set("citation", value)
    if err != nil {
        panic(err)
    }
}
// SetCitationTemplate sets the citationTemplate property value. The specific rule or regulation created by a jurisdiction used to determine whether certain labels and content should be retained or deleted.
func (m *FilePlanDescriptor) SetCitationTemplate(value CitationTemplateable)() {
    err := m.GetBackingStore().Set("citationTemplate", value)
    if err != nil {
        panic(err)
    }
}
// SetDepartment sets the department property value. Represents the file plan descriptor of type department applied to a particular retention label.
func (m *FilePlanDescriptor) SetDepartment(value FilePlanDepartmentable)() {
    err := m.GetBackingStore().Set("department", value)
    if err != nil {
        panic(err)
    }
}
// SetDepartmentTemplate sets the departmentTemplate property value. Specifies the  department or business unit of an organization to which a label belongs.
func (m *FilePlanDescriptor) SetDepartmentTemplate(value DepartmentTemplateable)() {
    err := m.GetBackingStore().Set("departmentTemplate", value)
    if err != nil {
        panic(err)
    }
}
// SetFilePlanReference sets the filePlanReference property value. Represents the file plan descriptor of type filePlanReference applied to a particular retention label.
func (m *FilePlanDescriptor) SetFilePlanReference(value FilePlanReferenceable)() {
    err := m.GetBackingStore().Set("filePlanReference", value)
    if err != nil {
        panic(err)
    }
}
// SetFilePlanReferenceTemplate sets the filePlanReferenceTemplate property value. Specifies a unique alpha-numeric identifier for an organization’s retention schedule.
func (m *FilePlanDescriptor) SetFilePlanReferenceTemplate(value FilePlanReferenceTemplateable)() {
    err := m.GetBackingStore().Set("filePlanReferenceTemplate", value)
    if err != nil {
        panic(err)
    }
}
type FilePlanDescriptorable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAuthority()(FilePlanAuthorityable)
    GetAuthorityTemplate()(AuthorityTemplateable)
    GetCategory()(FilePlanAppliedCategoryable)
    GetCategoryTemplate()(CategoryTemplateable)
    GetCitation()(FilePlanCitationable)
    GetCitationTemplate()(CitationTemplateable)
    GetDepartment()(FilePlanDepartmentable)
    GetDepartmentTemplate()(DepartmentTemplateable)
    GetFilePlanReference()(FilePlanReferenceable)
    GetFilePlanReferenceTemplate()(FilePlanReferenceTemplateable)
    SetAuthority(value FilePlanAuthorityable)()
    SetAuthorityTemplate(value AuthorityTemplateable)()
    SetCategory(value FilePlanAppliedCategoryable)()
    SetCategoryTemplate(value CategoryTemplateable)()
    SetCitation(value FilePlanCitationable)()
    SetCitationTemplate(value CitationTemplateable)()
    SetDepartment(value FilePlanDepartmentable)()
    SetDepartmentTemplate(value DepartmentTemplateable)()
    SetFilePlanReference(value FilePlanReferenceable)()
    SetFilePlanReferenceTemplate(value FilePlanReferenceTemplateable)()
}
