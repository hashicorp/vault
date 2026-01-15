/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { paginate } from 'core/utils/paginate-list';
import { KmipListScopesListEnum } from '@hashicorp/vault-client-typescript';

import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type Controller from '@ember/controller';
import type CapabilitiesService from 'vault/services/capabilities';

interface KmipScopesController extends Controller {
  pageFilter: string | undefined;
  page: number | undefined;
}

export default class KmipScopesRoute extends Route {
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

    try {
      const { keys } = await this.api.secrets.kmipListScopes(currentPath, KmipListScopesListEnum.TRUE);
      const scopes = keys ? paginate(keys, { page: Number(page) || 1, filter: pageFilter }) : [];
      // fetch capabilities for filtered scopes
      const paths = scopes.map((scope) =>
        this.capabilities.pathFor('kmipScope', { backend: currentPath, name: scope })
      );
      const capabilities = paths ? await this.capabilities.fetch(paths) : {};

      return { scopes, capabilities };
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return { scopes: [], capabilities: {} };
      }
      throw error;
    }
  }

  resetController(controller: KmipScopesController, isExiting: boolean) {
    if (isExiting) {
      controller.set('pageFilter', undefined);
      controller.set('page', undefined);
    }
  }
}
