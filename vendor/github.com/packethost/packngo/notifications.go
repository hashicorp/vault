package packngo

import (
	"path"
)

const notificationBasePath = "/notifications"

// Notification struct
type Notification struct {
	ID        string    `json:"id,omitempty"`
	Type      string    `json:"type,omitempty"`
	Body      string    `json:"body,omitempty"`
	Severity  string    `json:"severity,omitempty"`
	Read      bool      `json:"read,omitempty"`
	Context   string    `json:"context,omitempty"`
	CreatedAt Timestamp `json:"created_at,omitempty"`
	UpdatedAt Timestamp `json:"updated_at,omitempty"`
	User      Href      `json:"user,omitempty"`
	Href      string    `json:"href,omitempty"`
}

type notificationsRoot struct {
	Notifications []Notification `json:"notifications,omitempty"`
	Meta          meta           `json:"meta,omitempty"`
}

// NotificationService interface defines available event functions
type NotificationService interface {
	List(*ListOptions) ([]Notification, *Response, error)
	Get(string, *GetOptions) (*Notification, *Response, error)
	MarkAsRead(string) (*Notification, *Response, error)
}

// NotificationServiceOp implements NotificationService
type NotificationServiceOp struct {
	client *Client
}

// List returns all notifications
func (s *NotificationServiceOp) List(listOpt *ListOptions) ([]Notification, *Response, error) {
	return listNotifications(s.client, notificationBasePath, listOpt)
}

// Get returns a notification by ID
func (s *NotificationServiceOp) Get(notificationID string, opts *GetOptions) (*Notification, *Response, error) {
	if validateErr := ValidateUUID(notificationID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(notificationBasePath, notificationID)
	apiPathQuery := opts.WithQuery(endpointPath)
	return getNotifications(s.client, apiPathQuery)
}

// Marks notification as read by ID
func (s *NotificationServiceOp) MarkAsRead(notificationID string) (*Notification, *Response, error) {
	if validateErr := ValidateUUID(notificationID); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(notificationBasePath, notificationID)
	return markAsRead(s.client, apiPath)
}

// list helper function for all notification functions
func listNotifications(client *Client, endpointPath string, opts *ListOptions) ([]Notification, *Response, error) {
	root := new(notificationsRoot)

	apiPathQuery := opts.WithQuery(endpointPath)

	resp, err := client.DoRequest("GET", apiPathQuery, nil, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Notifications, resp, err
}

func getNotifications(client *Client, apiPath string) (*Notification, *Response, error) {

	notification := new(Notification)

	resp, err := client.DoRequest("GET", apiPath, nil, notification)
	if err != nil {
		return nil, resp, err
	}

	return notification, resp, err
}

func markAsRead(client *Client, apiPath string) (*Notification, *Response, error) {

	notification := new(Notification)

	resp, err := client.DoRequest("PUT", apiPath, nil, notification)
	if err != nil {
		return nil, resp, err
	}

	return notification, resp, err
}
