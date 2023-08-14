/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
// eslint-disable-next-line ember/no-mixins
import ClusterRoute from 'vault/mixins/cluster-route';

export default class VaultClusterDashboardRoute extends Route.extend(ClusterRoute) {
  @service store;
  @service version;

  async getVaultConfiguration() {
    try {
      const adapter = this.store.adapterFor('application');
      const configState = await adapter.ajax('/v1/sys/config/state/sanitized', 'GET');
      return configState.data;
    } catch (e) {
      return null;
    }
  }

  async model() {
    const vaultConfiguration = this.getVaultConfiguration();

    return hash({
      vaultConfiguration,
      secretsEngines: this.store.query('secret-engine', {}),
      version: this.version,
    });
  }
}
