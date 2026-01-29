/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';

export default class LeasesRoute extends Route {
  @service store;

  model() {
    return this.store.findRecord('capabilities', 'sys/leases/lookup/');
  }
}
