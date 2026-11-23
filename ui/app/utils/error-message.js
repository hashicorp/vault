/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * deprecated in favor of this.api.parseError for use with errors from the api service
 * new requests should be made with the api service
 */

// accepts an error and returns error.errors joined with a comma, error.message or a fallback message
export default function (error, fallbackMessage = 'An error occurred, please try again') {
  if (error instanceof Error && error?.errors && typeof error.errors[0] === 'string') {
    return error.errors.join(', ');
  }
  return error?.message || fallbackMessage;
}
