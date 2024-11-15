package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Onenote struct {
    Entity
}
// NewOnenote instantiates a new Onenote and sets the default values.
func NewOnenote()(*Onenote) {
    m := &Onenote{
        Entity: *NewEntity(),
    }
    return m
}
// CreateOnenoteFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOnenoteFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOnenote(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Onenote) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["notebooks"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateNotebookFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Notebookable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Notebookable)
                }
            }
            m.SetNotebooks(res)
        }
        return nil
    }
    res["operations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateOnenoteOperationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]OnenoteOperationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(OnenoteOperationable)
                }
            }
            m.SetOperations(res)
        }
        return nil
    }
    res["pages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateOnenotePageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]OnenotePageable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(OnenotePageable)
                }
            }
            m.SetPages(res)
        }
        return nil
    }
    res["resources"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateOnenoteResourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]OnenoteResourceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(OnenoteResourceable)
                }
            }
            m.SetResources(res)
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
    return res
}
// GetNotebooks gets the notebooks property value. The collection of OneNote notebooks that are owned by the user or group. Read-only. Nullable.
// returns a []Notebookable when successful
func (m *Onenote) GetNotebooks()([]Notebookable) {
    val, err := m.GetBackingStore().Get("notebooks")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Notebookable)
    }
    return nil
}
// GetOperations gets the operations property value. The status of OneNote operations. Getting an operations collection isn't supported, but you can get the status of long-running operations if the Operation-Location header is returned in the response. Read-only. Nullable.
// returns a []OnenoteOperationable when successful
func (m *Onenote) GetOperations()([]OnenoteOperationable) {
    val, err := m.GetBackingStore().Get("operations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]OnenoteOperationable)
    }
    return nil
}
// GetPages gets the pages property value. The pages in all OneNote notebooks that are owned by the user or group.  Read-only. Nullable.
// returns a []OnenotePageable when successful
func (m *Onenote) GetPages()([]OnenotePageable) {
    val, err := m.GetBackingStore().Get("pages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]OnenotePageable)
    }
    return nil
}
// GetResources gets the resources property value. The image and other file resources in OneNote pages. Getting a resources collection isn't supported, but you can get the binary content of a specific resource. Read-only. Nullable.
// returns a []OnenoteResourceable when successful
func (m *Onenote) GetResources()([]OnenoteResourceable) {
    val, err := m.GetBackingStore().Get("resources")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]OnenoteResourceable)
    }
    return nil
}
// GetSectionGroups gets the sectionGroups property value. The section groups in all OneNote notebooks that are owned by the user or group.  Read-only. Nullable.
// returns a []SectionGroupable when successful
func (m *Onenote) GetSectionGroups()([]SectionGroupable) {
    val, err := m.GetBackingStore().Get("sectionGroups")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SectionGroupable)
    }
    return nil
}
// GetSections gets the sections property value. The sections in all OneNote notebooks that are owned by the user or group.  Read-only. Nullable.
// returns a []OnenoteSectionable when successful
func (m *Onenote) GetSections()([]OnenoteSectionable) {
    val, err := m.GetBackingStore().Get("sections")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]OnenoteSectionable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Onenote) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetNotebooks() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetNotebooks()))
        for i, v := range m.GetNotebooks() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("notebooks", cast)
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
    if m.GetPages() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPages()))
        for i, v := range m.GetPages() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("pages", cast)
        if err != nil {
            return err
        }
    }
    if m.GetResources() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetResources()))
        for i, v := range m.GetResources() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("resources", cast)
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
    return nil
}
// SetNotebooks sets the notebooks property value. The collection of OneNote notebooks that are owned by the user or group. Read-only. Nullable.
func (m *Onenote) SetNotebooks(value []Notebookable)() {
    err := m.GetBackingStore().Set("notebooks", value)
    if err != nil {
        panic(err)
    }
}
// SetOperations sets the operations property value. The status of OneNote operations. Getting an operations collection isn't supported, but you can get the status of long-running operations if the Operation-Location header is returned in the response. Read-only. Nullable.
func (m *Onenote) SetOperations(value []OnenoteOperationable)() {
    err := m.GetBackingStore().Set("operations", value)
    if err != nil {
        panic(err)
    }
}
// SetPages sets the pages property value. The pages in all OneNote notebooks that are owned by the user or group.  Read-only. Nullable.
func (m *Onenote) SetPages(value []OnenotePageable)() {
    err := m.GetBackingStore().Set("pages", value)
    if err != nil {
        panic(err)
    }
}
// SetResources sets the resources property value. The image and other file resources in OneNote pages. Getting a resources collection isn't supported, but you can get the binary content of a specific resource. Read-only. Nullable.
func (m *Onenote) SetResources(value []OnenoteResourceable)() {
    err := m.GetBackingStore().Set("resources", value)
    if err != nil {
        panic(err)
    }
}
// SetSectionGroups sets the sectionGroups property value. The section groups in all OneNote notebooks that are owned by the user or group.  Read-only. Nullable.
func (m *Onenote) SetSectionGroups(value []SectionGroupable)() {
    err := m.GetBackingStore().Set("sectionGroups", value)
    if err != nil {
        panic(err)
    }
}
// SetSections sets the sections property value. The sections in all OneNote notebooks that are owned by the user or group.  Read-only. Nullable.
func (m *Onenote) SetSections(value []OnenoteSectionable)() {
    err := m.GetBackingStore().Set("sections", value)
    if err != nil {
        panic(err)
    }
}
type Onenoteable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetNotebooks()([]Notebookable)
    GetOperations()([]OnenoteOperationable)
    GetPages()([]OnenotePageable)
    GetResources()([]OnenoteResourceable)
    GetSectionGroups()([]SectionGroupable)
    GetSections()([]OnenoteSectionable)
    SetNotebooks(value []Notebookable)()
    SetOperations(value []OnenoteOperationable)()
    SetPages(value []OnenotePageable)()
    SetResources(value []OnenoteResourceable)()
    SetSectionGroups(value []SectionGroupable)()
    SetSections(value []OnenoteSectionable)()
}
