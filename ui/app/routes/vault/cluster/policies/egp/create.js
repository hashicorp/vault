/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import PolicyForm from 'vault/forms/policy';

export default Route.extend({
  router: service(),
  api: service(),
  version: service(),

  async model() {
    const policyType = 'egp';
    if (!this.version.hasSentinel) {
      return this.router.transitionTo(`vault.cluster.policies.${policyType}.create`);
    }

    const form = new PolicyForm(
      {
        enforcement_level: 'hard-mandatory',
      },
      { isNew: true }
    );
    form.policyType = policyType;
    return form;
  },

  setupController(controller) {
    this._super(...arguments);
    controller.set('policyType', 'egp');
  },
});
