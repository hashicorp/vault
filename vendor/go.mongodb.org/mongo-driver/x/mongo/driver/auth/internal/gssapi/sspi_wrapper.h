//+build gssapi,windows

#ifndef SSPI_WRAPPER_H
#define SSPI_WRAPPER_H

#define SECURITY_WIN32 1  /* Required for SSPI */

#include <windows.h>
#include <sspi.h>

#define SSPI_OK 0
#define SSPI_CONTINUE 1
#define SSPI_ERROR 2

typedef struct {
    CredHandle cred;
    CtxtHandle ctx;

    int has_ctx;

    SECURITY_STATUS status;
} sspi_client_state;

int sspi_init();

int sspi_client_init(
    sspi_client_state *client,
    char* username,
    char* password
);

int sspi_client_username(
    sspi_client_state *client,
    char** username
);

int sspi_client_negotiate(
    sspi_client_state *client,
    char* spn,
    PVOID input,
    ULONG input_length,
    PVOID* output,
    ULONG* output_length
);

int sspi_client_wrap_msg(
    sspi_client_state *client,
    PVOID input,
    ULONG input_length,
    PVOID* output,
    ULONG* output_length 
);

int sspi_client_destroy(
    sspi_client_state *client
);

#endif