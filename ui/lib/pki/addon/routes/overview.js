/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfig } from 'pki/decorators/check-config';
import { hash } from 'rsvp';

@withConfig()
export default class PkiOverviewRoute extends Route {
  @service secretMountPath;
  @service auth;
  @service store;

  async fetchAllCertificates() {
    try {
      return await this.store.query('pki/certificate/base', { backend: this.secretMountPath.currentPath });
    } catch (e) {
      return e.httpStatus;
    }
  }

  async fetchAllRoles() {
    try {
      return await this.store.query('pki/role', { backend: this.secretMountPath.currentPath });
    } catch (e) {
      return e.httpStatus;
    }
  }

  async fetchAllIssuers() {
    try {
      return await this.store.query('pki/issuer', { backend: this.secretMountPath.currentPath });
    } catch (e) {
      return e.httpStatus;
    }
  }

  async model() {
    return hash({
      hasConfig: this.shouldPromptConfig,
      engine: this.modelFor('application'),
      roles: this.fetchAllRoles(),
      issuers: this.fetchAllIssuers(),
      certificates: this.fetchAllCertificates(),
    });
  }
}
