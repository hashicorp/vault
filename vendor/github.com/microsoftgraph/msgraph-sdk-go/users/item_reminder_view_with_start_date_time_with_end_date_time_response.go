package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemReminderViewWithStartDateTimeWithEndDateTimeGetResponseable instead.
type ItemReminderViewWithStartDateTimeWithEndDateTimeResponse struct {
    ItemReminderViewWithStartDateTimeWithEndDateTimeGetResponse
}
// NewItemReminderViewWithStartDateTimeWithEndDateTimeResponse instantiates a new ItemReminderViewWithStartDateTimeWithEndDateTimeResponse and sets the default values.
func NewItemReminderViewWithStartDateTimeWithEndDateTimeResponse()(*ItemReminderViewWithStartDateTimeWithEndDateTimeResponse) {
    m := &ItemReminderViewWithStartDateTimeWithEndDateTimeResponse{
        ItemReminderViewWithStartDateTimeWithEndDateTimeGetResponse: *NewItemReminderViewWithStartDateTimeWithEndDateTimeGetResponse(),
    }
    return m
}
// CreateItemReminderViewWithStartDateTimeWithEndDateTimeResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemReminderViewWithStartDateTimeWithEndDateTimeResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemReminderViewWithStartDateTimeWithEndDateTimeResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemReminderViewWithStartDateTimeWithEndDateTimeGetResponseable instead.
type ItemReminderViewWithStartDateTimeWithEndDateTimeResponseable interface {
    ItemReminderViewWithStartDateTimeWithEndDateTimeGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}
