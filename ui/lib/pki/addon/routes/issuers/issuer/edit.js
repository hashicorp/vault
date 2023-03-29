/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

@withConfirmLeave()
export default class PkiIssuerDetailRoute extends Route {
  @service store;
  @service secretMountPath;
  @service pathHelp;

  beforeModel() {
    // Must call this promise before the model hook otherwise it doesn't add OpenApi to record.
    return this.pathHelp.getNewModel('pki/issuer', this.secretMountPath.currentPath);
  }

  model() {
    const { issuer_ref } = this.paramsFor('issuers/issuer');
    return this.store.queryRecord('pki/issuer', {
      backend: this.secretMountPath.currentPath,
      id: issuer_ref,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'issuers', route: 'issuers.index' },
      { label: resolvedModel.id, route: 'issuers.issuer.details' },
      { label: 'update' },
    ];
  }
}
