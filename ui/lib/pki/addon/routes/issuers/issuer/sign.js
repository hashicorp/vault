/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

@withConfirmLeave()
export default class PkiIssuerSignRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const { issuer_ref } = this.paramsFor('issuers/issuer');
    return this.store.createRecord('pki/sign-intermediate', { issuerRef: issuer_ref });
  }
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'issuers', route: 'issuers.index' },
      { label: resolvedModel.issuerRef, route: 'issuers.issuer.details' },
      { label: 'sign intermediate' },
    ];
  }
}
