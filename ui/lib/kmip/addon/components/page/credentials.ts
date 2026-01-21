/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { getOwner } from '@ember/owner';

import type RouterService from '@ember/routing/router-service';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type ApiService from 'vault/services/api';
import type { CapabilitiesMap, EngineOwner } from 'vault/app-types';
import type FlashMessageService from 'vault/services/flash-messages';
import type { KmipRevokeClientCertificateRequest } from '@hashicorp/vault-client-typescript';

interface Args {
  credentials: string[];
  capabilities: CapabilitiesMap;
  roleName: string;
  scopeName: string;
  filterValue: string;
}

export default class KmipScopesPageComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked credToRevoke: string | null = null;

  get mountPoint() {
    return (getOwner(this) as EngineOwner).mountPoint;
  }

  get paginationQueryParams() {
    return (page: number) => ({ page });
  }

  @action
  onFilterChange(pageFilter: string) {
    this.router.transitionTo({ queryParams: { pageFilter } });
  }

  @action
  async revoke() {
    try {
      const { roleName, scopeName } = this.args;
      const { currentPath } = this.secretMountPath;
      const payload = { serial_number: this.credToRevoke } as KmipRevokeClientCertificateRequest;

      await this.api.secrets.kmipRevokeClientCertificate(roleName, scopeName, currentPath, payload);
      this.flashMessages.success(`Successfully revoked credentials ${this.credToRevoke}`);
      this.credToRevoke = null;
      this.router.refresh();
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Error revoking credentials ${this.credToRevoke}: ${message}`);
    }
  }
}
