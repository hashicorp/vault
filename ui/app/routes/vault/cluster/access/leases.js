/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';

export default class LeasesRoute extends Route {
  @service capabilities;

  model() {
    return this.capabilities.fetchPathCapabilities('sys/leases/lookup/');
  }
}
