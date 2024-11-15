// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

//+build gssapi
//+build linux darwin
#ifndef GSS_WRAPPER_H
#define GSS_WRAPPER_H

#include <stdlib.h>
#ifdef GOOS_linux
#include <gssapi/gssapi.h>
#include <gssapi/gssapi_krb5.h>
#endif
#ifdef GOOS_darwin
#include <GSS/GSS.h>
#endif

#define GSSAPI_OK 0
#define GSSAPI_CONTINUE 1
#define GSSAPI_ERROR 2

typedef struct {
    gss_name_t spn;
    gss_cred_id_t cred;
    gss_ctx_id_t ctx;

    OM_uint32 maj_stat;
    OM_uint32 min_stat;
} gssapi_client_state;

int gssapi_error_desc(
    OM_uint32 maj_stat,
    OM_uint32 min_stat,
    char **desc
);

int gssapi_client_init(
    gssapi_client_state *client,
    char* spn,
    char* username,
    char* password
);

int gssapi_client_username(
    gssapi_client_state *client,
    char** username
);

int gssapi_client_negotiate(
    gssapi_client_state *client,
    void* input,
    size_t input_length,
    void** output,
    size_t* output_length
);

int gssapi_client_wrap_msg(
    gssapi_client_state *client,
    void* input,
    size_t input_length,
    void** output,
    size_t* output_length
);

int gssapi_client_destroy(
    gssapi_client_state *client
);

#endif
