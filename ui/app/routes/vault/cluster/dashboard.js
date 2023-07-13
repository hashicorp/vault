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
  @service namespace;

  async getVaultConfiguration() {
    try {
      const adapter = this.store.adapterFor('application');
      const configState = await adapter.ajax('/v1/sys/config/state/sanitized', 'GET');
      return configState.data;
    } catch (e) {
      return e.httpStatus;
    }
  }

  model() {
    const versionHeader = this.version.isEnterprise
      ? `Vault v${this.version.version.slice(0, this.version.version.indexOf('+'))}`
      : `Vault v${this.version.version}`;
    const vaultConfiguration = this.getVaultConfiguration();

    return hash({
      versionHeader,
      secretsEngines: this.store.query('secret-engine', {}),
      vaultConfiguration: typeof vaultConfiguration === 'number' ? vaultConfiguration : false,
      namespace: this.namespace,
    });
  }
}
