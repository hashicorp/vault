/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ROUTES } from 'vault/utils/routes';

export default class KubernetesRoleDetailsRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const backend = this.secretMountPath.currentPath;
    const { name } = this.paramsFor(ROUTES.ROLES_ROLE);
    return this.store.queryRecord('kubernetes/role', { backend, name });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: resolvedModel.backend, route: ROUTES.OVERVIEW },
      { label: 'Roles', route: ROUTES.ROLES, model: resolvedModel.backend },
      { label: resolvedModel.name },
    ];
  }
}
