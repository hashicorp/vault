/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { verifyCertificates } from 'vault/utils/parse-pki-cert';
import { hash } from 'rsvp';

export default class PkiIssuerDetailsRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const issuer = this.modelFor('issuers.issuer');
    return hash({
      issuer,
      pem: this.fetchCertByFormat(issuer.id, 'pem'),
      der: this.fetchCertByFormat(issuer.id, 'der'),
      isRotatable: this.isRoot(issuer),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'issuers', route: 'issuers.index' },
      { label: resolvedModel.issuer.id },
    ];
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
