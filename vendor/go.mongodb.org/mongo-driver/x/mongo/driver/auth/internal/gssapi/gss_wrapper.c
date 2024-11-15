// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

//+build gssapi
//+build linux darwin

#include <string.h>
#include <stdio.h>
#include "gss_wrapper.h"

OM_uint32 gssapi_canonicalize_name(
    OM_uint32* minor_status,
    char *input_name,
    gss_OID input_name_type,
    gss_name_t *output_name
)
{
    OM_uint32 major_status;
    gss_name_t imported_name = GSS_C_NO_NAME;
    gss_buffer_desc buffer = GSS_C_EMPTY_BUFFER;

    buffer.value = input_name;
    buffer.length = strlen(input_name);
    major_status = gss_import_name(minor_status, &buffer, input_name_type, &imported_name);
    if (GSS_ERROR(major_status)) {
        return major_status;
    }

    major_status = gss_canonicalize_name(minor_status, imported_name, (gss_OID)gss_mech_krb5, output_name);
    if (imported_name != GSS_C_NO_NAME) {
        OM_uint32 ignored;
        gss_release_name(&ignored, &imported_name);
    }

    return major_status;
}

int gssapi_error_desc(
    OM_uint32 maj_stat,
    OM_uint32 min_stat,
    char **desc
)
{
    OM_uint32 stat = maj_stat;
    int stat_type = GSS_C_GSS_CODE;
    if (min_stat != 0) {
        stat = min_stat;
        stat_type = GSS_C_MECH_CODE;
    }

    OM_uint32 local_maj_stat, local_min_stat;
    OM_uint32 msg_ctx = 0;
    gss_buffer_desc desc_buffer;
    do
    {
        local_maj_stat = gss_display_status(
            &local_min_stat,
            stat,
            stat_type,
            GSS_C_NO_OID,
            &msg_ctx,
            &desc_buffer
        );
        if (GSS_ERROR(local_maj_stat)) {
            return GSSAPI_ERROR;
        }

        if (*desc) {
            free(*desc);
        }

        *desc = malloc(desc_buffer.length+1);
        memcpy(*desc, desc_buffer.value, desc_buffer.length+1);

        gss_release_buffer(&local_min_stat, &desc_buffer);
    }
    while(msg_ctx != 0);

    return GSSAPI_OK;
}

int gssapi_client_init(
    gssapi_client_state *client,
    char* spn,
    char* username,
    char* password
)
{
    client->cred = GSS_C_NO_CREDENTIAL;
    client->ctx = GSS_C_NO_CONTEXT;

    client->maj_stat = gssapi_canonicalize_name(&client->min_stat, spn, GSS_C_NT_HOSTBASED_SERVICE, &client->spn);
    if (GSS_ERROR(client->maj_stat)) {
        return GSSAPI_ERROR;
    }

    if (username) {
        gss_name_t name;
        client->maj_stat = gssapi_canonicalize_name(&client->min_stat, username, GSS_C_NT_USER_NAME, &name);
        if (GSS_ERROR(client->maj_stat)) {
            return GSSAPI_ERROR;
        }

        if (password) {
            gss_buffer_desc password_buffer;
            password_buffer.value = password;
            password_buffer.length = strlen(password);
            client->maj_stat = gss_acquire_cred_with_password(&client->min_stat, name, &password_buffer, GSS_C_INDEFINITE, GSS_C_NO_OID_SET, GSS_C_INITIATE, &client->cred, NULL, NULL);
        } else {
            client->maj_stat = gss_acquire_cred(&client->min_stat, name, GSS_C_INDEFINITE, GSS_C_NO_OID_SET, GSS_C_INITIATE, &client->cred, NULL, NULL);
        }

        if (GSS_ERROR(client->maj_stat)) {
            return GSSAPI_ERROR;
        }

        OM_uint32 ignored;
        gss_release_name(&ignored, &name);
    }

    return GSSAPI_OK;
}

