/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { ResponseError } from '@hashicorp/vault-client-typescript';

export const getErrorResponse = <T>(error?: T, status?: number) => {
  const e = error || {
    errors: ['first error', 'second error'],
    message: 'there were some errors',
  };
  // url is readonly on Response so mock it and cast to Response type
  return new ResponseError({
    status: status || 404,
    url: `${document.location.origin}/v1/test/error/parsing`,
    json: () => Promise.resolve(e),
  } as Response);
};
