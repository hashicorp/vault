/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class PkiKeyDetailsRoute extends Route {
  @service secretMountPath;

  model() {
    return this.modelFor('keys.key');
  }
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: resolvedModel.backend },
      { label: 'Keys', route: 'keys.index', model: resolvedModel.backend },
      { label: resolvedModel.id },
    ];
  }
}
