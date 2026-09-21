/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import engineDisplayData from 'vault/helpers/engines-display-data';
import SecretsEngineResource from 'vault/resources/secrets/engine';
import { secretsBackendsListViewConfig } from 'vault/utils/constants/list-view-config/secrets-backends';

import type Controller from '@ember/controller';
import type ApiService from 'vault/services/api';

interface RouteParams extends Record<string, unknown> {
  page: string;
  pageSize: string;
}

interface BackendsController extends Controller {
  page: number;
  pageSize: number;
}

export default class SecretsBackendsRoute extends Route {
  @service declare readonly api: ApiService;

  queryParams = {
    page: { refreshModel: true },
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

    const listViewConfig = { ...secretsBackendsListViewConfig };
    listViewConfig.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'Secrets engines' },
    ];

    return { engines, listViewConfig, page: Number(page) || 1, pageSize: Number(pageSize) || 10 };
  }

  resetController(controller: BackendsController, isExiting: boolean) {
    if (isExiting) {
      controller.page = 1;
      controller.pageSize = 10;
    }
  }
}
