/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { SecretsApiPkiListIssuersListEnum } from '@hashicorp/vault-client-typescript';
import PkiRoleForm from 'vault/forms/secrets/pki/role';

export default class PkiRolesCreateRoute extends Route {
  @service api;
  @service secretMountPath;

  async model() {
    const backend = this.secretMountPath.currentPath;
    let issuers = [];
    try {
      const response = await this.api.secrets.pkiListIssuers(backend, SecretsApiPkiListIssuersListEnum.TRUE);
      issuers = this.api.keyInfoToArray(response, 'issuer_id');
    } catch (error) {
      if (error.response.status !== 404) {
        throw error;
      }
    }
    return {
      form: new PkiRoleForm({}, { isNew: true }),
      issuers,
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'Roles', route: 'roles.index', model: this.secretMountPath.currentPath },
      { label: 'Create' },
    ];
  }
}
