package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type MailSearchFolder struct {
    MailFolder
}
// NewMailSearchFolder instantiates a new MailSearchFolder and sets the default values.
func NewMailSearchFolder()(*MailSearchFolder) {
    m := &MailSearchFolder{
        MailFolder: *NewMailFolder(),
    }
    odataTypeValue := "#microsoft.graph.mailSearchFolder"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateMailSearchFolderFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMailSearchFolderFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMailSearchFolder(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MailSearchFolder) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MailFolder.GetFieldDeserializers()
    res["filterQuery"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFilterQuery(val)
        }
        return nil
    }
    res["includeNestedFolders"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIncludeNestedFolders(val)
        }
        return nil
    }
    res["isSupported"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSupported(val)
        }
        return nil
    }
    res["sourceFolderIds"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSourceFolderIds(res)
        }
        return nil
    }
    return res
}
// GetFilterQuery gets the filterQuery property value. The OData query to filter the messages.
// returns a *string when successful
func (m *MailSearchFolder) GetFilterQuery()(*string) {
    val, err := m.GetBackingStore().Get("filterQuery")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIncludeNestedFolders gets the includeNestedFolders property value. Indicates how the mailbox folder hierarchy should be traversed in the search. true means that a deep search should be done to include child folders in the hierarchy of each folder explicitly specified in sourceFolderIds. false means a shallow search of only each of the folders explicitly specified in sourceFolderIds.
// returns a *bool when successful
func (m *MailSearchFolder) GetIncludeNestedFolders()(*bool) {
    val, err := m.GetBackingStore().Get("includeNestedFolders")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsSupported gets the isSupported property value. Indicates whether a search folder is editable using REST APIs.
// returns a *bool when successful
func (m *MailSearchFolder) GetIsSupported()(*bool) {
    val, err := m.GetBackingStore().Get("isSupported")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSourceFolderIds gets the sourceFolderIds property value. The mailbox folders that should be mined.
// returns a []string when successful
func (m *MailSearchFolder) GetSourceFolderIds()([]string) {
    val, err := m.GetBackingStore().Get("sourceFolderIds")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MailSearchFolder) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MailFolder.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("filterQuery", m.GetFilterQuery())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("includeNestedFolders", m.GetIncludeNestedFolders())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isSupported", m.GetIsSupported())
        if err != nil {
            return err
        }
    }
    if m.GetSourceFolderIds() != nil {
        err = writer.WriteCollectionOfStringValues("sourceFolderIds", m.GetSourceFolderIds())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFilterQuery sets the filterQuery property value. The OData query to filter the messages.
func (m *MailSearchFolder) SetFilterQuery(value *string)() {
    err := m.GetBackingStore().Set("filterQuery", value)
    if err != nil {
        panic(err)
    }
}
// SetIncludeNestedFolders sets the includeNestedFolders property value. Indicates how the mailbox folder hierarchy should be traversed in the search. true means that a deep search should be done to include child folders in the hierarchy of each folder explicitly specified in sourceFolderIds. false means a shallow search of only each of the folders explicitly specified in sourceFolderIds.
func (m *MailSearchFolder) SetIncludeNestedFolders(value *bool)() {
    err := m.GetBackingStore().Set("includeNestedFolders", value)
    if err != nil {
        panic(err)
    }
}
// SetIsSupported sets the isSupported property value. Indicates whether a search folder is editable using REST APIs.
func (m *MailSearchFolder) SetIsSupported(value *bool)() {
    err := m.GetBackingStore().Set("isSupported", value)
    if err != nil {
        panic(err)
    }
}
// SetSourceFolderIds sets the sourceFolderIds property value. The mailbox folders that should be mined.
func (m *MailSearchFolder) SetSourceFolderIds(value []string)() {
    err := m.GetBackingStore().Set("sourceFolderIds", value)
    if err != nil {
        panic(err)
    }
}
type MailSearchFolderable interface {
    MailFolderable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetFilterQuery()(*string)
    GetIncludeNestedFolders()(*bool)
    GetIsSupported()(*bool)
    GetSourceFolderIds()([]string)
    SetFilterQuery(value *string)()
    SetIncludeNestedFolders(value *bool)()
    SetIsSupported(value *bool)()
    SetSourceFolderIds(value []string)()
}
