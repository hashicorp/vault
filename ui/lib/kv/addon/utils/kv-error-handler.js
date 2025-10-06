/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export function kvErrorHandler(status, errorResponse) {
  // if it's a legitimate error - throw it!
  if (errorResponse?.isControlGroupError) {
    throw errorResponse;
  }

  if (typeof errorResponse === 'object' && errorResponse !== null) {
    const { data } = errorResponse;

    if (status === 403) {
      return {
        failReadErrorCode: 403,
      };
    }

    // in the case of a deleted/destroyed secret the API returns a 404 because { data: null }
    // however, there could be a metadata block with important information like deletion_time
    // handleResponse below checks 404 status codes for metadata and updates the code to 200 if it exists.
    // we still end up in the good ol' catch() block, but instead of a 404 adapter error we've "caught"
    // the metadata that sneakily tried to hide from us
    if (data) {
      return data;
    }
  }

  // if we get here, it's likely either a script error or 404 because it doesn't exist
  throw errorResponse;
}
