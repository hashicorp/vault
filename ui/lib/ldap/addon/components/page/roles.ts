/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { getOwner } from '@ember/owner';
import { tracked } from '@glimmer/tracking';

import type { LdapRole } from 'vault/secrets/ldap';
import type FlashMessageService from 'vault/services/flash-messages';
import type { Breadcrumb, EngineOwner } from 'vault/vault/app-types';
import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

interface Args {
  roles: Array<LdapRole>;
  promptConfig: boolean;
  secretsEngine: SecretsEngineResource;
  breadcrumbs: Array<Breadcrumb>;
  pageFilter: string;
}

export default class LdapRolesPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  @tracked credsToRotate: LdapRole | null = null;
  @tracked roleToDelete: LdapRole | null = null;

  isHierarchical = (name: string) => name.endsWith('/');

  linkParams = (role: LdapRole) => {
    const route = this.isHierarchical(role.name) ? 'roles.subdirectory' : 'roles.role.details';
    return [route, role.type, role.completeRoleName];
  };

  get mountPoint(): string {
    const owner = getOwner(this) as EngineOwner;
    return owner.mountPoint;
  }

  get paginationQueryParams() {
    return (page: number) => ({ page });
  }

  @action
  onFilterChange(pageFilter: string) {
    // refresh route to re-request and filter response
    this.router.transitionTo(this.router?.currentRoute?.name, { queryParams: { pageFilter } });
  }

  @action
  async onRotate(role: LdapRole) {
    try {
      await this.api.secrets.ldapRotateStaticRole(
        role.completeRoleName,
        this.secretMountPath.currentPath,
        {}
      );
      this.flashMessages.success(`Successfully rotated credentials for ${role.completeRoleName}.`);
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Error rotating credentials \n ${message}`);
    } finally {
      this.credsToRotate = null;
    }
  }

  @action
  async onDelete(role: LdapRole) {
    try {
      const { currentPath } = this.secretMountPath;
      if (role.type === 'static') {
        await this.api.secrets.ldapDeleteStaticRole(role.completeRoleName, currentPath);
      } else {
        await this.api.secrets.ldapDeleteDynamicRole(role.completeRoleName, currentPath);
      }
      this.router.transitionTo('vault.cluster.secrets.backend.ldap.roles');
      this.flashMessages.success(`Successfully deleted role ${role.completeRoleName}.`);
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(`Error deleting role \n ${message}`);
    } finally {
      this.roleToDelete = null;
    }
  }
}
