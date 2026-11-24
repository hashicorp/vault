/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';

import type RouterService from '@ember/routing/router-service';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type { KmipRevokeClientCertificateRequest } from '@hashicorp/vault-client-typescript';

type KmipCredentials = {
  certificate: string;
  serial_number: string;
  ca_chain: string[];
};

interface Args {
  roleName: string;
  scopeName: string;
  credentials: KmipCredentials;
}

export default class KmipCredentialsShowPageComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @action
  async revoke() {
    const { roleName, scopeName, credentials } = this.args;
    const { serial_number } = credentials;
    const { currentPath } = this.secretMountPath;
    const payload = { serial_number } as KmipRevokeClientCertificateRequest;

    try {
      await this.api.secrets.kmipRevokeClientCertificate(roleName, scopeName, currentPath, payload);
      this.flashMessages.success('Successfully revoked credentials.');
      this.router.transitionTo('vault.cluster.secrets.backend.kmip.credentials.index', scopeName, roleName);
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Error revoking credentials ${serial_number}: ${message}`);
    }
  }
}
