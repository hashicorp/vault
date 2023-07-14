/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class VaultClusterDashboardRoute extends Route {
  @service store;

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
    const vaultConfiguration = this.getVaultConfiguration();

    return hash({
      vaultConfiguration,
      secretsEngines: this.store.query('secret-engine', {}),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(...arguments);
    controller.vaultConfiguration =
      typeof resolvedModel.vaultConfiguration === 'number' ? false : resolvedModel.vaultConfiguration;
  }
}
