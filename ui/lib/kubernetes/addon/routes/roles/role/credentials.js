/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
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
      { label: 'roles', route: 'roles' },
      { label: resolvedModel.roleName, route: 'roles.role.details' },
      { label: 'credentials' },
    ];
  }
}
