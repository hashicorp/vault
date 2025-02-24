/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Error from 'ember-data/adapter/errors';

export default class AdapterError extends Error {
  httpStatus: number;
  path: string;
  message: string;
  errors: Array<string | { [key: string]: unknown; title?: string; message?: string }>;
  data?: {
    [key: string]: unknown;
    error?: string;
  };
}
