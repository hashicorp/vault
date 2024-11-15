// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

//+build gssapi,windows

#include "sspi_wrapper.h"

static HINSTANCE sspi_secur32_dll = NULL;
static PSecurityFunctionTable sspi_functions = NULL;
static const LPSTR SSPI_PACKAGE_NAME = "kerberos";

int sspi_init(
)
{
	// Load the secur32.dll library using its exact path. Passing the exact DLL path rather than allowing LoadLibrary to
	// search in different locations removes the possibility of DLL preloading attacks. We use GetSystemDirectoryA and
	// LoadLibraryA rather than the GetSystemDirectory/LoadLibrary aliases to ensure the ANSI versions are used so we
	// don't have to account for variations in char sizes if UNICODE is enabled.

	// Passing a 0 size will return the required buffer length to hold the path, including the null terminator.
	int requiredLen = GetSystemDirectoryA(NULL, 0);
	if (!requiredLen) {
		return GetLastError();
	}

	// Allocate a buffer to hold the system directory + "\secur32.dll" (length 12, not including null terminator).
	int actualLen = requiredLen + 12;
	char *directoryBuffer = (char *) calloc(1, actualLen);
	int directoryLen = GetSystemDirectoryA(directoryBuffer, actualLen);
	if (!directoryLen) {
		free(directoryBuffer);
		return GetLastError();
	}

	// Append the DLL name to the buffer.
	char *dllName = "\\secur32.dll";
	strcpy_s(&(directoryBuffer[directoryLen]), actualLen - directoryLen, dllName);

	sspi_secur32_dll = LoadLibraryA(directoryBuffer);
	free(directoryBuffer);
	if (!sspi_secur32_dll) {
		return GetLastError();
	}

    INIT_SECURITY_INTERFACE init_security_interface = (INIT_SECURITY_INTERFACE)GetProcAddress(sspi_secur32_dll, SECURITY_ENTRYPOINT);
    if (!init_security_interface) {
        return -1;
    }

    sspi_functions = (*init_security_interface)();
    if (!sspi_functions) {
        return -2;
    }

	return SSPI_OK;
}

int sspi_client_init(
    sspi_client_state *client,
    char* username,
    char* password
)
{
	TimeStamp timestamp;

    if (username) {
        if (password) {
            SEC_WINNT_AUTH_IDENTITY auth_identity;

        #ifdef _UNICODE
            auth_identity.Flags = SEC_WINNT_AUTH_IDENTITY_UNICODE;
        #else
            auth_identity.Flags = SEC_WINNT_AUTH_IDENTITY_ANSI;
        #endif
            auth_identity.User = (LPSTR) username;
            auth_identity.UserLength = strlen(username);
            auth_identity.Password = (LPSTR) password;
            auth_identity.PasswordLength = strlen(password);
            auth_identity.Domain = NULL;
            auth_identity.DomainLength = 0;
            client->status = sspi_functions->AcquireCredentialsHandle(NULL, SSPI_PACKAGE_NAME, SECPKG_CRED_OUTBOUND, NULL, &auth_identity, NULL, NULL, &client->cred, &timestamp);
        } else {
            client->status = sspi_functions->AcquireCredentialsHandle(username, SSPI_PACKAGE_NAME, SECPKG_CRED_OUTBOUND, NULL, NULL, NULL, NULL, &client->cred, &timestamp);
        }
    } else {
        client->status = sspi_functions->AcquireCredentialsHandle(NULL, SSPI_PACKAGE_NAME, SECPKG_CRED_OUTBOUND, NULL, NULL, NULL, NULL, &client->cred, &timestamp);
    }

    if (client->status != SEC_E_OK) {
        return SSPI_ERROR;
    }

    return SSPI_OK;
}

int sspi_client_username(
    sspi_client_state *client,
    char** username
)
{
    SecPkgCredentials_Names names;
	client->status = sspi_functions->QueryCredentialsAttributes(&client->cred, SECPKG_CRED_ATTR_NAMES, &names);

	if (client->status != SEC_E_OK) {
		return SSPI_ERROR;
	}

	int len = strlen(names.sUserName) + 1;
	*username = malloc(len);
	memcpy(*username, names.sUserName, len);

	sspi_functions->FreeContextBuffer(names.sUserName);

    return SSPI_OK;
}

