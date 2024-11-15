package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessReviewNotificationRecipientQueryScope struct {
    AccessReviewNotificationRecipientScope
}
// NewAccessReviewNotificationRecipientQueryScope instantiates a new AccessReviewNotificationRecipientQueryScope and sets the default values.
func NewAccessReviewNotificationRecipientQueryScope()(*AccessReviewNotificationRecipientQueryScope) {
    m := &AccessReviewNotificationRecipientQueryScope{
        AccessReviewNotificationRecipientScope: *NewAccessReviewNotificationRecipientScope(),
    }
    odataTypeValue := "#microsoft.graph.accessReviewNotificationRecipientQueryScope"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAccessReviewNotificationRecipientQueryScopeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessReviewNotificationRecipientQueryScopeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessReviewNotificationRecipientQueryScope(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AccessReviewNotificationRecipientQueryScope) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AccessReviewNotificationRecipientScope.GetFieldDeserializers()
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
// GetQuery gets the query property value. Represents the query for who the recipients are. For example, /groups/{group id}/members for group members and /users/{user id} for a specific user.
// returns a *string when successful
func (m *AccessReviewNotificationRecipientQueryScope) GetQuery()(*string) {
    val, err := m.GetBackingStore().Get("query")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetQueryRoot gets the queryRoot property value. In the scenario where reviewers need to be specified dynamically, indicates the relative source of the query. This property is only required if a relative query (that is, ./manager) is specified.
// returns a *string when successful
func (m *AccessReviewNotificationRecipientQueryScope) GetQueryRoot()(*string) {
    val, err := m.GetBackingStore().Get("queryRoot")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetQueryType gets the queryType property value. Indicates the type of query. Allowed value is MicrosoftGraph.
// returns a *string when successful
func (m *AccessReviewNotificationRecipientQueryScope) GetQueryType()(*string) {
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
func (m *AccessReviewNotificationRecipientQueryScope) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AccessReviewNotificationRecipientScope.Serialize(writer)
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
// SetQuery sets the query property value. Represents the query for who the recipients are. For example, /groups/{group id}/members for group members and /users/{user id} for a specific user.
func (m *AccessReviewNotificationRecipientQueryScope) SetQuery(value *string)() {
    err := m.GetBackingStore().Set("query", value)
    if err != nil {
        panic(err)
    }
}
// SetQueryRoot sets the queryRoot property value. In the scenario where reviewers need to be specified dynamically, indicates the relative source of the query. This property is only required if a relative query (that is, ./manager) is specified.
func (m *AccessReviewNotificationRecipientQueryScope) SetQueryRoot(value *string)() {
    err := m.GetBackingStore().Set("queryRoot", value)
    if err != nil {
        panic(err)
    }
}
// SetQueryType sets the queryType property value. Indicates the type of query. Allowed value is MicrosoftGraph.
func (m *AccessReviewNotificationRecipientQueryScope) SetQueryType(value *string)() {
    err := m.GetBackingStore().Set("queryType", value)
    if err != nil {
        panic(err)
    }
}
type AccessReviewNotificationRecipientQueryScopeable interface {
    AccessReviewNotificationRecipientScopeable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetQuery()(*string)
    GetQueryRoot()(*string)
    GetQueryType()(*string)
    SetQuery(value *string)()
    SetQueryRoot(value *string)()
    SetQueryType(value *string)()
}
