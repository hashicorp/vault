/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { fetchMfaMethods } from 'vault/utils/mfa-login-enforcement-helpers';

import MfaLoginEnforcementForm from 'vault/forms/mfa/login-enforcement';

export default class MfaLoginEnforcementCreateRoute extends Route {
  @service api;

  async model() {
    const methods = await fetchMfaMethods(this.api);
    return { form: new MfaLoginEnforcementForm({}, { isNew: true }), methods };
  }
}
