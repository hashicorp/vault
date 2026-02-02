/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { verifyCertificates, parseCertificate } from 'vault/utils/parse-pki-cert';

export default class PkiIssuerIndexRoute extends Route {
  @service api;
  @service secretMountPath;

  model() {
    const { issuer_ref } = this.paramsFor('issuers/issuer');
    return this.api.secrets
      .pkiReadIssuer(issuer_ref, this.secretMountPath.currentPath)
      .then(async (issuer) => {
        const isRoot = await verifyCertificates(issuer.certificate, issuer.certificate);
        const parsedCertificate = parseCertificate(issuer.certificate);
        return {
          ...issuer,
          isRoot,
          parsedCertificate,
        };
      });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'Issuers', route: 'issuers.index', model: this.secretMountPath.currentPath },
    ];
  }
}
