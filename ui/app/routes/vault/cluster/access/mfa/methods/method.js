/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { hash } from 'rsvp';
import { service } from '@ember/service';
import {
  fetchMfaLoginEnforcements,
  getMfaMethodName,
  getMfaMethodIcon,
} from 'vault/utils/mfa-login-enforcement-helpers';

export default class MfaMethodRoute extends Route {
  @service api;

  async model({ id }) {
    let enforcements;
    let method;
    let error = [];

    try {
      method = await this.api.identity.mfaReadMethod(id);
      method.data.displayName = getMfaMethodName(method.data.type);
      method.data.icon = getMfaMethodIcon(method.data.type);

      enforcements = await fetchMfaLoginEnforcements(this.api);
    } catch (err) {
      const { status } = await this.api.parseError(err);
      if (status === 404) {
        error = [];
      } else {
        throw err;
      }
    }
    return hash({
      method: method.data || {},
      enforcements: enforcements || error,
    });
  }

  setupController(controller, model) {
    controller.set('model', model);
  }
}