int gssapi_client_username(
    gssapi_client_state *client,
    char** username
)
{
    OM_uint32 ignored;
    gss_name_t name = GSS_C_NO_NAME;

    client->maj_stat = gss_inquire_context(&client->min_stat, client->ctx, &name, NULL, NULL, NULL, NULL, NULL, NULL);
    if (GSS_ERROR(client->maj_stat)) {
        return GSSAPI_ERROR;
    }

    gss_buffer_desc name_buffer;
    client->maj_stat = gss_display_name(&client->min_stat, name, &name_buffer, NULL);
    if (GSS_ERROR(client->maj_stat)) {
        gss_release_name(&ignored, &name);
        return GSSAPI_ERROR;
    }

	*username = malloc(name_buffer.length+1);
	memcpy(*username, name_buffer.value, name_buffer.length+1);

    gss_release_buffer(&ignored, &name_buffer);
    gss_release_name(&ignored, &name);
    return GSSAPI_OK;
}

int gssapi_client_negotiate(
    gssapi_client_state *client,
    void* input,
    size_t input_length,
    void** output,
    size_t* output_length
)
{
    gss_buffer_desc input_buffer = GSS_C_EMPTY_BUFFER;
    gss_buffer_desc output_buffer = GSS_C_EMPTY_BUFFER;

    if (input) {
        input_buffer.value = input;
        input_buffer.length = input_length;
    }

    client->maj_stat = gss_init_sec_context(
        &client->min_stat,
        client->cred,
        &client->ctx,
        client->spn,
        GSS_C_NO_OID,
        GSS_C_MUTUAL_FLAG | GSS_C_SEQUENCE_FLAG,
        0,
        GSS_C_NO_CHANNEL_BINDINGS,
        &input_buffer,
        NULL,
        &output_buffer,
        NULL,
        NULL
    );

    if (output_buffer.length) {
        *output = malloc(output_buffer.length);
        *output_length = output_buffer.length;
        memcpy(*output, output_buffer.value, output_buffer.length);

        OM_uint32 ignored;
        gss_release_buffer(&ignored, &output_buffer);
    }

    if (GSS_ERROR(client->maj_stat)) {
        return GSSAPI_ERROR;
    } else if (client->maj_stat == GSS_S_CONTINUE_NEEDED) {
        return GSSAPI_CONTINUE;
    }

    return GSSAPI_OK;
}

int gssapi_client_wrap_msg(
    gssapi_client_state *client,
    void* input,
    size_t input_length,
    void** output,
    size_t* output_length
)
{
    gss_buffer_desc input_buffer = GSS_C_EMPTY_BUFFER;
    gss_buffer_desc output_buffer = GSS_C_EMPTY_BUFFER;

    input_buffer.value = input;
    input_buffer.length = input_length;

    client->maj_stat = gss_wrap(&client->min_stat, client->ctx, 0, GSS_C_QOP_DEFAULT, &input_buffer, NULL, &output_buffer);

    if (output_buffer.length) {
        *output = malloc(output_buffer.length);
        *output_length = output_buffer.length;
        memcpy(*output, output_buffer.value, output_buffer.length);

        gss_release_buffer(&client->min_stat, &output_buffer);
    }

    if (GSS_ERROR(client->maj_stat)) {
        return GSSAPI_ERROR;
    }

    return GSSAPI_OK;
}

int gssapi_client_destroy(
    gssapi_client_state *client
)
{
    OM_uint32 ignored;
    if (client->ctx != GSS_C_NO_CONTEXT) {
        gss_delete_sec_context(&ignored, &client->ctx, GSS_C_NO_BUFFER);
    }

    if (client->spn != GSS_C_NO_NAME) {
        gss_release_name(&ignored, &client->spn);
    }

    if (client->cred != GSS_C_NO_CREDENTIAL) {
        gss_release_cred(&ignored, &client->cred);
    }

    return GSSAPI_OK;
}
