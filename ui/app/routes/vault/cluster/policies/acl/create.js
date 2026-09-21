/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import UnsavedModelRoute from 'vault/mixins/unsaved-model-route';
import PolicyForm from 'vault/forms/policy';

export default Route.extend(UnsavedModelRoute, {
  api: service(),

  async model() {
    const policyType = 'acl';
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
    controller.set('policyType', 'acl');
  },
});
