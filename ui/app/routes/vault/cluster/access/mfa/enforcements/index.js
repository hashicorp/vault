/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { fetchMfaLoginEnforcements } from 'vault/utils/mfa-login-enforcement-helpers';

export default class MfaEnforcementsRoute extends Route {
  @service api;

  async model() {
    try {
      const enforcements = await fetchMfaLoginEnforcements(this.api);
      return { enforcements };
    } catch (err) {
      const { status } = await this.api.parseError(err);
      if (status === 404) {
        return { enforcements: [] };
      } else {
        throw err;
      }
    }
  }

  setupController(controller, model) {
    controller.set('model', model);

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'Multi-factor authentication', route: 'vault.cluster.access.mfa' },
      { label: 'Enforcements' },
    ];
  }
}
