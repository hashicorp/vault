/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

@withConfirmLeave()
export default class PkiTidyManualRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    return this.store.createRecord('pki/tidy', { backend: this.secretMountPath.currentPath });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: resolvedModel.backend },
      { label: 'Configuration', route: 'configuration.index', model: resolvedModel.backend },
      { label: 'Tidy', route: 'tidy', model: resolvedModel.backend },
      { label: 'Manual' },
    ];
  }
}
