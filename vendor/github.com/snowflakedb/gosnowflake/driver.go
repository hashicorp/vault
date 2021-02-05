// Copyright (c) 2017-2019 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"net/http"
)

// SnowflakeDriver is a context of Go Driver
type SnowflakeDriver struct{}

// Open creates a new connection.
func (d SnowflakeDriver) Open(dsn string) (driver.Conn, error) {
	glog.V(2).Info("Open")
	var err error
	sc := &snowflakeConn{
		SequenceCounter: 0,
	}
	ctx := context.TODO()
	sc.cfg, err = ParseDSN(dsn)
	if err != nil {
		sc.cleanup()
		return nil, err
	}
	st := SnowflakeTransport
	if sc.cfg.InsecureMode {
		// no revocation check with OCSP. Think twice when you want to enable this option.
		st = snowflakeInsecureTransport
	} else {
		// set OCSP fail open mode
		ocspResponseCacheLock.Lock()
		ocspFailOpen = sc.cfg.OCSPFailOpen
		ocspResponseCacheLock.Unlock()
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

	glog.V(2).Infof("Authenticating via %v", sc.cfg.Authenticator.String())
	switch sc.cfg.Authenticator {
	case AuthTypeExternalBrowser:
		samlResponse, proofKey, err = authenticateByExternalBrowser(
			ctx,
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
			ctx,
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
		ctx,
		sc,
		samlResponse,
		proofKey)
	if err != nil {
		sc.cleanup()
		return nil, err
	}

	sc.populateSessionParameters(authData.Parameters)
	sc.startHeartBeat()
	return sc, nil
}

func init() {
	sql.Register("snowflake", &SnowflakeDriver{})
}
