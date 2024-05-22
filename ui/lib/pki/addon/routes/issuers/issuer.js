/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class PkiIssuerIndexRoute extends Route {
  @service store;
  @service secretMountPath;

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
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'issuers', route: 'issuers.index', model: this.secretMountPath.currentPath },
    ];
  }
}
