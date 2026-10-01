/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import engineDisplayData from 'vault/helpers/engines-display-data';
import SecretsEngineResource from 'vault/resources/secrets/engine';
import { secretsBackendsListViewConfig } from 'vault/utils/constants/list-view-config/secrets-backends';

import type Controller from '@ember/controller';
import type ApiService from 'vault/services/api';
import type CapabilitiesService from 'vault/services/capabilities';

interface RouteParams extends Record<string, unknown> {
  page: string;
  pageSize: string;
}

interface BackendsController extends Controller {
  page: number;
  pageSize: number;
  resetSearchText(): void;
}

export default class SecretsBackendsRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly capabilities: CapabilitiesService;

  queryParams = {
    page: { refreshModel: true },
    sortBy: { refreshModel: false },
    sortOrder: { refreshModel: false },
  };

  async model(params: RouteParams) {
    const { page, pageSize } = params;

    const { secret } = await this.api.sys.internalUiListEnabledVisibleMounts();
    const engines = this.api
      .responseObjectToArray(secret, 'path')
      .map((raw) => {
        const engine = new SecretsEngineResource(raw);
        const displayData = engineDisplayData(engine.engineType);

        return {
          ...raw,
          id: engine.id,
          icon: engine.icon,
          engineType: engine.engineType,
          displayName: displayData?.displayName ?? engine.engineType,
          backendLink: engine.backendLink,
          backendConfigurationLink: engine.backendConfigurationLink,
          version: engine.running_plugin_version,
          // kept for filter and shouldIncludeInList check
          _resource: engine,
        };
      })
      .filter((e) => e._resource.shouldIncludeInList)
      .map(({ _resource: _r, ...rest }) => rest);

    const rowPath = (id: string) => this.capabilities.pathFor('secretsEngineMount', { path: id });
    const capabilitiesMap = engines.length
      ? await this.capabilities.fetch(engines.map((e) => rowPath(e.id)))
      : {};

    const enginesWithCapabilities = engines.map((engine) => ({
      ...engine,
      capabilities: {
        // cubbyhole/ is a built-in per-token mount that Vault refuses to disable.
        canDelete: engine.engineType !== 'cubbyhole' && !!capabilitiesMap[rowPath(engine.id)]?.canDelete,
      },
    }));

    const listViewConfig = { ...secretsBackendsListViewConfig };
    listViewConfig.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'Secrets engines' },
    ];

    return {
      engines: enginesWithCapabilities,
      listViewConfig,
      page: Number(page) || 1,
      pageSize: Number(pageSize) || 10,
    };
  }

  resetController(
    controller: BackendsController & { sortBy: string | undefined; sortOrder: string | undefined },
    isExiting: boolean
  ) {
    if (isExiting) {
      controller.page = 1;
      controller.pageSize = 10;
      controller.resetSearchText();
      controller.sortBy = undefined;
      controller.sortOrder = undefined;
    }
  }
}
