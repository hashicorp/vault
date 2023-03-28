/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import PkiIssuerIndexRoute from './index';
import { verifyCertificates } from 'vault/utils/parse-pki-cert';
export default class PkiIssuerDetailsRoute extends PkiIssuerIndexRoute {
  // Details route gets issuer data from PkiIssuerIndexRoute
  async setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs.push({ label: resolvedModel.id });
    const pem = await this.fetchCertByFormat(resolvedModel.id, 'pem');
    const der = await this.fetchCertByFormat(resolvedModel.id, 'der');
    const isRoot = await this.isRoot(resolvedModel);
    controller.pem = pem;
    controller.der = der;
    controller.isRotatable = isRoot;
  }

  /**
   * @private fetches cert by format so it's available for download
   */
  fetchCertByFormat(issuerId, format) {
    const endpoint = `/v1/${this.secretMountPath.currentPath}/issuer/${issuerId}/${format}`;
    const adapter = this.store.adapterFor('application');
    try {
      return adapter.rawRequest(endpoint, 'GET', { unauthenticated: true }).then(function (response) {
        if (format === 'der') {
          return response.blob();
        }
        return response.text();
      });
    } catch (e) {
      return null;
    }
  }

  async isRoot({ certificate, keyId }) {
    const isSelfSigned = await verifyCertificates(certificate, certificate);
    return isSelfSigned && !!keyId;
  }
}
