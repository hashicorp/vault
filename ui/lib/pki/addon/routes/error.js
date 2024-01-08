/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class PkiRolesErrorRoute extends Route {
  @service secretMountPath;

  setupController(controller) {
    super.setupController(...arguments);
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
    ];
    controller.tabs = [
      { label: 'Overview', route: 'overview' },
      { label: 'Roles', route: 'roles.index' },
      { label: 'Issuers', route: 'issuers.index' },
      { label: 'Keys', route: 'keys.index' },
      { label: 'Certificates', route: 'certificates.index' },
      { label: 'Tidy', route: 'tidy.index' },
      { label: 'Configuration', route: 'configuration.index' },
    ];
    controller.title = this.secretMountPath.currentPath;
  }
}
