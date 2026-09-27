/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { accessNamespacesListViewConfig } from 'vault/utils/constants/list-view-config/access-namespaces';
import { SystemApiSystemListNamespacesListEnum } from '@hashicorp/vault-client-typescript';

import type ApiService from 'vault/services/api';
import type Controller from '@ember/controller';
import FlagsService from 'vault/services/flags';
import NamespaceService from 'vault/services/namespace';

// RouteParams must extend Record<string, unknown> for the Ember Route base type
// page/pageSize are optional — Ember omits query params not present in the URL
interface RouteParams extends Record<string, unknown> {
  page?: string;
  pageSize?: string;
}

export interface NamespacesIndexModel {
  namespaces: Array<{ id: string }>;
  listViewConfig: typeof accessNamespacesListViewConfig;
  page: number;
  pageSize: number;
}

interface RouteController extends Controller {
  page: number;
  pageSize: number;
}

export default class NamespacesIndexRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly flags: FlagsService;
  @service declare namespace: NamespaceService;

  queryParams = {
    page: { refreshModel: true },
  };

  async model(params: RouteParams): Promise<NamespacesIndexModel> {
    const { page, pageSize } = params;

    let keys: string[] = [];
    try {
      const result = await this.api.sys.systemListNamespaces(SystemApiSystemListNamespacesListEnum.TRUE);
      keys = result.keys ?? [];
    } catch (err) {
      const { status } = await this.api.parseError(err);
      // 404 means no namespaces exist yet (CE build or unlicensed) — treat as empty list
      if (status !== 404) {
        throw err;
      }
    }

    // Strip trailing slashes to match the format used by namespace.accessibleNamespaces
    // (see services/namespace.ts → findNamespacesForUser which normalises the same way).
    const namespaces = keys.map((key) => ({ id: key.replace(/\/$/, '') }));

    const listViewConfig = { ...accessNamespacesListViewConfig };
    listViewConfig.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'Namespaces' },
    ];
    listViewConfig.badge = this.namespacePath;

    return { namespaces, listViewConfig, page: Number(page) || 1, pageSize: Number(pageSize) || 10 };
  }

  // show the full available namespace path e.g. "root/ns1/child2", "admin/ns1/child2"
  get namespacePath() {
    if (this.namespace.inRootNamespace) {
      return 'root';
    }
    if (!this.namespace.userRootNamespace && !this.flags.isHvdManaged) {
      return `root/${this.namespace.path}`;
    }
    return this.namespace.path;
  }

  resetController(controller: RouteController, isExiting: boolean) {
    if (isExiting) {
      controller.page = 1;
      controller.pageSize = 10;
    }
  }
}
