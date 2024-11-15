package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessReviewQueryScope struct {
    AccessReviewScope
}
// NewAccessReviewQueryScope instantiates a new AccessReviewQueryScope and sets the default values.
func NewAccessReviewQueryScope()(*AccessReviewQueryScope) {
    m := &AccessReviewQueryScope{
        AccessReviewScope: *NewAccessReviewScope(),
    }
    odataTypeValue := "#microsoft.graph.accessReviewQueryScope"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAccessReviewQueryScopeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessReviewQueryScopeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.accessReviewInactiveUsersQueryScope":
                        return NewAccessReviewInactiveUsersQueryScope(), nil
                }
            }
        }
    }
    return NewAccessReviewQueryScope(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AccessReviewQueryScope) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AccessReviewScope.GetFieldDeserializers()
    res["query"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQuery(val)
        }
        return nil
    }
    res["queryRoot"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQueryRoot(val)
        }
        return nil
    }
    res["queryType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQueryType(val)
        }
        return nil
    }
    return res
}
// GetQuery gets the query property value. The query representing what will be reviewed in an access review.
// returns a *string when successful
func (m *AccessReviewQueryScope) GetQuery()(*string) {
    val, err := m.GetBackingStore().Get("query")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetQueryRoot gets the queryRoot property value. In the scenario where reviewers need to be specified dynamically, this property is used to indicate the relative source of the query. This property is only required if a relative query is specified. For example, ./manager.
// returns a *string when successful
func (m *AccessReviewQueryScope) GetQueryRoot()(*string) {
    val, err := m.GetBackingStore().Get("queryRoot")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetQueryType gets the queryType property value. Indicates the type of query. Types include MicrosoftGraph and ARM.
// returns a *string when successful
func (m *AccessReviewQueryScope) GetQueryType()(*string) {
    val, err := m.GetBackingStore().Get("queryType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessReviewQueryScope) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AccessReviewScope.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("query", m.GetQuery())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("queryRoot", m.GetQueryRoot())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("queryType", m.GetQueryType())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetQuery sets the query property value. The query representing what will be reviewed in an access review.
func (m *AccessReviewQueryScope) SetQuery(value *string)() {
    err := m.GetBackingStore().Set("query", value)
    if err != nil {
        panic(err)
    }
}
// SetQueryRoot sets the queryRoot property value. In the scenario where reviewers need to be specified dynamically, this property is used to indicate the relative source of the query. This property is only required if a relative query is specified. For example, ./manager.
func (m *AccessReviewQueryScope) SetQueryRoot(value *string)() {
    err := m.GetBackingStore().Set("queryRoot", value)
    if err != nil {
        panic(err)
    }
}
// SetQueryType sets the queryType property value. Indicates the type of query. Types include MicrosoftGraph and ARM.
func (m *AccessReviewQueryScope) SetQueryType(value *string)() {
    err := m.GetBackingStore().Set("queryType", value)
    if err != nil {
        panic(err)
    }
}
type AccessReviewQueryScopeable interface {
    AccessReviewScopeable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetQuery()(*string)
    GetQueryRoot()(*string)
    GetQueryType()(*string)
    SetQuery(value *string)()
    SetQueryRoot(value *string)()
    SetQueryType(value *string)()
}
