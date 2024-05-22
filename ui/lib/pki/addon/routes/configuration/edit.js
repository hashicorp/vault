/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

@withConfirmLeave('model.config', ['model.urls', 'model.crl'])
export default class PkiConfigurationEditRoute extends Route {
  @service secretMountPath;

  model() {
    const { acme, cluster, urls, crl, engine } = this.modelFor('configuration');
    return {
      engineId: engine.id,
      acme,
      cluster,
      urls,
      crl,
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'configuration', route: 'configuration.index', model: this.secretMountPath.currentPath },
      { label: 'edit' },
    ];
  }
}
