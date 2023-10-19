/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';

export default class OpenApiIndex extends Controller {
  queryParams = ['filter'];

  filter = null;
}
