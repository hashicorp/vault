package gocbcore

import (
	"encoding/json"
	"strconv"
	"time"
)

type kvErrorMapAttribute string

type kvErrorMapRetry struct {
	Strategy    string
	Interval    int
	After       int
	Ceil        int
	MaxDuration int
}

func (retry kvErrorMapRetry) CalculateRetryDelay(retryCount uint32) time.Duration {
	duraCeil := time.Duration(retry.Ceil) * time.Millisecond

	var dura time.Duration
	if retryCount == 0 {
		dura = time.Duration(retry.After) * time.Millisecond
	} else {
		interval := time.Duration(retry.Interval) * time.Millisecond
		if retry.Strategy == "constant" {
			dura = interval
		} else if retry.Strategy == "linear" {
			dura = interval * time.Duration(retryCount)
		} else if retry.Strategy == "exponential" {
			dura = interval
			for i := uint32(0); i < retryCount-1; i++ {
				// Need to multiply by the original value, not the scaled one
				dura = dura * time.Duration(retry.Interval)

				// We have to check this here to make sure we do not overflow
				if duraCeil > 0 && dura > duraCeil {
					dura = duraCeil
					break
				}
			}
		}
	}

	if duraCeil > 0 && dura > duraCeil {
		dura = duraCeil
	}

	return dura
}

type kvErrorMapError struct {
	Name        string
	Description string
	Attributes  []kvErrorMapAttribute
	Retry       kvErrorMapRetry
}

type kvErrorMap struct {
	Version  int
	Revision int
	Errors   map[uint16]kvErrorMapError
}

type cfgKvErrorMapError struct {
	Name  string   `json:"name"`
	Desc  string   `json:"desc"`
	Attrs []string `json:"attrs"`
	Retry struct {
		Strategy    string `json:"strategy"`
		Interval    int    `json:"interval"`
		After       int    `json:"after"`
		Ceil        int    `json:"ceil"`
		MaxDuration int    `json:"max-duration"`
	} `json:"retry"`
}

type cfgKvErrorMap struct {
	Version  int `json:"version"`
	Revision int `json:"revision"`
	Errors   map[string]cfgKvErrorMapError
}

func parseKvErrorMap(data []byte) (*kvErrorMap, error) {
	var cfg cfgKvErrorMap
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	var errMap kvErrorMap
	errMap.Version = cfg.Version
	errMap.Revision = cfg.Revision
	errMap.Errors = make(map[uint16]kvErrorMapError)
	for errCodeStr, errData := range cfg.Errors {
		errCode, err := strconv.ParseInt(errCodeStr, 16, 64)
		if err != nil {
			return nil, err
		}

		var errInfo kvErrorMapError
		errInfo.Name = errData.Name
		errInfo.Description = errData.Desc
		errInfo.Attributes = make([]kvErrorMapAttribute, len(errData.Attrs))
		for i, attr := range errData.Attrs {
			errInfo.Attributes[i] = kvErrorMapAttribute(attr)
		}
		errInfo.Retry.Strategy = errData.Retry.Strategy
		errInfo.Retry.Interval = errData.Retry.Interval
		errInfo.Retry.After = errData.Retry.After
		errInfo.Retry.Ceil = errData.Retry.Ceil
		errInfo.Retry.MaxDuration = errData.Retry.MaxDuration
		errMap.Errors[uint16(errCode)] = errInfo
	}

	return &errMap, nil
}