int sspi_client_negotiate(
    sspi_client_state *client,
    char* spn,
    PVOID input,
    ULONG input_length,
    PVOID* output,
    ULONG* output_length
)
{
    SecBufferDesc inbuf;
	SecBuffer in_bufs[1];
	SecBufferDesc outbuf;
	SecBuffer out_bufs[1];

	if (client->has_ctx > 0) {
		inbuf.ulVersion = SECBUFFER_VERSION;
		inbuf.cBuffers = 1;
		inbuf.pBuffers = in_bufs;
		in_bufs[0].pvBuffer = input;
		in_bufs[0].cbBuffer = input_length;
		in_bufs[0].BufferType = SECBUFFER_TOKEN;
	}

	outbuf.ulVersion = SECBUFFER_VERSION;
	outbuf.cBuffers = 1;
	outbuf.pBuffers = out_bufs;
	out_bufs[0].pvBuffer = NULL;
	out_bufs[0].cbBuffer = 0;
	out_bufs[0].BufferType = SECBUFFER_TOKEN;

	ULONG context_attr = 0;

	client->status = sspi_functions->InitializeSecurityContext(
        &client->cred,
        client->has_ctx > 0 ? &client->ctx : NULL,
        (LPSTR) spn,
        ISC_REQ_ALLOCATE_MEMORY | ISC_REQ_MUTUAL_AUTH,
        0,
        SECURITY_NETWORK_DREP,
        client->has_ctx > 0 ? &inbuf : NULL,
        0,
        &client->ctx,
        &outbuf,
        &context_attr,
        NULL);

    if (client->status != SEC_E_OK && client->status != SEC_I_CONTINUE_NEEDED) {
        return SSPI_ERROR;
    }

    client->has_ctx = 1;

	*output = malloc(out_bufs[0].cbBuffer);
	*output_length = out_bufs[0].cbBuffer;
	memcpy(*output, out_bufs[0].pvBuffer, *output_length);
    sspi_functions->FreeContextBuffer(out_bufs[0].pvBuffer);

    if (client->status == SEC_I_CONTINUE_NEEDED) {
        return SSPI_CONTINUE;
    }

    return SSPI_OK;
}

int sspi_client_wrap_msg(
    sspi_client_state *client,
    PVOID input,
    ULONG input_length,
    PVOID* output,
    ULONG* output_length
)
{
    SecPkgContext_Sizes sizes;

	client->status = sspi_functions->QueryContextAttributes(&client->ctx, SECPKG_ATTR_SIZES, &sizes);
	if (client->status != SEC_E_OK) {
		return SSPI_ERROR;
	}

	char *msg = malloc((sizes.cbSecurityTrailer + input_length + sizes.cbBlockSize) * sizeof(char));
	memcpy(&msg[sizes.cbSecurityTrailer], input, input_length);

	SecBuffer wrap_bufs[3];
	SecBufferDesc wrap_buf_desc;
	wrap_buf_desc.cBuffers = 3;
	wrap_buf_desc.pBuffers = wrap_bufs;
	wrap_buf_desc.ulVersion = SECBUFFER_VERSION;

	wrap_bufs[0].cbBuffer = sizes.cbSecurityTrailer;
	wrap_bufs[0].BufferType = SECBUFFER_TOKEN;
	wrap_bufs[0].pvBuffer = msg;

	wrap_bufs[1].cbBuffer = input_length;
	wrap_bufs[1].BufferType = SECBUFFER_DATA;
	wrap_bufs[1].pvBuffer = msg + sizes.cbSecurityTrailer;

	wrap_bufs[2].cbBuffer = sizes.cbBlockSize;
	wrap_bufs[2].BufferType = SECBUFFER_PADDING;
	wrap_bufs[2].pvBuffer = msg + sizes.cbSecurityTrailer + input_length;

	client->status = sspi_functions->EncryptMessage(&client->ctx, SECQOP_WRAP_NO_ENCRYPT, &wrap_buf_desc, 0);
	if (client->status != SEC_E_OK) {
		free(msg);
		return SSPI_ERROR;
	}

	*output_length = wrap_bufs[0].cbBuffer + wrap_bufs[1].cbBuffer + wrap_bufs[2].cbBuffer;
	*output = malloc(*output_length);

	memcpy(*output, wrap_bufs[0].pvBuffer, wrap_bufs[0].cbBuffer);
	memcpy(*output + wrap_bufs[0].cbBuffer, wrap_bufs[1].pvBuffer, wrap_bufs[1].cbBuffer);
	memcpy(*output + wrap_bufs[0].cbBuffer + wrap_bufs[1].cbBuffer, wrap_bufs[2].pvBuffer, wrap_bufs[2].cbBuffer);

	free(msg);

	return SSPI_OK;
}

int sspi_client_destroy(
    sspi_client_state *client
)
{
    if (client->has_ctx > 0) {
        sspi_functions->DeleteSecurityContext(&client->ctx);
    }

    sspi_functions->FreeCredentialsHandle(&client->cred);

    return SSPI_OK;
}
