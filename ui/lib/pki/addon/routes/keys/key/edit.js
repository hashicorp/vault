/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { withConfirmLeave } from 'core/decorators/confirm-leave';
import Route from '@ember/routing/route';
import { service } from '@ember/service';

@withConfirmLeave()
export default class PkiKeyEditRoute extends Route {
  @service secretMountPath;

  model() {
    return this.modelFor('keys.key');
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'Keys', route: 'keys.index', model: this.secretMountPath.currentPath },
      { label: resolvedModel.id },
    ];
  }
}
