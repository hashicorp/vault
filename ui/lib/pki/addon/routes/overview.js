/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfig } from 'pki/decorators/check-issuers';
import { hash } from 'rsvp';
import {
  PkiListCertsListEnum,
  PkiListRolesListEnum,
  PkiListIssuersListEnum,
} from '@hashicorp/vault-client-typescript';

export const PKI_DEFAULT_EMPTY_STATE_MSG =
  "This PKI mount hasn't yet been configured with a certificate issuer.";

export const getCliMessage = (msg) => {
  if (!msg) return PKI_DEFAULT_EMPTY_STATE_MSG;

  return `${PKI_DEFAULT_EMPTY_STATE_MSG} There are existing ${msg}. Use the CLI to perform any operations with them until an issuer is configured.`;
};

@withConfig()
export default class PkiOverviewRoute extends Route {
  @service secretMountPath;
  @service auth;
  @service api;

  async fetchAllCertificates() {
    try {
      const { keys } = await this.api.secrets.pkiListCerts(
        this.secretMountPath.currentPath,
        PkiListCertsListEnum.TRUE
      );
      return keys;
    } catch (e) {
      return e.response.status;
    }
  }

  async fetchAllRoles() {
    try {
      const { keys } = await this.api.secrets.pkiListRoles(
        this.secretMountPath.currentPath,
        PkiListRolesListEnum.TRUE
      );
      return keys;
    } catch (e) {
      return e.response.status;
    }
  }

  async fetchAllIssuers() {
    try {
      const { keys } = await this.api.secrets.pkiListIssuers(
        this.secretMountPath.currentPath,
        PkiListIssuersListEnum.TRUE
      );
      return keys;
    } catch (e) {
      return e.response.status;
    }
  }

  async model() {
    return hash({
      hasConfig: this.pkiMountHasConfig,
      engine: this.modelFor('application'),
      roles: this.fetchAllRoles(),
      issuers: this.fetchAllIssuers(),
      certificates: this.fetchAllCertificates(),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const roles = resolvedModel.roles;
    const certificates = resolvedModel.certificates;

    controller.notConfiguredMessage = getCliMessage();

    if (roles?.length) controller.notConfiguredMessage = getCliMessage('roles');
    if (certificates?.length) controller.notConfiguredMessage = getCliMessage('certificates');
    if (roles?.length && certificates?.length)
      controller.notConfiguredMessage = getCliMessage('roles and certificates');
  }
}
