/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ROUTES } from 'vault/utils/routes';
export default class KubernetesRoleCredentialsRoute extends Route {
  @service secretMountPath;

  model() {
    return {
      roleName: this.paramsFor(ROUTES.ROLES_ROLE).name,
      backend: this.secretMountPath.currentPath,
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: resolvedModel.backend, route: ROUTES.OVERVIEW },
      { label: 'Roles', route: ROUTES.ROLES, model: resolvedModel.backend },
      { label: resolvedModel.roleName, route: ROUTES.ROLES_ROLE_DETAILS },
      { label: 'Credentials' },
    ];
  }
}
