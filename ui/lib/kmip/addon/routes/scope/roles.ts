/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { paginate } from 'core/utils/paginate-list';
import { KmipListRolesListEnum } from '@hashicorp/vault-client-typescript';

import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type Controller from '@ember/controller';
import type CapabilitiesService from 'vault/services/capabilities';

interface KmipScopeRolesController extends Controller {
  pageFilter: string | undefined;
  page: number | undefined;
}

export default class KmipScopeRolesRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly capabilities: CapabilitiesService;

  queryParams = {
    page: {
      refreshModel: true,
    },
    pageFilter: {
      refreshModel: true,
    },
  };

  async model(params: { page: number; pageFilter: string }) {
    const { page, pageFilter } = params;
    const { currentPath } = this.secretMountPath;
    const { scope_name: scope } = this.paramsFor('scope');

    try {
      const { keys } = await this.api.secrets.kmipListRoles(
        scope as string,
        currentPath,
        KmipListRolesListEnum.TRUE
      );
      const roles = keys ? paginate(keys, { page: Number(page) || 1, filter: pageFilter }) : [];
      // fetch capabilities for filtered scopes
      const paths = roles.map((role) =>
        this.capabilities.pathFor('kmipRole', { backend: currentPath, scope, name: role })
      );
      const capabilities = paths ? await this.capabilities.fetch(paths) : {};

      return { roles, capabilities, scope };
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return { roles: [], capabilities: {}, scope };
      }
      throw error;
    }
  }

  resetController(controller: KmipScopeRolesController, isExiting: boolean) {
    if (isExiting) {
      controller.set('pageFilter', undefined);
      controller.set('page', undefined);
    }
  }
}
