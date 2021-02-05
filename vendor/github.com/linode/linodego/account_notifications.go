package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// Notification represents a notification on an Account
type Notification struct {
	Label    string               `json:"label"`
	Body     *string              `json:"body"`
	Message  string               `json:"message"`
	Type     NotificationType     `json:"type"`
	Severity NotificationSeverity `json:"severity"`
	Entity   *NotificationEntity  `json:"entity"`
	Until    *time.Time           `json:"-"`
	When     *time.Time           `json:"-"`
}

// NotificationEntity adds detailed information about the Notification.
// This could refer to the ticket that triggered the notification, for example.
type NotificationEntity struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

// NotificationSeverity constants start with Notification and include all known Linode API Notification Severities.
type NotificationSeverity string

// NotificationSeverity constants represent the actions that cause a Notification. New severities may be added in the future.
const (
	NotificationMinor    NotificationSeverity = "minor"
	NotificationMajor    NotificationSeverity = "major"
	NotificationCritical NotificationSeverity = "critical"
)

// NotificationType constants start with Notification and include all known Linode API Notification Types.
type NotificationType string

// NotificationType constants represent the actions that cause a Notification. New types may be added in the future.
const (
	NotificationMigrationScheduled NotificationType = "migration_scheduled"
	NotificationMigrationImminent  NotificationType = "migration_imminent"
	NotificationMigrationPending   NotificationType = "migration_pending"
	NotificationRebootScheduled    NotificationType = "reboot_scheduled"
	NotificationOutage             NotificationType = "outage"
	NotificationPaymentDue         NotificationType = "payment_due"
	NotificationTicketImportant    NotificationType = "ticket_important"
	NotificationTicketAbuse        NotificationType = "ticket_abuse"
	NotificationNotice             NotificationType = "notice"
	NotificationMaintenance        NotificationType = "maintenance"
)

// NotificationsPagedResponse represents a paginated Notifications API response
type NotificationsPagedResponse struct {
	*PageOptions
	Data []Notification `json:"data"`
}

// endpoint gets the endpoint URL for Notification
func (NotificationsPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.Notifications.Endpoint()
	if err != nil {
		panic(err)
	}

	return endpoint
}

// appendData appends Notifications when processing paginated Notification responses
func (resp *NotificationsPagedResponse) appendData(r *NotificationsPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// ListNotifications gets a collection of Notification objects representing important,
// often time-sensitive items related to the Account. An account cannot interact directly with
// Notifications, and a Notification will disappear when the circumstances causing it
// have been resolved. For example, if the account has an important Ticket open, a response
// to the Ticket will dismiss the Notification.
func (c *Client) ListNotifications(ctx context.Context, opts *ListOptions) ([]Notification, error) {
	response := NotificationsPagedResponse{}
	err := c.listHelper(ctx, &response, opts)

	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *Notification) UnmarshalJSON(b []byte) error {
	type Mask Notification

	p := struct {
		*Mask
		Until *parseabletime.ParseableTime `json:"until"`
		When  *parseabletime.ParseableTime `json:"when"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Until = (*time.Time)(p.Until)
	i.When = (*time.Time)(p.When)

	return nil
}
