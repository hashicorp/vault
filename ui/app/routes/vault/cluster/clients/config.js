/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class ConfigRoute extends Route {
  @service store;

  model() {
    return this.store.queryRecord('clients/config', {});
  }
}
