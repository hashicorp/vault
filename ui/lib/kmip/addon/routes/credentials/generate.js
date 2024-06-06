/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default Route.extend({
  store: service(),
  secretMountPath: service(),
  pathHelp: service(),
  model() {
    const params = this.paramsFor('credentials');
    return this.store.createRecord('kmip/credential', {
      backend: this.secretMountPath.currentPath,
      scope: params.scope_name,
      role: params.role_name,
    });
  },

  setupController(controller) {
    this._super(...arguments);
    const { scope_name: scope, role_name: role } = this.paramsFor('credentials');
    controller.setProperties({ role, scope });
  },
});
