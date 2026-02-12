/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type CapabilitiesService from 'vault/services/capabilities';

export default class KmipCredentialsShowRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly capabilities: CapabilitiesService;

  async model(params: { serial: string }) {
    const { role_name, scope_name } = this.paramsFor('credentials');
    const { serial: serial_number } = params;
    const { currentPath } = this.secretMountPath;

    const { data: credentials } = await this.api.secrets.kmipRetrieveClientCertificate(
      role_name as string,
      scope_name as string,
      currentPath,
      (context) => this.api.addQueryParams(context, { serial_number })
    );
    const capabilities = await this.capabilities.for('kmipCredentialsRevoke', {
      backend: currentPath,
      role: role_name,
      scope: scope_name,
    });

    return { credentials, capabilities, scopeName: scope_name, roleName: role_name };
  }
}
