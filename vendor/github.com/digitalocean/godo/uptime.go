package godo

import (
	"context"
	"fmt"
	"net/http"
	"path"
)

const (
	uptimeChecksBasePath = "/v2/uptime/checks"
	// UptimeAlertGreaterThan is the comparison >
	UptimeAlertGreaterThan UptimeAlertComp = "greater_than"
	// UptimeAlertLessThan is the comparison <
	UptimeAlertLessThan UptimeAlertComp = "less_than"
)

// UptimeChecksService is an interface for creating and managing Uptime checks with the DigitalOcean API.
// See: https://docs.digitalocean.com/reference/api/api-reference/#tag/Uptime
type UptimeChecksService interface {
	List(context.Context, *ListOptions) ([]UptimeCheck, *Response, error)
	Get(context.Context, string) (*UptimeCheck, *Response, error)
	GetState(context.Context, string) (*UptimeCheckState, *Response, error)
	Create(context.Context, *CreateUptimeCheckRequest) (*UptimeCheck, *Response, error)
	Update(context.Context, string, *UpdateUptimeCheckRequest) (*UptimeCheck, *Response, error)
	Delete(context.Context, string) (*Response, error)
	GetAlert(context.Context, string, string) (*UptimeAlert, *Response, error)
	ListAlerts(context.Context, string, *ListOptions) ([]UptimeAlert, *Response, error)
	CreateAlert(context.Context, string, *CreateUptimeAlertRequest) (*UptimeAlert, *Response, error)
	UpdateAlert(context.Context, string, string, *UpdateUptimeAlertRequest) (*UptimeAlert, *Response, error)
	DeleteAlert(context.Context, string, string) (*Response, error)
}

// UptimeChecksServiceOp handles communication with Uptime Check methods of the DigitalOcean API.
type UptimeChecksServiceOp struct {
	client *Client
}

// UptimeCheck represents a DigitalOcean UptimeCheck configuration.
type UptimeCheck struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Target  string   `json:"target"`
	Regions []string `json:"regions"`
	Enabled bool     `json:"enabled"`
}

// UptimeAlert represents a DigitalOcean Uptime Alert configuration.
type UptimeAlert struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Type          string          `json:"type"`
	Threshold     int             `json:"threshold"`
	Comparison    UptimeAlertComp `json:"comparison"`
	Notifications *Notifications  `json:"notifications"`
	Period        string          `json:"period"`
}

// Notifications represents a DigitalOcean Notifications configuration.
type Notifications struct {
	Email []string       `json:"email"`
	Slack []SlackDetails `json:"slack"`
}

// UptimeCheckState represents a DigitalOcean Uptime Check's state configuration.
type UptimeCheckState struct {
	Regions        map[string]UptimeRegion `json:"regions"`
	PreviousOutage UptimePreviousOutage    `json:"previous_outage"`
}

type UptimeRegion struct {
	Status                    string  `json:"status"`
	StatusChangedAt           string  `json:"status_changed_at"`
	ThirtyDayUptimePercentage float32 `json:"thirty_day_uptime_percentage"`
}

// UptimePreviousOutage represents a DigitalOcean Uptime Check's previous outage configuration.
type UptimePreviousOutage struct {
	Region          string `json:"region"`
	StartedAt       string `json:"started_at"`
	EndedAt         string `json:"ended_at"`
	DurationSeconds int    `json:"duration_seconds"`
}

// CreateUptimeCheckRequest represents the request to create a new uptime check.
type CreateUptimeCheckRequest struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Target  string   `json:"target"`
	Regions []string `json:"regions"`
	Enabled bool     `json:"enabled"`
}

// UpdateUptimeCheckRequest represents the request to update uptime check information.
type UpdateUptimeCheckRequest struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Target  string   `json:"target"`
	Regions []string `json:"regions"`
	Enabled bool     `json:"enabled"`
}

// CreateUptimeUptimeAlertRequest represents the request to create a new Uptime Alert.
type CreateUptimeAlertRequest struct {
	Name          string          `json:"name"`
	Type          string          `json:"type"`
	Threshold     int             `json:"threshold"`
	Comparison    UptimeAlertComp `json:"comparison"`
	Notifications *Notifications  `json:"notifications"`
	Period        string          `json:"period"`
}

// UpdateUptimeAlertRequest represents the request to update an alert.
type UpdateUptimeAlertRequest struct {
	Name          string          `json:"name"`
	Type          string          `json:"type"`
	Threshold     int             `json:"threshold"`
	Comparison    UptimeAlertComp `json:"comparison"`
	Notifications *Notifications  `json:"notifications"`
	Period        string          `json:"period"`
}

// UptimeAlertComp represents an uptime alert comparison operation
type UptimeAlertComp string

type uptimeChecksRoot struct {
	UptimeChecks []UptimeCheck `json:"checks"`
	Links        *Links        `json:"links"`
	Meta         *Meta         `json:"meta"`
}

type uptimeCheckStateRoot struct {
	UptimeCheckState UptimeCheckState `json:"state"`
}

type uptimeAlertsRoot struct {
	UptimeAlerts []UptimeAlert `json:"alerts"`
	Links        *Links        `json:"links"`
	Meta         *Meta         `json:"meta"`
}

