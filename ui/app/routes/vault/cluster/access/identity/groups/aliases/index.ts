/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { IdentityApiGroupListAliasesByIdListEnum } from '@hashicorp/vault-client-typescript';
import { accessIdentityGroupAliasesListViewConfig } from 'vault/utils/constants/list-view-config/access-identity-group-aliases';
import { attachAliasCapabilities } from 'vault/utils/identity-helpers';

import type ApiService from 'vault/services/api';
import type Controller from '@ember/controller';
import type CapabilitiesService from 'vault/services/capabilities';
import type { ListViewDisplayConfig } from 'core/types/list-view-config';

interface AliasKeyInfo {
  id: string;
  name: string;
  mount_type?: string;
  mount_accessor?: string;
  mount_path?: string;
  canonical_id?: string;
}

// RouteParams must extend Record<string, unknown> for the Ember Route base type.
// page/pageSize are optional — Ember omits query params that are not present in the URL.
interface RouteParams extends Record<string, unknown> {
  page?: string;
  pageSize?: string;
}

export type GroupAliasListItem = AliasKeyInfo & {
  canDelete: boolean;
  canEdit: boolean;
};

export interface IdentityGroupAliasesIndexModel {
  aliases: GroupAliasListItem[];
  listViewConfig: ListViewDisplayConfig;
  page: number;
  pageSize: number;
}

export default class IdentityGroupAliasesIndexRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly capabilities: CapabilitiesService;

  queryParams = {
    page: { refreshModel: true },
  };

  async model(params: RouteParams): Promise<IdentityGroupAliasesIndexModel> {
    const { page, pageSize } = params;

    const listViewConfig = { ...accessIdentityGroupAliasesListViewConfig };
    listViewConfig.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'Groups' },
    ];

    try {
      const response = await this.api.identity.groupListAliasesById(
        IdentityApiGroupListAliasesByIdListEnum.TRUE
      );
      const items = this.api.keyInfoToArray<AliasKeyInfo>(response);

      const aliases = (await attachAliasCapabilities({
        aliases: items,
        identityType: 'group',
        capabilities: this.capabilities,
      })) as GroupAliasListItem[];

      return {
        aliases,
        listViewConfig,
        page: Number(page) || 1,
        pageSize: Number(pageSize) || 10,
      };
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return { aliases: [], listViewConfig, page: 1, pageSize: Number(pageSize) || 10 };
      }
      throw error;
    }
  }

  resetController(controller: Controller & { page: number; pageSize: number }, isExiting: boolean) {
    if (isExiting) {
      controller.page = 1;
      controller.pageSize = 10;
    }
  }
}
