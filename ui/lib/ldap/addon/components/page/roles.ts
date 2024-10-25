/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { getOwner } from '@ember/owner';
import errorMessage from 'vault/utils/error-message';
import { tracked } from '@glimmer/tracking';
import { sanitizePath } from 'core/utils/sanitize-path';

import type LdapRoleModel from 'vault/models/ldap/role';
import type SecretEngineModel from 'vault/models/secret-engine';
import type FlashMessageService from 'vault/services/flash-messages';
import type { Breadcrumb, EngineOwner } from 'vault/vault/app-types';
import type RouterService from '@ember/routing/router-service';
import type PaginationService from 'vault/services/pagination';

interface Args {
  roles: Array<LdapRoleModel>;
  promptConfig: boolean;
  backendModel: SecretEngineModel;
  breadcrumbs: Array<Breadcrumb>;
  pageFilter: string;
  roleName?: string;
}

export default class LdapRolesPageComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly pagination: PaginationService;

  @tracked credsToRotate: LdapRoleModel | null = null;
  @tracked roleToDelete: LdapRoleModel | null = null;

  // this component renders in the roles index route and in subdirectory lists for hierarchal roles
  linkParams = (role: LdapRoleModel) => {
    const route = role.name.endsWith('/') ? 'roles.subdirectory' : 'roles.role.details';
    // remove trailing forward slash from URL
    const sanitizedName = sanitizePath(role.name);
    // if there is a pathToRole we're in a subdirectory
    // and must concat the ancestors with the leaf to get the full name
    const roleName = role.pathToRole ? `${sanitizePath(role.pathToRole)}/${sanitizedName}` : sanitizedName;
    return [route, role.type, roleName];
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
    // refresh route, which fires off lazyPaginatedQuery to re-request and filter response
    this.router.transitionTo(this.router?.currentRoute?.name, { queryParams: { pageFilter } });
  }

  @action
  async onRotate(model: LdapRoleModel) {
    try {
      const message = `Successfully rotated credentials for ${model.name}.`;
      await model.rotateStaticPassword();
      this.flashMessages.success(message);
    } catch (error) {
      this.flashMessages.danger(`Error rotating credentials \n ${errorMessage(error)}`);
    } finally {
      this.credsToRotate = null;
    }
  }

  @action
  async onDelete(model: LdapRoleModel) {
    try {
      const message = `Successfully deleted role ${model.name}.`;
      await model.destroyRecord();
      this.pagination.clearDataset('ldap/role');
      this.router.transitionTo('vault.cluster.secrets.backend.ldap.roles');
      this.flashMessages.success(message);
    } catch (error) {
      this.flashMessages.danger(`Error deleting role \n ${errorMessage(error)}`);
    } finally {
      this.roleToDelete = null;
    }
  }
}
