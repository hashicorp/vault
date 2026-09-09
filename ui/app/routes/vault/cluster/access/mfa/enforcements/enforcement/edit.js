/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { fetchMfaMethods } from 'vault/utils/mfa-login-enforcement-helpers';

import MfaLoginEnforcementForm from 'vault/forms/mfa/login-enforcement';

export default class MfaLoginEnforcementEditRoute extends Route {
  @service api;

  async model() {
    const methods = await fetchMfaMethods(this.api);
    const { enforcement } = this.modelFor('vault.cluster.access.mfa.enforcements.enforcement');
    const selectedMethods = methods.filter((method) =>
      (enforcement.mfa_method_ids || []).includes(method.id)
    );
    const formData = { ...enforcement, mfa_methods: selectedMethods };
    return { form: new MfaLoginEnforcementForm(formData, { isNew: false }), methods, name: enforcement.name };
  }
}
