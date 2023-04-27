/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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
      { label: 'Certificates', route: 'certificates.index' },
      { label: 'Keys', route: 'keys.index' },
    ];
    controller.title = this.secretMountPath.currentPath;
  }
}
