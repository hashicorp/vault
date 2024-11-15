package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SitePage struct {
    BaseSitePage
}
// NewSitePage instantiates a new SitePage and sets the default values.
func NewSitePage()(*SitePage) {
    m := &SitePage{
        BaseSitePage: *NewBaseSitePage(),
    }
    return m
}
// CreateSitePageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSitePageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSitePage(), nil
}
// GetCanvasLayout gets the canvasLayout property value. Indicates the layout of the content in a given SharePoint page, including horizontal sections and vertical sections.
// returns a CanvasLayoutable when successful
func (m *SitePage) GetCanvasLayout()(CanvasLayoutable) {
    val, err := m.GetBackingStore().Get("canvasLayout")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CanvasLayoutable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SitePage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseSitePage.GetFieldDeserializers()
    res["canvasLayout"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCanvasLayoutFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCanvasLayout(val.(CanvasLayoutable))
        }
        return nil
    }
    res["promotionKind"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePagePromotionType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPromotionKind(val.(*PagePromotionType))
        }
        return nil
    }
    res["reactions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateReactionsFacetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReactions(val.(ReactionsFacetable))
        }
        return nil
    }
    res["showComments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowComments(val)
        }
        return nil
    }
    res["showRecommendedPages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowRecommendedPages(val)
        }
        return nil
    }
    res["thumbnailWebUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetThumbnailWebUrl(val)
        }
        return nil
    }
    res["titleArea"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTitleAreaFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTitleArea(val.(TitleAreaable))
        }
        return nil
    }
    res["webParts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWebPartFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WebPartable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WebPartable)
                }
            }
            m.SetWebParts(res)
        }
        return nil
    }
    return res
}
// GetPromotionKind gets the promotionKind property value. Indicates the promotion kind of the sitePage. The possible values are: microsoftReserved, page, newsPost, unknownFutureValue.
// returns a *PagePromotionType when successful
func (m *SitePage) GetPromotionKind()(*PagePromotionType) {
    val, err := m.GetBackingStore().Get("promotionKind")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PagePromotionType)
    }
    return nil
}
// GetReactions gets the reactions property value. Reactions information for the page.
// returns a ReactionsFacetable when successful
func (m *SitePage) GetReactions()(ReactionsFacetable) {
    val, err := m.GetBackingStore().Get("reactions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ReactionsFacetable)
    }
    return nil
}
// GetShowComments gets the showComments property value. Determines whether or not to show comments at the bottom of the page.
// returns a *bool when successful
func (m *SitePage) GetShowComments()(*bool) {
    val, err := m.GetBackingStore().Get("showComments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetShowRecommendedPages gets the showRecommendedPages property value. Determines whether or not to show recommended pages at the bottom of the page.
// returns a *bool when successful
func (m *SitePage) GetShowRecommendedPages()(*bool) {
    val, err := m.GetBackingStore().Get("showRecommendedPages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetThumbnailWebUrl gets the thumbnailWebUrl property value. Url of the sitePage's thumbnail image
// returns a *string when successful
func (m *SitePage) GetThumbnailWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("thumbnailWebUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTitleArea gets the titleArea property value. Title area on the SharePoint page.
// returns a TitleAreaable when successful
func (m *SitePage) GetTitleArea()(TitleAreaable) {
    val, err := m.GetBackingStore().Get("titleArea")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TitleAreaable)
    }
    return nil
}
// GetWebParts gets the webParts property value. Collection of webparts on the SharePoint page.
// returns a []WebPartable when successful
func (m *SitePage) GetWebParts()([]WebPartable) {
    val, err := m.GetBackingStore().Get("webParts")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WebPartable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SitePage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseSitePage.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("canvasLayout", m.GetCanvasLayout())
        if err != nil {
            return err
        }
    }
    if m.GetPromotionKind() != nil {
        cast := (*m.GetPromotionKind()).String()
        err = writer.WriteStringValue("promotionKind", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("reactions", m.GetReactions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("showComments", m.GetShowComments())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("showRecommendedPages", m.GetShowRecommendedPages())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("thumbnailWebUrl", m.GetThumbnailWebUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("titleArea", m.GetTitleArea())
        if err != nil {
            return err
        }
    }
    if m.GetWebParts() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetWebParts()))
        for i, v := range m.GetWebParts() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("webParts", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCanvasLayout sets the canvasLayout property value. Indicates the layout of the content in a given SharePoint page, including horizontal sections and vertical sections.
func (m *SitePage) SetCanvasLayout(value CanvasLayoutable)() {
    err := m.GetBackingStore().Set("canvasLayout", value)
    if err != nil {
        panic(err)
    }
}
// SetPromotionKind sets the promotionKind property value. Indicates the promotion kind of the sitePage. The possible values are: microsoftReserved, page, newsPost, unknownFutureValue.
func (m *SitePage) SetPromotionKind(value *PagePromotionType)() {
    err := m.GetBackingStore().Set("promotionKind", value)
    if err != nil {
        panic(err)
    }
}
// SetReactions sets the reactions property value. Reactions information for the page.
func (m *SitePage) SetReactions(value ReactionsFacetable)() {
    err := m.GetBackingStore().Set("reactions", value)
    if err != nil {
        panic(err)
    }
}
// SetShowComments sets the showComments property value. Determines whether or not to show comments at the bottom of the page.
func (m *SitePage) SetShowComments(value *bool)() {
    err := m.GetBackingStore().Set("showComments", value)
    if err != nil {
        panic(err)
    }
}
// SetShowRecommendedPages sets the showRecommendedPages property value. Determines whether or not to show recommended pages at the bottom of the page.
func (m *SitePage) SetShowRecommendedPages(value *bool)() {
    err := m.GetBackingStore().Set("showRecommendedPages", value)
    if err != nil {
        panic(err)
    }
}
// SetThumbnailWebUrl sets the thumbnailWebUrl property value. Url of the sitePage's thumbnail image
func (m *SitePage) SetThumbnailWebUrl(value *string)() {
    err := m.GetBackingStore().Set("thumbnailWebUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetTitleArea sets the titleArea property value. Title area on the SharePoint page.
func (m *SitePage) SetTitleArea(value TitleAreaable)() {
    err := m.GetBackingStore().Set("titleArea", value)
    if err != nil {
        panic(err)
    }
}
// SetWebParts sets the webParts property value. Collection of webparts on the SharePoint page.
func (m *SitePage) SetWebParts(value []WebPartable)() {
    err := m.GetBackingStore().Set("webParts", value)
    if err != nil {
        panic(err)
    }
}
type SitePageable interface {
    BaseSitePageable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCanvasLayout()(CanvasLayoutable)
    GetPromotionKind()(*PagePromotionType)
    GetReactions()(ReactionsFacetable)
    GetShowComments()(*bool)
    GetShowRecommendedPages()(*bool)
    GetThumbnailWebUrl()(*string)
    GetTitleArea()(TitleAreaable)
    GetWebParts()([]WebPartable)
    SetCanvasLayout(value CanvasLayoutable)()
    SetPromotionKind(value *PagePromotionType)()
    SetReactions(value ReactionsFacetable)()
    SetShowComments(value *bool)()
    SetShowRecommendedPages(value *bool)()
    SetThumbnailWebUrl(value *string)()
    SetTitleArea(value TitleAreaable)()
    SetWebParts(value []WebPartable)()
}
