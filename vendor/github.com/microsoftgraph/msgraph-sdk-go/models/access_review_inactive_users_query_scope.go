package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessReviewInactiveUsersQueryScope struct {
    AccessReviewQueryScope
}
// NewAccessReviewInactiveUsersQueryScope instantiates a new AccessReviewInactiveUsersQueryScope and sets the default values.
func NewAccessReviewInactiveUsersQueryScope()(*AccessReviewInactiveUsersQueryScope) {
    m := &AccessReviewInactiveUsersQueryScope{
        AccessReviewQueryScope: *NewAccessReviewQueryScope(),
    }
    odataTypeValue := "#microsoft.graph.accessReviewInactiveUsersQueryScope"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAccessReviewInactiveUsersQueryScopeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessReviewInactiveUsersQueryScopeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessReviewInactiveUsersQueryScope(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AccessReviewInactiveUsersQueryScope) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AccessReviewQueryScope.GetFieldDeserializers()
    res["inactiveDuration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInactiveDuration(val)
        }
        return nil
    }
    return res
}
// GetInactiveDuration gets the inactiveDuration property value. Defines the duration of inactivity. Inactivity is based on the last sign in date of the user compared to the access review instance's start date. If this property is not specified, it's assigned the default value PT0S.
// returns a *ISODuration when successful
func (m *AccessReviewInactiveUsersQueryScope) GetInactiveDuration()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("inactiveDuration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessReviewInactiveUsersQueryScope) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AccessReviewQueryScope.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteISODurationValue("inactiveDuration", m.GetInactiveDuration())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetInactiveDuration sets the inactiveDuration property value. Defines the duration of inactivity. Inactivity is based on the last sign in date of the user compared to the access review instance's start date. If this property is not specified, it's assigned the default value PT0S.
func (m *AccessReviewInactiveUsersQueryScope) SetInactiveDuration(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("inactiveDuration", value)
    if err != nil {
        panic(err)
    }
}
type AccessReviewInactiveUsersQueryScopeable interface {
    AccessReviewQueryScopeable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetInactiveDuration()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    SetInactiveDuration(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
}
