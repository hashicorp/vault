/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';

export default class SyncSecretsDestinationsIndexController extends Controller {
  queryParams = ['name', 'type', 'page'];
}
