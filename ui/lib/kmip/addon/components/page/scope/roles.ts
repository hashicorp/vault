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
import FlashMessageService from 'vault/services/flash-messages';

interface Args {
  scope: string;
  roles: string[];
  capabilities: CapabilitiesMap;
  filterValue: string | undefined;
}

export default class KmipScopeRolesPageComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked roleToDelete: string | null = null;

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
  async deleteRole() {
    try {
      await this.api.secrets.kmipDeleteRole(
        this.roleToDelete as string,
        this.args.scope,
        this.secretMountPath.currentPath
      );
      this.flashMessages.success(`Successfully deleted role ${this.roleToDelete}`);
      this.roleToDelete = null;
      this.router.refresh();
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Error deleting role ${this.roleToDelete}: ${message}`);
    }
  }
}
