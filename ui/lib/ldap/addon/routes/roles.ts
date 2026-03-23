/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { paginate, PaginateOptions } from 'core/utils/paginate-list';
import sortObjects from 'vault/utils/sort-objects';
import { fetchRoleCapabilities } from 'ldap/utils/capabilities-helper';
import {
  SecretsApiLdapListStaticRolesListEnum,
  SecretsApiLdapListDynamicRolesListEnum,
  SecretsApiLdapListStaticRolePathListEnum,
  SecretsApiLdapListRolePathListEnum,
} from '@hashicorp/vault-client-typescript';

import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type FlashMessageService from 'ember-cli-flash/services/flash-messages';
import type CapabilitiesService from 'vault/services/capabilities';
import type { LdapRole } from 'vault/secrets/ldap';

// Base class for roles/index and roles/subdirectory routes
export default class LdapRolesRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly capabilities: CapabilitiesService;

  async fetchRolesAndCapabilities(
    paginationOptions: PaginateOptions,
    roleAncestry?: { path_to_role: string; type: string }
  ) {
    // fetch both static and dynamic roles and combine the results
    const { currentPath } = this.secretMountPath;
    const { path_to_role: path, type: subType } = roleAncestry || {};
    let requests = [];
    if (path) {
      requests =
        subType === 'static'
          ? [
              this.api.secrets.ldapListStaticRolePath(
                currentPath,
                path,
                SecretsApiLdapListStaticRolePathListEnum.TRUE
              ),
            ]
          : [this.api.secrets.ldapListRolePath(currentPath, path, SecretsApiLdapListRolePathListEnum.TRUE)];
    } else {
      requests = [
        this.api.secrets.ldapListStaticRoles(currentPath, SecretsApiLdapListStaticRolesListEnum.TRUE),
        this.api.secrets.ldapListDynamicRoles(currentPath, SecretsApiLdapListDynamicRolesListEnum.TRUE),
      ];
    }
    const results = await Promise.allSettled(requests);
    const errors: string[] = [];
    const roles: LdapRole[] = [];
    for (const result of results) {
      if (result.status === 'fulfilled') {
        const type = subType ? subType : results.indexOf(result) === 0 ? 'static' : 'dynamic';
        if (result.value.keys) {
          roles.push(
            ...result.value.keys.map((name) => ({ name, type, completeRoleName: `${path || ''}${name}` }))
          );
        }
      } else if (result.status === 'rejected' && result.reason.response.status !== 404) {
        const { path, message } = await this.api.parseError(result.reason);
        errors.push(`${path}: ${message}`);
      }
    }

    if (errors.length) {
      if (errors.length === 2) {
        // throw error as normal if both requests fail
        // ignore status code and concat errors to be displayed in Page::Error component with generic message
        throw { message: 'Error fetching roles:', errors };
      } else if (!roleAncestry) {
        // if only one request fails when listing all roles, surface the error to the user as an info level flash message
        // this may help for permissions errors where a users policy may be incorrect
        this.flashMessages.info(`Error fetching roles from ${errors.join(', ')}`);
      }
    }

    const paginatedRoles = paginate<LdapRole>(sortObjects(roles, 'name'), {
      ...paginationOptions,
      filterKey: 'name',
    });
    // fetch capabilities for the roles being displayed on the current page
    const capabilities = await fetchRoleCapabilities(this.capabilities, currentPath, paginatedRoles);

    return { roles: paginatedRoles, capabilities };
  }
}
