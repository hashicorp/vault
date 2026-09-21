/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { policiesAclListViewConfig } from 'vault/utils/constants/list-view-config/policies-acl';
import { SystemApiPoliciesListAclPoliciesListEnum } from '@hashicorp/vault-client-typescript';

import type ApiService from 'vault/services/api';
import type CapabilitiesService from 'vault/services/capabilities';
import type Controller from '@ember/controller';
import type { Capabilities } from 'vault/app-types';

// RouteParams must extend Record<string, unknown> for the Ember Route base type
// page/pageSize are optional — Ember omits query params that are not present in the URL
interface RouteParams extends Record<string, unknown> {
  page?: string;
  pageSize?: string;
}

export interface PoliciesAclIndexModel {
  policies: Array<{
    name: string;
    capabilities: Capabilities | null;
  }>;
  listViewConfig: typeof policiesAclListViewConfig;
  page: number;
  pageSize: number;
}

export default class PoliciesAclIndexRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly capabilities: CapabilitiesService;

  queryParams = {
    page: { refreshModel: true },
  };

  async model(params: RouteParams): Promise<PoliciesAclIndexModel> {
    const { page, pageSize } = params;

    const { keys } = await this.api.sys.policiesListAclPolicies(
      SystemApiPoliciesListAclPoliciesListEnum.TRUE
    );

    const names = keys ?? [];

    // Fetch per-row capabilities for read, update, and delete on each policy path.
    const paths = names.map((name) => this.capabilities.pathFor('policy', { policyType: 'acl', id: name }));
    const capabilitiesMap = paths.length ? await this.capabilities.fetch(paths) : {};

    const policies = names.map((name, i) => {
      const path = paths[i];
      return {
        name,
        capabilities: (path ? capabilitiesMap[path] : null) ?? null,
      };
    });

    const listViewConfig = { ...policiesAclListViewConfig };
    listViewConfig.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'ACL policies' },
    ];

    return { policies, listViewConfig, page: Number(page) || 1, pageSize: Number(pageSize) || 10 };
  }

  resetController(controller: Controller & { page: number; pageSize: number }, isExiting: boolean) {
    if (isExiting) {
      controller.page = 1;
      controller.pageSize = 10;
    }
  }
}
