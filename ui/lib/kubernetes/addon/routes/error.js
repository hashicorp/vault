/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ROUTES } from 'vault/utils/routes';

export default class KubernetesErrorRoute extends Route {
  @service secretMountPath;

  setupController(controller) {
    super.setupController(...arguments);
    controller.breadcrumbs = [
      { label: 'Secrets', route: ROUTES.SECRETS, linkExternal: true },
      { label: this.secretMountPath.currentPath, route: ROUTES.OVERVIEW },
    ];
    controller.backend = this.modelFor('application');
  }
}
