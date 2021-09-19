package api

import (
	"context"
	"fmt"
	"time"
)

func (c *Sys) ListCounters() (*Counter, error) {
	r := c.c.NewRequest("GET", "/v1/sys/internal/counters/config")

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Counter
	err = resp.DecodeJSON(&result)

	return &result, err
}

func (c *Sys) ConfigCounters(config CountersConfig) error {
	r := c.c.NewRequest("POST", fmt.Sprintf("/v1/sys/internal/counters/config"))
	if err := r.SetJSONBody(config); err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *Sys) RequestsCounters() (*CountersRequests, error) {
	r := c.c.NewRequest("GET", "/v1/sys/internal/counters/requests")

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result CountersRequests
	err = resp.DecodeJSON(&result)

	return &result, err
}

func (c *Sys) EntitiesCounters() (*CountersEntities, error) {
	r := c.c.NewRequest("GET", "/v1/sys/internal/counters/entities")

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result CountersEntities
	err = resp.DecodeJSON(&result)

	return &result, err
}

func (c *Sys) TokensCounters() (*CountersToken, error) {
	r := c.c.NewRequest("GET", "/v1/sys/internal/counters/tokens")

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result CountersToken
	err = resp.DecodeJSON(&result)

	return &result, err
}

func (c *Sys) ActivityCounters() (*CountersActivity, error) {
	r := c.c.NewRequest("GET", "/v1/sys/internal/counters/activity")

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result CountersActivity
	err = resp.DecodeJSON(&result)

	return &result, err
}

func (c *Sys) ActivityMonthlyCounters() (*CountersActivityMonthly, error) {
	r := c.c.NewRequest("GET", "/v1/sys/internal/counters/activity/monthly")

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result CountersActivityMonthly
	err = resp.DecodeJSON(&result)

	return &result, err
}

type Counter struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
		DefaultReportMonths int    `json:"default_report_months"`
		Enabled             string `json:"enabled"`
		QueriesAvailable    bool   `json:"queries_available"`
		RetentionMonths     int    `json:"retention_months"`
	} `json:"data"`
	WrapInfo interface{} `json:"wrap_info"`
	Warnings interface{} `json:"warnings"`
	Auth     interface{} `json:"auth"`
}

type CountersConfig struct {
	Default_report_months int    `json:"default_report_months"`
	Enabled               string `json:"enabled"`
	Retention_months      int    `json:"retention_months"`
}

type CountersRequests struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
		Counters []struct {
			StartTime time.Time `json:"start_time"`
			Total     int       `json:"total"`
		} `json:"counters"`
	} `json:"data"`
	WrapInfo interface{} `json:"wrap_info"`
	Warnings interface{} `json:"warnings"`
	Auth     interface{} `json:"auth"`
}

type CountersEntities struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
		Counters struct {
			Entities struct {
				Total int `json:"total"`
			} `json:"entities"`
		} `json:"counters"`
	} `json:"data"`
	WrapInfo interface{} `json:"wrap_info"`
	Warnings interface{} `json:"warnings"`
	Auth     interface{} `json:"auth"`
}

type CountersToken struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
		Counters struct {
			ServiceTokens struct {
				Total int `json:"total"`
			} `json:"service_tokens"`
		} `json:"counters"`
	} `json:"data"`
	WrapInfo interface{} `json:"wrap_info"`
	Warnings interface{} `json:"warnings"`
	Auth     interface{} `json:"auth"`
}

type CountersActivity struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
		Total     struct {
			DistinctEntities int `json:"distinct_entities"`
			NonEntityTokens  int `json:"non_entity_tokens"`
			Clients          int `json:"clients"`
		} `json:"total"`
		ByNamespace []struct {
			NamespaceID   string `json:"namespace_id"`
			NamespacePath string `json:"namespace_path"`
			Counts        struct {
				DistinctEntities int `json:"distinct_entities"`
				NonEntityTokens  int `json:"non_entity_tokens"`
				Clients          int `json:"clients"`
			} `json:"counts"`
		} `json:"by_namespace"`
	} `json:"data"`
	WrapInfo interface{} `json:"wrap_info"`
	Warnings interface{} `json:"warnings"`
	Auth     interface{} `json:"auth"`
}

type CountersActivityMonthly struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
		DistinctEntities int `json:"distinct_entities"`
		NonEntityTokens  int `json:"non_entity_tokens"`
		Clients          int `json:"clients"`
	} `json:"data"`
	WrapInfo interface{} `json:"wrap_info"`
	Warnings interface{} `json:"warnings"`
	Auth     interface{} `json:"auth"`
}
