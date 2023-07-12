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
    return hash({
      secretsEngines: this.store.query('secret-engine', {}),
      version: this.version,
      vaultConfiguration: this.getVaultConfiguration(),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const { version } = resolvedModel;

    controller.versionHeader = version.isEnterprise
      ? `Vault v${version.version.slice(0, version.version.indexOf('+'))}`
      : `Vault v${version.version}`;
  }
}
