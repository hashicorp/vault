// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

//+build gssapi
//+build linux darwin

package gssapi

/*
#cgo linux CFLAGS: -DGOOS_linux
#cgo linux LDFLAGS: -lgssapi_krb5 -lkrb5
#cgo darwin CFLAGS: -DGOOS_darwin
#cgo darwin LDFLAGS: -framework GSS
#include "gss_wrapper.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"strings"
	"unsafe"
)

// New creates a new SaslClient. The target parameter should be a hostname with no port.
func New(target, username, password string, passwordSet bool, props map[string]string) (*SaslClient, error) {
	serviceName := "mongodb"

	for key, value := range props {
		switch strings.ToUpper(key) {
		case "CANONICALIZE_HOST_NAME":
			return nil, fmt.Errorf("CANONICALIZE_HOST_NAME is not supported when using gssapi on %s", runtime.GOOS)
		case "SERVICE_REALM":
			return nil, fmt.Errorf("SERVICE_REALM is not supported when using gssapi on %s", runtime.GOOS)
		case "SERVICE_NAME":
			serviceName = value
		case "SERVICE_HOST":
			target = value
		default:
			return nil, fmt.Errorf("unknown mechanism property %s", key)
		}
	}

	servicePrincipalName := fmt.Sprintf("%s@%s", serviceName, target)

	return &SaslClient{
		servicePrincipalName: servicePrincipalName,
		username:             username,
		password:             password,
		passwordSet:          passwordSet,
	}, nil
}

type SaslClient struct {
	servicePrincipalName string
	username             string
	password             string
	passwordSet          bool

	// state
	state           C.gssapi_client_state
	contextComplete bool
	done            bool
}

func (sc *SaslClient) Close() {
	C.gssapi_client_destroy(&sc.state)
}

func (sc *SaslClient) Start() (string, []byte, error) {
	const mechName = "GSSAPI"

	cservicePrincipalName := C.CString(sc.servicePrincipalName)
	defer C.free(unsafe.Pointer(cservicePrincipalName))
	var cusername *C.char
	var cpassword *C.char
	if sc.username != "" {
		cusername = C.CString(sc.username)
		defer C.free(unsafe.Pointer(cusername))
		if sc.passwordSet {
			cpassword = C.CString(sc.password)
			defer C.free(unsafe.Pointer(cpassword))
		}
	}
	status := C.gssapi_client_init(&sc.state, cservicePrincipalName, cusername, cpassword)

	if status != C.GSSAPI_OK {
		return mechName, nil, sc.getError("unable to initialize client")
	}

	payload, err := sc.Next(nil)

	return mechName, payload, err
}

func (sc *SaslClient) Next(challenge []byte) ([]byte, error) {

	var buf unsafe.Pointer
	var bufLen C.size_t
	var outBuf unsafe.Pointer
	var outBufLen C.size_t

	if sc.contextComplete {
		if sc.username == "" {
			var cusername *C.char
			status := C.gssapi_client_username(&sc.state, &cusername)
			if status != C.GSSAPI_OK {
				return nil, sc.getError("unable to acquire username")
			}
			defer C.free(unsafe.Pointer(cusername))
			sc.username = C.GoString((*C.char)(unsafe.Pointer(cusername)))
		}

		bytes := append([]byte{1, 0, 0, 0}, []byte(sc.username)...)
		buf = unsafe.Pointer(&bytes[0])
		bufLen = C.size_t(len(bytes))
		status := C.gssapi_client_wrap_msg(&sc.state, buf, bufLen, &outBuf, &outBufLen)
		if status != C.GSSAPI_OK {
			return nil, sc.getError("unable to wrap authz")
		}

		sc.done = true
	} else {
		if len(challenge) > 0 {
			buf = unsafe.Pointer(&challenge[0])
			bufLen = C.size_t(len(challenge))
		}

		status := C.gssapi_client_negotiate(&sc.state, buf, bufLen, &outBuf, &outBufLen)
		switch status {
		case C.GSSAPI_OK:
			sc.contextComplete = true
		case C.GSSAPI_CONTINUE:
		default:
			return nil, sc.getError("unable to negotiate with server")
		}
	}

	if outBuf != nil {
		defer C.free(outBuf)
	}

	return C.GoBytes(outBuf, C.int(outBufLen)), nil
}

func (sc *SaslClient) Completed() bool {
	return sc.done
}

func (sc *SaslClient) getError(prefix string) error {
	var desc *C.char

	status := C.gssapi_error_desc(sc.state.maj_stat, sc.state.min_stat, &desc)
	if status != C.GSSAPI_OK {
		if desc != nil {
			C.free(unsafe.Pointer(desc))
		}

		return fmt.Errorf("%s: (%v, %v)", prefix, sc.state.maj_stat, sc.state.min_stat)
	}
	defer C.free(unsafe.Pointer(desc))

	return fmt.Errorf("%s: %v(%v,%v)", prefix, C.GoString(desc), int32(sc.state.maj_stat), int32(sc.state.min_stat))
}
