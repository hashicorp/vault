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
import { ROUTES } from 'vault/utils/routes';

interface Args {
  issuers: PkiIssuerModel[] | number;
  roles: PkiRoleModel[] | number;
  engine: string;
}

export default class PkiOverview extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly store: Store;

  @tracked rolesValue = '';
  @tracked certificateValue = '';
  @tracked issuerValue = '';

  @action
  transitionToViewCertificates() {
    this.router.transitionTo(
      ROUTES.VAULT_CLUSTER_SECRETS_BACKEND_PKI_CERTIFICATES_CERTIFICATE_DETAILS,
      this.certificateValue
    );
  }
  @action
  transitionToIssueCertificates() {
    this.router.transitionTo(ROUTES.VAULT_CLUSTER_SECRETS_BACKEND_PKI_ROLES_ROLE_GENERATE, this.rolesValue);
  }

  @action
  transitionToIssuerDetails() {
    this.router.transitionTo(
      ROUTES.VAULT_CLUSTER_SECRETS_BACKEND_PKI_ISSUERS_ISSUER_DETAILS,
      this.issuerValue
    );
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
