/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

@withConfirmLeave()
export default class PkiIssuerCrossSignRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    return this.modelFor('issuers.issuer');
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'Issuers', route: 'issuers.index', model: this.secretMountPath.currentPath },
      {
        label: resolvedModel.issuer_id,
        route: 'issuers.issuer.details',
        models: [this.secretMountPath.currentPath, resolvedModel.issuer_id],
      },
      { label: 'Cross-sign' },
    ];
  }
}
