/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class TidyAutoIndexRoute extends Route {
  @service secretMountPath;

  // inherits model from tidy/auto

  setupController(controller, resolvedModel) {
    const { currentPath } = this.secretMountPath;
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'overview', model: currentPath },
      { label: 'Tidy', route: 'tidy.index', model: currentPath },
      { label: 'Auto' },
    ];
    controller.backend = currentPath;
  }
}