type uptimeCheckRoot struct {
	UptimeCheck *UptimeCheck `json:"check"`
}

type uptimeAlertRoot struct {
	UptimeAlert *UptimeAlert `json:"alert"`
}

var _ UptimeChecksService = &UptimeChecksServiceOp{}

// List Checks.
func (p *UptimeChecksServiceOp) List(ctx context.Context, opts *ListOptions) ([]UptimeCheck, *Response, error) {
	path, err := addOptions(uptimeChecksBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(uptimeChecksRoot)
	resp, err := p.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.UptimeChecks, resp, err
}

// GetState of uptime check.
func (p *UptimeChecksServiceOp) GetState(ctx context.Context, uptimeCheckID string) (*UptimeCheckState, *Response, error) {
	path := path.Join(uptimeChecksBasePath, uptimeCheckID, "/state")

	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(uptimeCheckStateRoot)
	resp, err := p.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return &root.UptimeCheckState, resp, err
}

// Get retrieves a single uptime check by its ID.
func (p *UptimeChecksServiceOp) Get(ctx context.Context, uptimeCheckID string) (*UptimeCheck, *Response, error) {
	path := path.Join(uptimeChecksBasePath, uptimeCheckID)

	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(uptimeCheckRoot)
	resp, err := p.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.UptimeCheck, resp, err
}

// Create a new uptime check.
func (p *UptimeChecksServiceOp) Create(ctx context.Context, cr *CreateUptimeCheckRequest) (*UptimeCheck, *Response, error) {
	req, err := p.client.NewRequest(ctx, http.MethodPost, uptimeChecksBasePath, cr)
	if err != nil {
		return nil, nil, err
	}

	root := new(uptimeCheckRoot)
	resp, err := p.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.UptimeCheck, resp, err
}

// Update an uptime check.
func (p *UptimeChecksServiceOp) Update(ctx context.Context, uptimeCheckID string, ur *UpdateUptimeCheckRequest) (*UptimeCheck, *Response, error) {
	path := path.Join(uptimeChecksBasePath, uptimeCheckID)
	req, err := p.client.NewRequest(ctx, http.MethodPut, path, ur)
	if err != nil {
		return nil, nil, err
	}

	root := new(uptimeCheckRoot)
	resp, err := p.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.UptimeCheck, resp, err
}

// Delete an existing uptime check.
func (p *UptimeChecksServiceOp) Delete(ctx context.Context, uptimeCheckID string) (*Response, error) {
	path := path.Join(uptimeChecksBasePath, uptimeCheckID)
	req, err := p.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return p.client.Do(ctx, req, nil)
}

// alerts

// ListAlerts lists alerts for a check.
func (p *UptimeChecksServiceOp) ListAlerts(ctx context.Context, uptimeCheckID string, opts *ListOptions) ([]UptimeAlert, *Response, error) {
	fullPath := path.Join(uptimeChecksBasePath, uptimeCheckID, "/alerts")
	path, err := addOptions(fullPath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(uptimeAlertsRoot)
	resp, err := p.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.UptimeAlerts, resp, err
}

// CreateAlert creates a new check alert.
func (p *UptimeChecksServiceOp) CreateAlert(ctx context.Context, uptimeCheckID string, cr *CreateUptimeAlertRequest) (*UptimeAlert, *Response, error) {
	fullPath := path.Join(uptimeChecksBasePath, uptimeCheckID, "/alerts")
	req, err := p.client.NewRequest(ctx, http.MethodPost, fullPath, cr)
	if err != nil {
		return nil, nil, err
	}

	root := new(uptimeAlertRoot)
	resp, err := p.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.UptimeAlert, resp, err
}

// GetAlert retrieves a single uptime check alert by its ID.
func (p *UptimeChecksServiceOp) GetAlert(ctx context.Context, uptimeCheckID string, alertID string) (*UptimeAlert, *Response, error) {
	path := fmt.Sprintf("v2/uptime/checks/%s/alerts/%s", uptimeCheckID, alertID)

	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(uptimeAlertRoot)
	resp, err := p.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.UptimeAlert, resp, err
}

// UpdateAlert updates an check's alert.
func (p *UptimeChecksServiceOp) UpdateAlert(ctx context.Context, uptimeCheckID string, alertID string, ur *UpdateUptimeAlertRequest) (*UptimeAlert, *Response, error) {
	path := path.Join(uptimeChecksBasePath, uptimeCheckID, "/alerts/", alertID)
	req, err := p.client.NewRequest(ctx, http.MethodPut, path, ur)
	if err != nil {
		return nil, nil, err
	}

	root := new(uptimeAlertRoot)
	resp, err := p.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.UptimeAlert, resp, err
}

// DeleteAlert deletes an existing check's alert.
func (p *UptimeChecksServiceOp) DeleteAlert(ctx context.Context, uptimeCheckID string, alertID string) (*Response, error) {
	path := path.Join(uptimeChecksBasePath, uptimeCheckID, "/alerts/", alertID)
	req, err := p.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return p.client.Do(ctx, req, nil)
}
