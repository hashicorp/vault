/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default Route.extend({
  store: service(),
  secretMountPath: service(),
  credParams() {
    const { role_name: role, scope_name: scope } = this.paramsFor('credentials');
    return {
      role,
      scope,
    };
  },
  model(params) {
    const { role, scope } = this.credParams();
    return this.store.queryRecord('kmip/credential', {
      role,
      scope,
      backend: this.secretMountPath.currentPath,
      id: params.serial,
    });
  },

  setupController(controller) {
    const { role, scope } = this.credParams();
    this._super(...arguments);
    controller.setProperties({ role, scope });
  },
});
