package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type OnenotePage struct {
    OnenoteEntitySchemaObjectModel
}
// NewOnenotePage instantiates a new OnenotePage and sets the default values.
func NewOnenotePage()(*OnenotePage) {
    m := &OnenotePage{
        OnenoteEntitySchemaObjectModel: *NewOnenoteEntitySchemaObjectModel(),
    }
    odataTypeValue := "#microsoft.graph.onenotePage"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateOnenotePageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOnenotePageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOnenotePage(), nil
}
// GetContent gets the content property value. The page's HTML content.
// returns a []byte when successful
func (m *OnenotePage) GetContent()([]byte) {
    val, err := m.GetBackingStore().Get("content")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetContentUrl gets the contentUrl property value. The URL for the page's HTML content.  Read-only.
// returns a *string when successful
func (m *OnenotePage) GetContentUrl()(*string) {
    val, err := m.GetBackingStore().Get("contentUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCreatedByAppId gets the createdByAppId property value. The unique identifier of the application that created the page. Read-only.
// returns a *string when successful
func (m *OnenotePage) GetCreatedByAppId()(*string) {
    val, err := m.GetBackingStore().Get("createdByAppId")
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
func (m *OnenotePage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.OnenoteEntitySchemaObjectModel.GetFieldDeserializers()
    res["content"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContent(val)
        }
        return nil
    }
    res["contentUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContentUrl(val)
        }
        return nil
    }
    res["createdByAppId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedByAppId(val)
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
    res["level"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLevel(val)
        }
        return nil
    }
    res["links"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePageLinksFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLinks(val.(PageLinksable))
        }
        return nil
    }
    res["order"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOrder(val)
        }
        return nil
    }
    res["parentNotebook"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateNotebookFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParentNotebook(val.(Notebookable))
        }
        return nil
    }
    res["parentSection"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOnenoteSectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParentSection(val.(OnenoteSectionable))
        }
        return nil
    }
    res["title"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTitle(val)
        }
        return nil
    }
    res["userTags"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetUserTags(res)
        }
        return nil
    }
    return res
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. The date and time when the page was last modified. The timestamp represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
// returns a *Time when successful
func (m *OnenotePage) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetLevel gets the level property value. The indentation level of the page. Read-only.
// returns a *int32 when successful
func (m *OnenotePage) GetLevel()(*int32) {
    val, err := m.GetBackingStore().Get("level")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetLinks gets the links property value. Links for opening the page. The oneNoteClientURL link opens the page in the OneNote native client if it 's installed. The oneNoteWebUrl link opens the page in OneNote on the web. Read-only.
// returns a PageLinksable when successful
func (m *OnenotePage) GetLinks()(PageLinksable) {
    val, err := m.GetBackingStore().Get("links")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PageLinksable)
    }
    return nil
}
// GetOrder gets the order property value. The order of the page within its parent section. Read-only.
// returns a *int32 when successful
func (m *OnenotePage) GetOrder()(*int32) {
    val, err := m.GetBackingStore().Get("order")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetParentNotebook gets the parentNotebook property value. The notebook that contains the page.  Read-only.
// returns a Notebookable when successful
func (m *OnenotePage) GetParentNotebook()(Notebookable) {
    val, err := m.GetBackingStore().Get("parentNotebook")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Notebookable)
    }
    return nil
}
// GetParentSection gets the parentSection property value. The section that contains the page. Read-only.
// returns a OnenoteSectionable when successful
func (m *OnenotePage) GetParentSection()(OnenoteSectionable) {
    val, err := m.GetBackingStore().Get("parentSection")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(OnenoteSectionable)
    }
    return nil
}
// GetTitle gets the title property value. The title of the page.
// returns a *string when successful
func (m *OnenotePage) GetTitle()(*string) {
    val, err := m.GetBackingStore().Get("title")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserTags gets the userTags property value. The userTags property
// returns a []string when successful
func (m *OnenotePage) GetUserTags()([]string) {
    val, err := m.GetBackingStore().Get("userTags")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *OnenotePage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.OnenoteEntitySchemaObjectModel.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteByteArrayValue("content", m.GetContent())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("contentUrl", m.GetContentUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("createdByAppId", m.GetCreatedByAppId())
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
    {
        err = writer.WriteInt32Value("level", m.GetLevel())
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
    {
        err = writer.WriteInt32Value("order", m.GetOrder())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("parentNotebook", m.GetParentNotebook())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("parentSection", m.GetParentSection())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("title", m.GetTitle())
        if err != nil {
            return err
        }
    }
    if m.GetUserTags() != nil {
        err = writer.WriteCollectionOfStringValues("userTags", m.GetUserTags())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetContent sets the content property value. The page's HTML content.
func (m *OnenotePage) SetContent(value []byte)() {
    err := m.GetBackingStore().Set("content", value)
    if err != nil {
        panic(err)
    }
}
// SetContentUrl sets the contentUrl property value. The URL for the page's HTML content.  Read-only.
func (m *OnenotePage) SetContentUrl(value *string)() {
    err := m.GetBackingStore().Set("contentUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedByAppId sets the createdByAppId property value. The unique identifier of the application that created the page. Read-only.
func (m *OnenotePage) SetCreatedByAppId(value *string)() {
    err := m.GetBackingStore().Set("createdByAppId", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. The date and time when the page was last modified. The timestamp represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
func (m *OnenotePage) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLevel sets the level property value. The indentation level of the page. Read-only.
func (m *OnenotePage) SetLevel(value *int32)() {
    err := m.GetBackingStore().Set("level", value)
    if err != nil {
        panic(err)
    }
}
// SetLinks sets the links property value. Links for opening the page. The oneNoteClientURL link opens the page in the OneNote native client if it 's installed. The oneNoteWebUrl link opens the page in OneNote on the web. Read-only.
func (m *OnenotePage) SetLinks(value PageLinksable)() {
    err := m.GetBackingStore().Set("links", value)
    if err != nil {
        panic(err)
    }
}
// SetOrder sets the order property value. The order of the page within its parent section. Read-only.
func (m *OnenotePage) SetOrder(value *int32)() {
    err := m.GetBackingStore().Set("order", value)
    if err != nil {
        panic(err)
    }
}
// SetParentNotebook sets the parentNotebook property value. The notebook that contains the page.  Read-only.
func (m *OnenotePage) SetParentNotebook(value Notebookable)() {
    err := m.GetBackingStore().Set("parentNotebook", value)
    if err != nil {
        panic(err)
    }
}
// SetParentSection sets the parentSection property value. The section that contains the page. Read-only.
func (m *OnenotePage) SetParentSection(value OnenoteSectionable)() {
    err := m.GetBackingStore().Set("parentSection", value)
    if err != nil {
        panic(err)
    }
}
// SetTitle sets the title property value. The title of the page.
func (m *OnenotePage) SetTitle(value *string)() {
    err := m.GetBackingStore().Set("title", value)
    if err != nil {
        panic(err)
    }
}
// SetUserTags sets the userTags property value. The userTags property
func (m *OnenotePage) SetUserTags(value []string)() {
    err := m.GetBackingStore().Set("userTags", value)
    if err != nil {
        panic(err)
    }
}
type OnenotePageable interface {
    OnenoteEntitySchemaObjectModelable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetContent()([]byte)
    GetContentUrl()(*string)
    GetCreatedByAppId()(*string)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLevel()(*int32)
    GetLinks()(PageLinksable)
    GetOrder()(*int32)
    GetParentNotebook()(Notebookable)
    GetParentSection()(OnenoteSectionable)
    GetTitle()(*string)
    GetUserTags()([]string)
    SetContent(value []byte)()
    SetContentUrl(value *string)()
    SetCreatedByAppId(value *string)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLevel(value *int32)()
    SetLinks(value PageLinksable)()
    SetOrder(value *int32)()
    SetParentNotebook(value Notebookable)()
    SetParentSection(value OnenoteSectionable)()
    SetTitle(value *string)()
    SetUserTags(value []string)()
}
