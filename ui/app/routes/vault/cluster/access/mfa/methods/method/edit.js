/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import MfaCreateTotpMethodForm from 'vault/forms/mfa/method/totp';
import MfaCreateDuoMethodForm from 'vault/forms/mfa/method/duo';
import MfaCreateOktaMethodForm from 'vault/forms/mfa/method/okta';
import MfaCreatePingIdMethodForm from 'vault/forms/mfa/method/ping-id';

export default class MfaMethodEditRoute extends Route {
  @service api;

  async model() {
    const { method } = this.modelFor('vault.cluster.access.mfa.methods.method');

    let form;

    if (method.type === 'totp') {
      form = new MfaCreateTotpMethodForm(method, { isNew: false });
    } else if (method.type === 'duo') {
      form = new MfaCreateDuoMethodForm(method, { isNew: false });
    } else if (method.type === 'okta') {
      form = new MfaCreateOktaMethodForm(method, { isNew: false });
    } else if (method.type === 'pingid') {
      form = new MfaCreatePingIdMethodForm(method, { isNew: false });
    }

    return {
      form,
      method,
    };
  }
}
