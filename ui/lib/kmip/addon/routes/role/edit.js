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
  beforeModel() {
    return this.pathHelp.getNewModel('kmip/role', this.secretMountPath.currentPath);
  },
  model() {
    const params = this.paramsFor(this.routeName);
    return this.store.queryRecord('kmip/role', {
      backend: this.secretMountPath.currentPath,
      scope: params.scope_name,
      id: params.role_name,
    });
  },

  setupController(controller) {
    this._super(...arguments);
    const { scope_name: scope, role_name: role } = this.paramsFor(this.routeName);
    controller.setProperties({ role, scope });
  },
});
