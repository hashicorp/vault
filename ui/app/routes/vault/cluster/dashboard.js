/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class VaultClusterDashboardRoute extends Route {
  @service store;
  @service version;

  model() {
    return hash({
      secretsEngines: this.store.query('secret-engine', {}),
      version: this.version.version,
    });
  }
}
