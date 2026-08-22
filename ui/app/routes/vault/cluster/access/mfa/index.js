/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class MfaConfigureRoute extends Route {
  @service router;
  @service api;

  async beforeModel() {
    try {
      // if response then they should transition to the methods page instead of staying on the configure page.
      await this.api.identity.mfaListMethods(true);
      this.router.transitionTo('vault.cluster.access.mfa.methods.index');
    } catch (e) {
      // stay on the landing page
    }
  }
}
