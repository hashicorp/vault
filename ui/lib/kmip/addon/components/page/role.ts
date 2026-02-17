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
import type { Capabilities } from 'vault/app-types';
import type FlashMessageService from 'vault/services/flash-messages';
import type { KmipWriteRoleRequest } from '@hashicorp/vault-client-typescript';

interface Args {
  role: KmipWriteRoleRequest;
  roleName: string;
  scopeName: string;
  capabilities: Capabilities;
}

export default class KmipScopesPageComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  iconState = (field: keyof KmipWriteRoleRequest) => {
    const { operation_all, operation_none } = this.args.role;
    const isEnabled = operation_all || (!operation_none && this.args.role[field]);

    return {
      class: isEnabled ? 'hds-foreground-success' : 'hds-foreground-faint',
      name: isEnabled ? 'check-circle' : 'x-square',
      label: isEnabled ? 'Enabled' : 'Disabled',
    };
  };

  @action
  async deleteRole() {
    const { roleName, scopeName } = this.args;
    const { currentPath } = this.secretMountPath;
    try {
      await this.api.secrets.kmipDeleteRole(roleName, scopeName, currentPath);
      this.flashMessages.success(`Successfully deleted role ${roleName}`);
      this.router.transitionTo('vault.cluster.secrets.backend.kmip.scope.roles', scopeName);
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Error deleting role ${roleName}: ${message}`);
    }
  }
}
