package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type BaseSitePage struct {
    BaseItem
}
// NewBaseSitePage instantiates a new BaseSitePage and sets the default values.
func NewBaseSitePage()(*BaseSitePage) {
    m := &BaseSitePage{
        BaseItem: *NewBaseItem(),
    }
    odataTypeValue := "#microsoft.graph.baseSitePage"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateBaseSitePageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateBaseSitePageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.sitePage":
                        return NewSitePage(), nil
                }
            }
        }
    }
    return NewBaseSitePage(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *BaseSitePage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseItem.GetFieldDeserializers()
    res["pageLayout"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePageLayoutType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPageLayout(val.(*PageLayoutType))
        }
        return nil
    }
    res["publishingState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePublicationFacetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPublishingState(val.(PublicationFacetable))
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
    return res
}
// GetPageLayout gets the pageLayout property value. The name of the page layout of the page. The possible values are: microsoftReserved, article, home, unknownFutureValue.
// returns a *PageLayoutType when successful
func (m *BaseSitePage) GetPageLayout()(*PageLayoutType) {
    val, err := m.GetBackingStore().Get("pageLayout")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PageLayoutType)
    }
    return nil
}
// GetPublishingState gets the publishingState property value. The publishing status and the MM.mm version of the page.
// returns a PublicationFacetable when successful
func (m *BaseSitePage) GetPublishingState()(PublicationFacetable) {
    val, err := m.GetBackingStore().Get("publishingState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PublicationFacetable)
    }
    return nil
}
// GetTitle gets the title property value. Title of the sitePage.
// returns a *string when successful
func (m *BaseSitePage) GetTitle()(*string) {
    val, err := m.GetBackingStore().Get("title")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *BaseSitePage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseItem.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetPageLayout() != nil {
        cast := (*m.GetPageLayout()).String()
        err = writer.WriteStringValue("pageLayout", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("publishingState", m.GetPublishingState())
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
    return nil
}
// SetPageLayout sets the pageLayout property value. The name of the page layout of the page. The possible values are: microsoftReserved, article, home, unknownFutureValue.
func (m *BaseSitePage) SetPageLayout(value *PageLayoutType)() {
    err := m.GetBackingStore().Set("pageLayout", value)
    if err != nil {
        panic(err)
    }
}
// SetPublishingState sets the publishingState property value. The publishing status and the MM.mm version of the page.
func (m *BaseSitePage) SetPublishingState(value PublicationFacetable)() {
    err := m.GetBackingStore().Set("publishingState", value)
    if err != nil {
        panic(err)
    }
}
// SetTitle sets the title property value. Title of the sitePage.
func (m *BaseSitePage) SetTitle(value *string)() {
    err := m.GetBackingStore().Set("title", value)
    if err != nil {
        panic(err)
    }
}
type BaseSitePageable interface {
    BaseItemable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetPageLayout()(*PageLayoutType)
    GetPublishingState()(PublicationFacetable)
    GetTitle()(*string)
    SetPageLayout(value *PageLayoutType)()
    SetPublishingState(value PublicationFacetable)()
    SetTitle(value *string)()
}
