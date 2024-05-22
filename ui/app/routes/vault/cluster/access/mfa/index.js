/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class MfaConfigureRoute extends Route {
  @service router;
  @service store;

  beforeModel() {
    return this.store
      .query('mfa-method', {})
      .then(() => {
        // if response then they should transition to the methods page instead of staying on the configure page.
        this.router.transitionTo('vault.cluster.access.mfa.methods.index');
      })
      .catch(() => {
        // stay on the landing page
      });
  }
}
