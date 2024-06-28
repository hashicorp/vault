/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import type Store from '@ember-data/store';
import type RouterService from '@ember/routing/router-service';
import type PkiIssuerModel from 'vault/models/pki/issuer';
import type PkiRoleModel from 'vault/models/pki/role';

interface Args {
  issuers: PkiIssuerModel[] | number;
  roles: PkiRoleModel[] | number;
  engine: string;
}

export default class PkiOverview extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly store: Store;

  @tracked rolesValue = '';
  @tracked certificateValue = '';
  @tracked issuerValue = '';

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
