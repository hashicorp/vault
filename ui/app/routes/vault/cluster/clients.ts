/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ModelFrom } from 'vault/route';

import type ApiService from 'vault/services/api';
import type CapabilitiesService from 'vault/services/capabilities';
import { SystemApiVersionHistoryListEnum } from '@hashicorp/vault-client-typescript';

export type ClientsRouteModel = ModelFrom<ClientsRoute>;

export default class ClientsRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly capabilities: CapabilitiesService;

  async model() {
    const { canRead: canReadConfig, canUpdate: canUpdateConfig } =
      await this.capabilities.for('clientsConfig');
    const response = await this.api.sys
      .versionHistory(SystemApiVersionHistoryListEnum.TRUE)
      .catch(() => undefined);
    const versionHistory = response ? this.api.keyInfoToArray(response, 'version') : [];
    const config = await this.api.sys.internalClientActivityReadConfiguration().catch(() => ({}));
    return { canReadConfig, canUpdateConfig, versionHistory, config };
  }
}
