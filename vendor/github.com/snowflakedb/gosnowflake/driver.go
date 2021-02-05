// Copyright (c) 2017-2021 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"net/http"
	"sync/atomic"
)

// SnowflakeDriver is a context of Go Driver
type SnowflakeDriver struct {
}

// Open creates a new connection.
func (d SnowflakeDriver) Open(dsn string) (driver.Conn, error) {
	logger.Info("Open")
	ctx := context.TODO()
	cfg, err := ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	return d.OpenWithConfig(ctx, *cfg)
}

// OpenWithConfig creates a new connection with the given Config.
func (d SnowflakeDriver) OpenWithConfig(ctx context.Context, config Config) (driver.Conn, error) {
	logger.Info("OpenWithConfig")
	var err error
	sc := &snowflakeConn{
		SequenceCounter: 0,
		ctx:             ctx,
		cfg:             &config,
	}
	var st http.RoundTripper = SnowflakeTransport
	if sc.cfg.Transporter == nil {
		if sc.cfg.InsecureMode {
			// no revocation check with OCSP. Think twice when you want to enable this option.
			st = snowflakeInsecureTransport
		} else {
			// set OCSP fail open mode
			ocspResponseCacheLock.Lock()
			atomic.StoreUint32((*uint32)(&ocspFailOpen), uint32(sc.cfg.OCSPFailOpen))
			ocspResponseCacheLock.Unlock()
		}
	} else {
		// use the custom transport
		st = sc.cfg.Transporter
	}
	// authenticate
	sc.rest = &snowflakeRestful{
		Host:     sc.cfg.Host,
		Port:     sc.cfg.Port,
		Protocol: sc.cfg.Protocol,
		Client: &http.Client{
			// request timeout including reading response body
			Timeout:   defaultClientTimeout,
			Transport: st,
		},
		LoginTimeout:        sc.cfg.LoginTimeout,
		RequestTimeout:      sc.cfg.RequestTimeout,
		FuncPost:            postRestful,
		FuncGet:             getRestful,
		FuncPostQuery:       postRestfulQuery,
		FuncPostQueryHelper: postRestfulQueryHelper,
		FuncRenewSession:    renewRestfulSession,
		FuncPostAuth:        postAuth,
		FuncCloseSession:    closeSession,
		FuncCancelQuery:     cancelQuery,
		FuncPostAuthSAML:    postAuthSAML,
		FuncPostAuthOKTA:    postAuthOKTA,
		FuncGetSSO:          getSSO,
	}
	var authData *authResponseMain
	var samlResponse []byte
	var proofKey []byte

	logger.Infof("Authenticating via %v", sc.cfg.Authenticator.String())
	switch sc.cfg.Authenticator {
	case AuthTypeExternalBrowser:
		samlResponse, proofKey, err = authenticateByExternalBrowser(
			sc.ctx,
			sc.rest,
			sc.cfg.Authenticator.String(),
			sc.cfg.Application,
			sc.cfg.Account,
			sc.cfg.User,
			sc.cfg.Password)
		if err != nil {
			sc.cleanup()
			return nil, err
		}
	case AuthTypeOkta:
		samlResponse, err = authenticateBySAML(
			sc.ctx,
			sc.rest,
			sc.cfg.OktaURL,
			sc.cfg.Application,
			sc.cfg.Account,
			sc.cfg.User,
			sc.cfg.Password)
		if err != nil {
			sc.cleanup()
			return nil, err
		}
	}
	authData, err = authenticate(
		sc.ctx,
		sc,
		samlResponse,
		proofKey)
	if err != nil {
		sc.cleanup()
		return nil, err
	}

	sc.populateSessionParameters(authData.Parameters)
	sc.ctx = context.WithValue(sc.ctx, SFSessionIDKey, authData.SessionID)
	sc.startHeartBeat()
	return sc, nil
}

var logger = CreateDefaultLogger()

func init() {
	sql.Register("snowflake", &SnowflakeDriver{})
	logger.SetLogLevel("error")
}
