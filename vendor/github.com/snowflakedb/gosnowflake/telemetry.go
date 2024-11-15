// Copyright (c) 2021-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	telemetryPath           = "/telemetry/send"
	defaultTelemetryTimeout = 10 * time.Second
	defaultFlushSize        = 100
)

const (
	typeKey          = "type"
	sourceKey        = "source"
	queryIDKey       = "QueryID"
	driverTypeKey    = "DriverType"
	driverVersionKey = "DriverVersion"
	golangVersionKey = "GolangVersion"
	sqlStateKey      = "SQLState"
	reasonKey        = "reason"
	errorNumberKey   = "ErrorNumber"
	stacktraceKey    = "Stacktrace"
)

const (
	telemetrySource      = "golang_driver"
	sqlException         = "client_sql_exception"
	connectionParameters = "client_connection_parameters"
)

type telemetryData struct {
	Timestamp int64             `json:"timestamp,omitempty"`
	Message   map[string]string `json:"message,omitempty"`
}

type snowflakeTelemetry struct {
	logs      []*telemetryData
	flushSize int
	sr        *snowflakeRestful
	mutex     *sync.Mutex
	enabled   bool
}

func (st *snowflakeTelemetry) addLog(data *telemetryData) error {
	if !st.enabled {
		return fmt.Errorf("telemetry disabled; not adding log")
	}
	st.mutex.Lock()
	st.logs = append(st.logs, data)
	st.mutex.Unlock()
	if len(st.logs) >= st.flushSize {
		if err := st.sendBatch(); err != nil {
			return err
		}
	}
	return nil
}

func (st *snowflakeTelemetry) sendBatch() error {
	if !st.enabled {
		err := fmt.Errorf("telemetry disabled; not sending log")
		logger.Debug(err)
		return err
	}
	type telemetry struct {
		Logs []*telemetryData `json:"logs"`
	}

	st.mutex.Lock()
	logsToSend := st.logs
	st.logs = make([]*telemetryData, 0)
	st.mutex.Unlock()

	if len(logsToSend) == 0 {
		logger.Debug("nothing to send to telemetry")
		return nil
	}

	s := &telemetry{logsToSend}
	body, err := json.Marshal(s)
	if err != nil {
		return err
	}
	logger.Debugf("sending %v logs to telemetry. inband telemetry payload "+
		"being sent: %v", len(logsToSend), string(body))

	headers := getHeaders()
	if token, _, _ := st.sr.TokenAccessor.GetTokens(); token != "" {
		headers[headerAuthorizationKey] = fmt.Sprintf(headerSnowflakeToken, token)
	}
	resp, err := st.sr.FuncPost(context.Background(), st.sr,
		st.sr.getFullURL(telemetryPath, nil), headers, body,
		defaultTelemetryTimeout, defaultTimeProvider, nil)
	if err != nil {
		logger.Info("failed to upload metrics to telemetry. err: %v", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("non-successful response from telemetry server: %v. "+
			"disabling telemetry", resp.StatusCode)
		logger.Info(err)
		st.enabled = false
		return err
	}
	var respd telemetryResponse
	if err = json.NewDecoder(resp.Body).Decode(&respd); err != nil {
		logger.Info(err)
		st.enabled = false
		return err
	}
	if !respd.Success {
		err = fmt.Errorf("telemetry send failed with error code: %v, message: %v",
			respd.Code, respd.Message)
		logger.Info(err)
		st.enabled = false
		return err
	}
	logger.Debug("successfully uploaded metrics to telemetry")
	return nil
}
