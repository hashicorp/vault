/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';

export default class KvSecretDetailsController extends Controller {
  queryParams = ['version'];
  version = null;
}
