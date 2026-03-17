/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class PkiRolesErrorRoute extends Route {
  @service secretMountPath;

  setupController(controller) {
    super.setupController(...arguments);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
    ];
    controller.tabs = [
      { label: 'Overview', route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'Roles', route: 'roles.index', model: this.secretMountPath.currentPath },
      { label: 'Issuers', route: 'issuers.index', model: this.secretMountPath.currentPath },
      { label: 'Keys', route: 'keys.index', model: this.secretMountPath.currentPath },
      { label: 'Certificates', route: 'certificates.index', model: this.secretMountPath.currentPath },
      { label: 'Tidy', route: 'tidy.index', model: this.secretMountPath.currentPath },
      { label: 'Configuration', route: 'configuration.index', model: this.secretMountPath.currentPath },
    ];
    controller.title = this.secretMountPath.currentPath;
  }
}
