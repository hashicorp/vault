/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { service } from '@ember/service';

import type { CertificateRouteModel } from 'pki/routes/external/certificates/certificate';
import type RouterService from '@ember/routing/router-service';

interface Args {
  model: CertificateRouteModel;
}

export default class ExternalPkiPageCertificateComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;

  @action
  refresh() {
    this.router.refresh('vault.cluster.secrets.backend.pki.external.certificates.certificate');
  }

  get orderParams() {
    const { order_status, role_name, identifiers } = this.args.model.certLookup;
    return { details: { order_status, role_name, identifiers } };
  }

  get certParams() {
    // If the request to fetch the actual certificate failed, but the serial number lookup succeeded
    // return the cert validity dates returned by the API
    if (this.args.model.certificate?.error && !this.args.model.certificate.details) {
      const { not_before, not_after } = this.args.model.certLookup;
      return { details: { not_before, not_after }, error: this.args.model.certificate?.error };
    }
    return this.args.model.certificate;
  }
}
