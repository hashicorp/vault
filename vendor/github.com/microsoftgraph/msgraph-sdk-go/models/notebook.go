package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Notebook struct {
    OnenoteEntityHierarchyModel
}
// NewNotebook instantiates a new Notebook and sets the default values.
func NewNotebook()(*Notebook) {
    m := &Notebook{
        OnenoteEntityHierarchyModel: *NewOnenoteEntityHierarchyModel(),
    }
    odataTypeValue := "#microsoft.graph.notebook"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateNotebookFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateNotebookFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewNotebook(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Notebook) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.OnenoteEntityHierarchyModel.GetFieldDeserializers()
    res["isDefault"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsDefault(val)
        }
        return nil
    }
    res["isShared"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsShared(val)
        }
        return nil
    }
    res["links"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateNotebookLinksFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLinks(val.(NotebookLinksable))
        }
        return nil
    }
    res["sectionGroups"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSectionGroupFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SectionGroupable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SectionGroupable)
                }
            }
            m.SetSectionGroups(res)
        }
        return nil
    }
    res["sectionGroupsUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSectionGroupsUrl(val)
        }
        return nil
    }
    res["sections"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateOnenoteSectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]OnenoteSectionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(OnenoteSectionable)
                }
            }
            m.SetSections(res)
        }
        return nil
    }
    res["sectionsUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSectionsUrl(val)
        }
        return nil
    }
    res["userRole"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseOnenoteUserRole)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserRole(val.(*OnenoteUserRole))
        }
        return nil
    }
    return res
}
// GetIsDefault gets the isDefault property value. Indicates whether this is the user's default notebook. Read-only.
// returns a *bool when successful
func (m *Notebook) GetIsDefault()(*bool) {
    val, err := m.GetBackingStore().Get("isDefault")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsShared gets the isShared property value. Indicates whether the notebook is shared. If true, the contents of the notebook can be seen by people other than the owner. Read-only.
// returns a *bool when successful
func (m *Notebook) GetIsShared()(*bool) {
    val, err := m.GetBackingStore().Get("isShared")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLinks gets the links property value. Links for opening the notebook. The oneNoteClientURL link opens the notebook in the OneNote native client if it's installed. The oneNoteWebURL link opens the notebook in OneNote on the web.
// returns a NotebookLinksable when successful
func (m *Notebook) GetLinks()(NotebookLinksable) {
    val, err := m.GetBackingStore().Get("links")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(NotebookLinksable)
    }
    return nil
}
// GetSectionGroups gets the sectionGroups property value. The section groups in the notebook. Read-only. Nullable.
// returns a []SectionGroupable when successful
func (m *Notebook) GetSectionGroups()([]SectionGroupable) {
    val, err := m.GetBackingStore().Get("sectionGroups")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SectionGroupable)
    }
    return nil
}
// GetSectionGroupsUrl gets the sectionGroupsUrl property value. The URL for the sectionGroups navigation property, which returns all the section groups in the notebook. Read-only.
// returns a *string when successful
func (m *Notebook) GetSectionGroupsUrl()(*string) {
    val, err := m.GetBackingStore().Get("sectionGroupsUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSections gets the sections property value. The sections in the notebook. Read-only. Nullable.
// returns a []OnenoteSectionable when successful
func (m *Notebook) GetSections()([]OnenoteSectionable) {
    val, err := m.GetBackingStore().Get("sections")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]OnenoteSectionable)
    }
    return nil
}
// GetSectionsUrl gets the sectionsUrl property value. The URL for the sections navigation property, which returns all the sections in the notebook. Read-only.
// returns a *string when successful
func (m *Notebook) GetSectionsUrl()(*string) {
    val, err := m.GetBackingStore().Get("sectionsUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserRole gets the userRole property value. Possible values are: Owner, Contributor, Reader, None. Owner represents owner-level access to the notebook. Contributor represents read/write access to the notebook. Reader represents read-only access to the notebook. Read-only.
// returns a *OnenoteUserRole when successful
func (m *Notebook) GetUserRole()(*OnenoteUserRole) {
    val, err := m.GetBackingStore().Get("userRole")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*OnenoteUserRole)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Notebook) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.OnenoteEntityHierarchyModel.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("isDefault", m.GetIsDefault())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isShared", m.GetIsShared())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("links", m.GetLinks())
        if err != nil {
            return err
        }
    }
    if m.GetSectionGroups() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSectionGroups()))
        for i, v := range m.GetSectionGroups() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("sectionGroups", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("sectionGroupsUrl", m.GetSectionGroupsUrl())
        if err != nil {
            return err
        }
    }
    if m.GetSections() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSections()))
        for i, v := range m.GetSections() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("sections", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("sectionsUrl", m.GetSectionsUrl())
        if err != nil {
            return err
        }
    }
    if m.GetUserRole() != nil {
        cast := (*m.GetUserRole()).String()
        err = writer.WriteStringValue("userRole", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetIsDefault sets the isDefault property value. Indicates whether this is the user's default notebook. Read-only.
func (m *Notebook) SetIsDefault(value *bool)() {
    err := m.GetBackingStore().Set("isDefault", value)
    if err != nil {
        panic(err)
    }
}
// SetIsShared sets the isShared property value. Indicates whether the notebook is shared. If true, the contents of the notebook can be seen by people other than the owner. Read-only.
func (m *Notebook) SetIsShared(value *bool)() {
    err := m.GetBackingStore().Set("isShared", value)
    if err != nil {
        panic(err)
    }
}
// SetLinks sets the links property value. Links for opening the notebook. The oneNoteClientURL link opens the notebook in the OneNote native client if it's installed. The oneNoteWebURL link opens the notebook in OneNote on the web.
func (m *Notebook) SetLinks(value NotebookLinksable)() {
    err := m.GetBackingStore().Set("links", value)
    if err != nil {
        panic(err)
    }
}
// SetSectionGroups sets the sectionGroups property value. The section groups in the notebook. Read-only. Nullable.
func (m *Notebook) SetSectionGroups(value []SectionGroupable)() {
    err := m.GetBackingStore().Set("sectionGroups", value)
    if err != nil {
        panic(err)
    }
}
// SetSectionGroupsUrl sets the sectionGroupsUrl property value. The URL for the sectionGroups navigation property, which returns all the section groups in the notebook. Read-only.
func (m *Notebook) SetSectionGroupsUrl(value *string)() {
    err := m.GetBackingStore().Set("sectionGroupsUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetSections sets the sections property value. The sections in the notebook. Read-only. Nullable.
func (m *Notebook) SetSections(value []OnenoteSectionable)() {
    err := m.GetBackingStore().Set("sections", value)
    if err != nil {
        panic(err)
    }
}
// SetSectionsUrl sets the sectionsUrl property value. The URL for the sections navigation property, which returns all the sections in the notebook. Read-only.
func (m *Notebook) SetSectionsUrl(value *string)() {
    err := m.GetBackingStore().Set("sectionsUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetUserRole sets the userRole property value. Possible values are: Owner, Contributor, Reader, None. Owner represents owner-level access to the notebook. Contributor represents read/write access to the notebook. Reader represents read-only access to the notebook. Read-only.
func (m *Notebook) SetUserRole(value *OnenoteUserRole)() {
    err := m.GetBackingStore().Set("userRole", value)
    if err != nil {
        panic(err)
    }
}
type Notebookable interface {
    OnenoteEntityHierarchyModelable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetIsDefault()(*bool)
    GetIsShared()(*bool)
    GetLinks()(NotebookLinksable)
    GetSectionGroups()([]SectionGroupable)
    GetSectionGroupsUrl()(*string)
    GetSections()([]OnenoteSectionable)
    GetSectionsUrl()(*string)
    GetUserRole()(*OnenoteUserRole)
    SetIsDefault(value *bool)()
    SetIsShared(value *bool)()
    SetLinks(value NotebookLinksable)()
    SetSectionGroups(value []SectionGroupable)()
    SetSectionGroupsUrl(value *string)()
    SetSections(value []OnenoteSectionable)()
    SetSectionsUrl(value *string)()
    SetUserRole(value *OnenoteUserRole)()
}
