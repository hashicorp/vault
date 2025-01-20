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
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'keys', route: 'keys.index' },
      { label: resolvedModel.id },
    ];
  }
}
