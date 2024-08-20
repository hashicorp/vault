/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class TidyAutoIndexRoute extends Route {
  @service secretMountPath;
  @service store;

  // inherits model from tidy/auto

  setupController(controller, resolvedModel) {
    // autoTidyConfig id is the backend path
    const { id: backend } = resolvedModel;
    super.setupController(...arguments);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: backend },
      { label: 'tidy', route: 'tidy.index', model: backend },
      { label: 'auto' },
    ];
    controller.title = this.secretMountPath.currentPath;
  }
}
