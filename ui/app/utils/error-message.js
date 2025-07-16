/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * deprecated in favor of api-error-message for use with the api service
 * this util will be removed once Ember Data is fully removed
 * new requests should be made with api service
 * if Ember Data is still needed during the migration then this util may be used
 */

// accepts an error and returns error.errors joined with a comma, error.message or a fallback message
export default function (error, fallbackMessage = 'An error occurred, please try again') {
  if (error instanceof Error && error?.errors && typeof error.errors[0] === 'string') {
    return error.errors.join(', ');
  }
  return error?.message || fallbackMessage;
}
