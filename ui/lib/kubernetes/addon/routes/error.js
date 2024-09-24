/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class KubernetesErrorRoute extends Route {
  @service secretMountPath;

  setupController(controller) {
    super.setupController(...arguments);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
    ];
    controller.backend = this.modelFor('application');
  }
}
