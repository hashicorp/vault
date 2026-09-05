/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { prepareTargets, fetchMfaMethods } from 'vault/utils/mfa-login-enforcement-helpers';

export default class MfaLoginEnforcementRoute extends Route {
  @service api;

  async model({ name }) {
    const { data } = await this.api.identity.mfaReadLoginEnforcement(name);
    const targets = await prepareTargets(data, this.api);
    const methods = await fetchMfaMethods(this.api);

    return { enforcement: data, name, targets, methods };
  }
}
