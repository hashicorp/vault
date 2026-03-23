/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import KmipRoleForm from 'vault/forms/secrets/kmip/role';

import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type { KmipWriteRoleRequest } from '@hashicorp/vault-client-typescript';

export default class KmipRoleEditRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  async model(params: { scope_name: string; role_name: string }) {
    const { currentPath } = this.secretMountPath;
    const { scope_name: scopeName, role_name: roleName } = params;

    const { data: role } = await this.api.secrets.kmipReadRole(roleName, scopeName, currentPath);

    return {
      form: new KmipRoleForm(role as KmipWriteRoleRequest),
      roleName,
      scopeName,
    };
  }
}
