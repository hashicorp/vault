/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class KmipRoleRoute extends Route {
  @service store;
  @service secretMountPath;
  @service pathHelp;

  beforeModel() {
    return this.pathHelp.getNewModel('kmip/role', this.secretMountPath.currentPath);
  }

  model(params) {
    return this.store.queryRecord('kmip/role', {
      backend: this.secretMountPath.currentPath,
      scope: params.scope_name,
      id: params.role_name,
    });
  }

  setupController(controller) {
    super.setupController(...arguments);
    const { scope_name: scope, role_name: role } = this.paramsFor('role');
    controller.setProperties({ role, scope });
  }
}
