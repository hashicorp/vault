/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type CapabilitiesService from 'vault/services/capabilities';

export default class KmipRoleRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly capabilities: CapabilitiesService;

  async model(params: { scope_name: string; role_name: string }) {
    const { currentPath } = this.secretMountPath;
    const { scope_name, role_name } = params;

    const { data: role } = await this.api.secrets.kmipReadRole(
      params.role_name,
      params.scope_name,
      currentPath
    );
    const capabilities = await this.capabilities.for('kmipRole', {
      backend: currentPath,
      scope: scope_name,
      name: role_name,
    });

    return {
      role,
      roleName: role_name,
      scopeName: scope_name,
      capabilities,
    };
  }
}
