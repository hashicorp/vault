/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class KubernetesRoleEditRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const backend = this.secretMountPath.currentPath;
    const { name } = this.paramsFor('roles.role');
    return this.store.queryRecord('kubernetes/role', { backend, name });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: resolvedModel.backend, route: 'overview' },
      { label: 'Roles', route: 'roles', model: resolvedModel.backend },
      { label: resolvedModel.name, route: 'roles.role' },
      { label: 'Edit' },
    ];
  }
}
