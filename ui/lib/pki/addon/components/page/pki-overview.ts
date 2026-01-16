/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';

import type RouterService from '@ember/routing/router-service';

interface Args {
  issuers: string[] | number;
  roles: string[] | number;
  certificates: string[] | number;
  engine: string;
}

export default class PkiOverview extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;

  @tracked rolesValue = '';
  @tracked certificateValue = '';
  @tracked issuerValue = '';

  // format issuers, roles and certificates to pass in to SearchSelect components
  get searchSelectOptions() {
    const issuers = Array.isArray(this.args.issuers) ? this.args.issuers : [];
    const roles = Array.isArray(this.args.roles) ? this.args.roles : [];
    const certificates = Array.isArray(this.args.certificates) ? this.args.certificates : [];

    return {
      issuers: issuers.map((issuer) => ({ name: issuer, id: issuer })),
      roles: roles.map((role) => ({ name: role, id: role })),
      certificates: certificates.map((certificate) => ({ name: certificate, id: certificate })),
    };
  }

  @action
  transitionToViewCertificates() {
    this.router.transitionTo(
      'vault.cluster.secrets.backend.pki.certificates.certificate.details',
      this.certificateValue
    );
  }
  @action
  transitionToIssueCertificates() {
    this.router.transitionTo('vault.cluster.secrets.backend.pki.roles.role.generate', this.rolesValue);
  }

  @action
  transitionToIssuerDetails() {
    this.router.transitionTo('vault.cluster.secrets.backend.pki.issuers.issuer.details', this.issuerValue);
  }

  @action
  handleRolesInput(roles: string) {
    if (Array.isArray(roles)) {
      this.rolesValue = roles[0];
    } else {
      this.rolesValue = roles;
    }
  }

  @action
  handleCertificateInput(certificate: string) {
    if (Array.isArray(certificate)) {
      this.certificateValue = certificate[0];
    } else {
      this.certificateValue = certificate;
    }
  }

  @action
  handleIssuerSearch(issuers: string) {
    if (Array.isArray(issuers)) {
      this.issuerValue = issuers[0];
    } else {
      this.issuerValue = issuers;
    }
  }
}
