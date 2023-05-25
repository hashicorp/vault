/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import Component from '@glimmer/component';
// TYPES
import Store from '@ember-data/store';
import RouterService from '@ember/routing/router-service';
import PkiIssuerModel from 'vault/models/pki/issuer';
import PkiRoleModel from 'vault/models/pki/role';

interface Args {
  issuers: PkiIssuerModel | number;
  roles: PkiRoleModel | number;
  engine: string;
}

export default class PkiOverview extends Component<Args> {
  @service declare readonly router: RouterService;
  @service declare readonly store: Store;

  @tracked rolesValue = '';
  @tracked certificateValue = '';

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
}
