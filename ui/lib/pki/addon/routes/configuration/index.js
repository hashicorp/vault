/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { withConfig } from 'pki/decorators/check-issuers';
import { PKI_DEFAULT_EMPTY_STATE_MSG } from 'pki/routes/overview';
import { service } from '@ember/service';

@withConfig()
export default class ConfigurationIndexRoute extends Route {
  @service('app-router') router;

  async model() {
    return {
      hasConfig: this.pkiMountHasConfig,
      mountConfig: this.modelFor('application'),
      ...this.modelFor('configuration'),
    };
  }

  afterModel(resolvedModel) {
    if (!resolvedModel.hasConfig) {
      this.router.transitionTo('vault.cluster.secrets.backend.pki.configuration.create');
    }
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.notConfiguredMessage = PKI_DEFAULT_EMPTY_STATE_MSG;
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.engine.id, route: 'overview', model: resolvedModel.engine.id },
      { label: 'Configuration' },
    ];
  }
}
