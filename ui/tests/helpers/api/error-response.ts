/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { ResponseError } from '@hashicorp/vault-client-typescript';

const DEFAULT_ERROR = {
  errors: ['first error', 'second error'],
  message: 'there were some errors',
};

const EMPTY_ERROR = { errors: [] };

export const getErrorResponse = (error = DEFAULT_ERROR, status = 404) => {
  // 404 responses do not return any errors
  const e = status === 404 ? EMPTY_ERROR : error;

  // url is readonly on Response so mock it and cast to Response type
  const response = {
    status,
    url: `${document.location.origin}/v1/test/error/parsing`,
    json: () => Promise.resolve(e),
  } as Response;

  return new ResponseError({
    ...response,
    clone: () => response,
  });
};
