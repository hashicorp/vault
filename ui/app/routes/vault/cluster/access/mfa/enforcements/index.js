/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class MfaEnforcementsRoute extends Route {
  @service store;

  model() {
    return this.store.query('mfa-login-enforcement', {}).catch((err) => {
      if (err.httpStatus === 404) {
        return [];
      } else {
        throw err;
      }
    });
  }
  setupController(controller, model) {
    controller.set('model', model);
  }
}
