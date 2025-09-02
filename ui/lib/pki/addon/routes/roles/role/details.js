/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class RolesRoleDetailsRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const { role } = this.paramsFor('roles/role');
    return this.store.queryRecord('pki/role', {
      backend: this.secretMountPath.currentPath,
      id: role,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const { id } = resolvedModel;
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'Roles', route: 'roles.index', model: this.secretMountPath.currentPath },
      { label: id },
    ];
  }
}
