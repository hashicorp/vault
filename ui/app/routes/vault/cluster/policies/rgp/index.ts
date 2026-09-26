/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { policiesRgpListViewConfig } from 'vault/utils/constants/list-view-config/policies-rgp';
import { SystemApiSystemListPoliciesRgpListEnum } from '@hashicorp/vault-client-typescript';

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

export interface PoliciesRgpIndexModel {
  policies: Array<{
    name: string;
    capabilities: Capabilities | null;
  }>;
  listViewConfig: typeof policiesRgpListViewConfig;
  page: number;
  pageSize: number;
}

export default class PoliciesRgpIndexRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly capabilities: CapabilitiesService;

  queryParams = {
    page: { refreshModel: true },
  };

  async model(params: RouteParams): Promise<PoliciesRgpIndexModel> {
    const { page, pageSize } = params;

    let keys: string[] = [];
    try {
      const result = await this.api.sys.systemListPoliciesRgp(SystemApiSystemListPoliciesRgpListEnum.TRUE);
      keys = result.keys ?? [];
    } catch (err) {
      const { status } = await this.api.parseError(err);
      // Sentinel policies can be 404 when none exist yet
      if (status !== 404) {
        throw err;
      }
    }

    // Fetch per-row capabilities for read, update, and delete on each policy path.
    const paths = keys.map((name) => this.capabilities.pathFor('policy', { policyType: 'rgp', id: name }));
    const capabilitiesMap = paths.length ? await this.capabilities.fetch(paths) : {};

    const policies = keys.map((name, i) => {
      const path = paths[i];
      return {
        name,
        capabilities: (path ? capabilitiesMap[path] : null) ?? null,
      };
    });

    const listViewConfig = { ...policiesRgpListViewConfig };
    listViewConfig.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'RGP policies' },
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
