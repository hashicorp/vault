/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { fetchMfaMethods } from 'vault/utils/mfa-login-enforcement-helpers';

export default class MfaMethodsRoute extends Route {
  @service router;
  @service api;

  async model() {
    try {
      const methods = await fetchMfaMethods(this.api);
      return { methods };
    } catch (err) {
      const { status } = await this.api.parseError(err);
      if (status === 404) {
        this.router.transitionTo('vault.cluster.access.mfa.index');

        return { methods: [] };
      } else {
        throw err;
      }
    }
  }

  afterModel(model) {
    if (model.length === 0) {
      this.router.transitionTo('vault.cluster.access.mfa');
    }
  }

  setupController(controller, model) {
    controller.set('model', model);

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'Multi-factor authentication' },
    ];
  }
}
