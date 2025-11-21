/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import PkiIssuersSignIntermediateForm from 'vault/forms/secrets/pki/issuers/sign-intermediate';

export default class PkiIssuerSignRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const { issuer_ref } = this.paramsFor('issuers/issuer');
    return {
      form: new PkiIssuersSignIntermediateForm({}, { isNew: true }),
      issuerRef: issuer_ref,
    };
  }
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'Issuers', route: 'issuers.index', model: this.secretMountPath.currentPath },
      {
        label: resolvedModel.issuerRef,
        route: 'issuers.issuer.details',
        models: [this.secretMountPath.currentPath, resolvedModel.issuerRef],
      },
      { label: 'Sign Intermediate' },
    ];
  }
}
