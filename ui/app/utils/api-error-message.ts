/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * this util was derived from error-message and updated to handle the error context returned from the api service
 * once Ember Data is fully removed, the error-message util will also be removed
 * for all requests made with the api service, use this util to display error messages from server
 */

import { ErrorContext, ApiError } from 'vault/api';
import ENV from 'vault/config/environment';

// accepts an error and returns error.errors joined with a comma, error.message or a fallback message
export default async function (error: unknown, fallbackMessage = 'An error occurred, please try again') {
  const messageOrFallback = (message?: string) => message || fallbackMessage;

  // log out the error for ease of debugging in dev env
  if (ENV.environment === 'development') {
    console.error('API Error:', error); // eslint-disable-line no-console
  }

  if ((error as ErrorContext).response instanceof Response) {
    const apiError: ApiError = await (error as ErrorContext).response?.json();

    if (apiError.errors && typeof apiError.errors[0] === 'string') {
      return apiError.errors.join(', ');
    }
    return messageOrFallback(apiError.message);
  }

  return messageOrFallback((error as Error)?.message);
}
