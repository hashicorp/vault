/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

@withConfirmLeave()
export default class PkiTidyAutoConfigureRoute extends Route {
  @service store;
  @service secretMountPath;

  // inherits model from tidy/auto

  setupController(controller, resolvedModel) {
    // autoTidyConfig id is the backend path
    const { id: backend } = resolvedModel;
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: backend },
      { label: 'configuration', route: 'configuration.index', model: backend },
      { label: 'tidy', route: 'tidy', model: backend },
      { label: 'auto', route: 'tidy.auto', model: backend },
      { label: 'configure' },
    ];
  }
}
