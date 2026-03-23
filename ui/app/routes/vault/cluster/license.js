/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';

export default class LicenseRoute extends Route {
  @service store;
  @service version;
  @service router;

  beforeModel() {
    if (this.version.isCommunity) {
      this.router.transitionTo('vault.cluster');
    }
  }

  model() {
    return this.store.queryRecord('license', {});
  }
}
