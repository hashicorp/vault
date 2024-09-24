/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
export default class KubernetesRoleCredentialsRoute extends Route {
  @service secretMountPath;

  model() {
    return {
      roleName: this.paramsFor('roles.role').name,
      backend: this.secretMountPath.currentPath,
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: resolvedModel.backend, route: 'overview' },
      { label: 'Roles', route: 'roles', model: resolvedModel.backend },
      { label: resolvedModel.roleName, route: 'roles.role.details' },
      { label: 'Credentials' },
    ];
  }
}
