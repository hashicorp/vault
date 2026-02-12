/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import PkiIssuerForm from 'vault/forms/secrets/pki/issuers/issuer';

export default class PkiIssuerEditRoute extends Route {
  @service api;
  @service secretMountPath;

  async model() {
    const { issuer_ref } = this.paramsFor('issuers/issuer');
    const issuer = await this.api.secrets.pkiReadIssuer(issuer_ref, this.secretMountPath.currentPath);
    return {
      form: new PkiIssuerForm(issuer),
      issuerRef: issuer_ref,
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'Issuers', route: 'issuers.index', model: this.secretMountPath.currentPath },
      {
        label: resolvedModel.issuerRef,
        route: 'issuers.issuer.details',
        models: [this.secretMountPath.currentPath, resolvedModel.issuerRef],
      },
      { label: 'Update' },
    ];
  }
}
