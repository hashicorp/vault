/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

@withConfirmLeave()
export default class PkiIssuersGenerateIntermediateRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    return this.store.createRecord('pki/action', { actionType: 'generate-csr' });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'issuers', route: 'issuers.index', model: this.secretMountPath.currentPath },
      { label: 'generate CSR' },
    ];
  }
}
