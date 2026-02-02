/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type CapabilitiesService from 'vault/services/capabilities';
import type SecretMountPath from 'vault/services/secret-mount-path';

export default class KmipCredentialsGenerateRoute extends Route {
  @service declare readonly capabilities: CapabilitiesService;
  @service declare readonly secretMountPath: SecretMountPath;

  async model() {
    const { scope_name, role_name } = this.paramsFor('credentials');
    const capabilities = await this.capabilities.for('kmipCredentialsRevoke', {
      backend: this.secretMountPath.currentPath,
      role: role_name,
      scope: scope_name,
    });
    return {
      scopeName: scope_name,
      roleName: role_name,
      capabilities,
    };
  }
}
