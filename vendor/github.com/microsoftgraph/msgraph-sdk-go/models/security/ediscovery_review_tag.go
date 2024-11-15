package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EdiscoveryReviewTag struct {
    Tag
}
// NewEdiscoveryReviewTag instantiates a new EdiscoveryReviewTag and sets the default values.
func NewEdiscoveryReviewTag()(*EdiscoveryReviewTag) {
    m := &EdiscoveryReviewTag{
        Tag: *NewTag(),
    }
    odataTypeValue := "#microsoft.graph.security.ediscoveryReviewTag"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEdiscoveryReviewTagFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEdiscoveryReviewTagFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEdiscoveryReviewTag(), nil
}
// GetChildSelectability gets the childSelectability property value. Indicates whether a single or multiple child tags can be associated with a document. Possible values are: One, Many.  This value controls whether the UX presents the tags as checkboxes or a radio button group.
// returns a *ChildSelectability when successful
func (m *EdiscoveryReviewTag) GetChildSelectability()(*ChildSelectability) {
    val, err := m.GetBackingStore().Get("childSelectability")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ChildSelectability)
    }
    return nil
}
// GetChildTags gets the childTags property value. Returns the tags that are a child of a tag.
// returns a []EdiscoveryReviewTagable when successful
func (m *EdiscoveryReviewTag) GetChildTags()([]EdiscoveryReviewTagable) {
    val, err := m.GetBackingStore().Get("childTags")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]EdiscoveryReviewTagable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EdiscoveryReviewTag) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Tag.GetFieldDeserializers()
    res["childSelectability"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseChildSelectability)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetChildSelectability(val.(*ChildSelectability))
        }
        return nil
    }
    res["childTags"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateEdiscoveryReviewTagFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]EdiscoveryReviewTagable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(EdiscoveryReviewTagable)
                }
            }
            m.SetChildTags(res)
        }
        return nil
    }
    res["parent"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEdiscoveryReviewTagFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParent(val.(EdiscoveryReviewTagable))
        }
        return nil
    }
    return res
}
// GetParent gets the parent property value. Returns the parent tag of the specified tag.
// returns a EdiscoveryReviewTagable when successful
func (m *EdiscoveryReviewTag) GetParent()(EdiscoveryReviewTagable) {
    val, err := m.GetBackingStore().Get("parent")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EdiscoveryReviewTagable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EdiscoveryReviewTag) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Tag.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetChildSelectability() != nil {
        cast := (*m.GetChildSelectability()).String()
        err = writer.WriteStringValue("childSelectability", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetChildTags() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetChildTags()))
        for i, v := range m.GetChildTags() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("childTags", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("parent", m.GetParent())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetChildSelectability sets the childSelectability property value. Indicates whether a single or multiple child tags can be associated with a document. Possible values are: One, Many.  This value controls whether the UX presents the tags as checkboxes or a radio button group.
func (m *EdiscoveryReviewTag) SetChildSelectability(value *ChildSelectability)() {
    err := m.GetBackingStore().Set("childSelectability", value)
    if err != nil {
        panic(err)
    }
}
// SetChildTags sets the childTags property value. Returns the tags that are a child of a tag.
func (m *EdiscoveryReviewTag) SetChildTags(value []EdiscoveryReviewTagable)() {
    err := m.GetBackingStore().Set("childTags", value)
    if err != nil {
        panic(err)
    }
}
// SetParent sets the parent property value. Returns the parent tag of the specified tag.
func (m *EdiscoveryReviewTag) SetParent(value EdiscoveryReviewTagable)() {
    err := m.GetBackingStore().Set("parent", value)
    if err != nil {
        panic(err)
    }
}
type EdiscoveryReviewTagable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    Tagable
    GetChildSelectability()(*ChildSelectability)
    GetChildTags()([]EdiscoveryReviewTagable)
    GetParent()(EdiscoveryReviewTagable)
    SetChildSelectability(value *ChildSelectability)()
    SetChildTags(value []EdiscoveryReviewTagable)()
    SetParent(value EdiscoveryReviewTagable)()
}
