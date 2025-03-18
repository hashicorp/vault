/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ROUTES } from 'vault/utils/routes';

export default class MfaMethodsRoute extends Route {
  @service store;
  @service router;

  model() {
    return this.store.query('mfa-method', {}).catch((err) => {
      if (err.httpStatus === 404) {
        return [];
      } else {
        throw err;
      }
    });
  }

  afterModel(model) {
    if (model.length === 0) {
      this.router.transitionTo(ROUTES.VAULT_CLUSTER_ACCESS_MFA);
    }
  }
  setupController(controller, model) {
    controller.set('model', model);
  }
}
