/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { IdentityApiGroupListByIdListEnum } from '@hashicorp/vault-client-typescript';
import { accessIdentityGroupsListViewConfig } from 'vault/utils/constants/list-view-config/access-identity-groups';

import type ApiService from 'vault/services/api';
import type Controller from '@ember/controller';
import type CapabilitiesService from 'vault/services/capabilities';
import type { Group } from 'vault/vault/identity';

// keyInfoToArray needs a concrete type to avoid index-signature errors on .id access.
interface GroupKeyInfo {
  id: string;
  name: string;
}

interface RouteParams extends Record<string, unknown> {
  page: string;
  pageSize: string;
}

export type GroupListItem = Pick<Group, 'id' | 'name' | 'type' | 'alias'> & {
  canDelete: boolean;
  canEdit: boolean;
  canAddAlias: boolean;
};

export interface IdentityGroupsIndexModel {
  groups: GroupListItem[];
  listViewConfig: typeof accessIdentityGroupsListViewConfig;
  page: number;
  pageSize: number;
}

interface RouteController extends Controller {
  page: number;
  pageSize: number;
}

export default class IdentityGroupsIndexRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly capabilities: CapabilitiesService;

  queryParams = {
    page: { refreshModel: true },
  };

  async model(params: RouteParams) {
    const { page, pageSize } = params;

    try {
      const response = await this.api.identity.groupListById(IdentityApiGroupListByIdListEnum.TRUE);
      const items = this.api.keyInfoToArray<GroupKeyInfo>(response);

      const capabilityPaths = items.flatMap((item) => [
        this.capabilities.pathFor('identityCapabilities', { identityType: 'group', id: item.id }),
        this.capabilities.pathFor('groupAlias'),
      ]);

      const [capabilitiesMap, detailsResponses] = await Promise.all([
        this.capabilities.fetch(capabilityPaths),
        Promise.all(items.map((item) => this.api.identity.groupReadById(item.id))),
      ]);

      const detailsMap: Record<string, Partial<Pick<Group, 'type' | 'alias'>>> = {};
      detailsResponses.forEach((res) => {
        const d = (res as unknown as { data: Group })?.data;
        if (d?.id) {
          detailsMap[d.id] = { type: d.type, alias: d.alias };
        }
      });

      const groupAliasPath = this.capabilities.pathFor('groupAlias');
      const groupAliasCapabilities = capabilitiesMap[groupAliasPath];

      const groups = items.map((item) => {
        const capPath = this.capabilities.pathFor('identityCapabilities', {
          identityType: 'group',
          id: item.id,
        });
        const itemCapabilities = capabilitiesMap[capPath];

        return {
          ...item,
          ...(detailsMap[item.id] || {}),
          canDelete: itemCapabilities?.canDelete || false,
          canEdit: itemCapabilities?.canUpdate || false,
          canAddAlias: groupAliasCapabilities?.canCreate || false,
        };
      }) as GroupListItem[];

      const listViewConfig = { ...accessIdentityGroupsListViewConfig };
      listViewConfig.breadcrumbs = [
        { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
        { label: 'Groups' },
      ];

      return {
        groups,
        listViewConfig,
        page: Number(page) || 1,
        pageSize: Number(pageSize) || 10,
      };
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        const listViewConfig = { ...accessIdentityGroupsListViewConfig };
        listViewConfig.breadcrumbs = [
          { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
          { label: 'Groups' },
        ];
        return { groups: [], listViewConfig, page: 1, pageSize: Number(pageSize) || 10 };
      }
      throw error;
    }
  }

  resetController(controller: RouteController, isExiting: boolean) {
    if (isExiting) {
      controller.page = 1;
      controller.pageSize = 10;
    }
  }

  @action
  reload() {
    this.refresh();
  }
}
