/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { hash } from 'rsvp';
import { service } from '@ember/service';

export default class MfaMethodRoute extends Route {
  @service store;

  model({ id }) {
    return hash({
      method: this.store.findRecord('mfa-method', id).then((data) => data),
      enforcements: this.store
        .query('mfa-login-enforcement', {})
        .then((data) => {
          const filteredEnforcements = data.filter((item) => {
            const results = item.hasMany('mfa_methods').ids();
            return results.includes(id);
          });
          return filteredEnforcements;
        })
        .catch(() => {
          // Do nothing
        }),
    });
  }
  setupController(controller, model) {
    controller.set('model', model);
  }
}
