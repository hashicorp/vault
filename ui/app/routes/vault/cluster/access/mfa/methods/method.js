/**
 * Copyright IBM Corp. 2016, 2026
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
    let enforcements = [];
    let method;

    try {
      method = await this.api.identity.mfaReadMethod(id);
      if (method?.data) {
        method.data.displayName = getMfaMethodName(method.data.type);
        method.data.icon = getMfaMethodIcon(method.data.type);
      }
    } catch (err) {
      const { status } = await this.api.parseError(err);
      if (status !== 404) {
        throw err;
      }
    }

    try {
      const allEnforcements = await fetchMfaLoginEnforcements(this.api);
      enforcements = allEnforcements.filter((item) => {
        return item.mfa_method_ids?.includes(id);
      });
    } catch (err) {
      // Ignore errors loading enforcements
    }

    return hash({
      method: method?.data || {},
      enforcements: enforcements || [],
    });
  }

  setupController(controller, model) {
    controller.set('model', model);
  }
}
