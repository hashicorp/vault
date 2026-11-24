/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { SecretsApiPkiListIssuersListEnum } from '@hashicorp/vault-client-typescript';
import PkiRoleForm from 'vault/forms/secrets/pki/role';

export default class PkiRoleEditRoute extends Route {
  @service api;
  @service secretMountPath;

  async model() {
    const { role: name } = this.paramsFor('roles/role');
    const backend = this.secretMountPath.currentPath;

    const role = await this.api.secrets.pkiReadRole(name, backend).then((role) => ({ name, ...role }));

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
      form: new PkiRoleForm(role),
      issuers,
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const { form } = resolvedModel;
    const { name } = form.data;
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview', model: this.secretMountPath.currentPath },
      { label: 'Roles', route: 'roles.index', model: this.secretMountPath.currentPath },
      { label: name, route: 'roles.role.details', models: [this.secretMountPath.currentPath, name] },
      { label: 'Edit' },
    ];
  }
}
