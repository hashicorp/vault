/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import {
  SystemListSyncDestinationsListEnum,
  SystemListSyncAssociationsListEnum,
} from '@hashicorp/vault-client-typescript';
import { listDestinationsTransform } from 'sync/utils/api-transforms';

import type FlagsService from 'vault/services/flags';
import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type VersionService from 'vault/services/version';
import type CapabilitiesService from 'vault/services/capabilities';
import type { Capabilities } from 'vault/app-types';
import type {
  SystemListSyncDestinationsResponse,
  SystemListSyncAssociationsResponse,
} from '@hashicorp/vault-client-typescript';

export default class SyncSecretsOverviewRoute extends Route {
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly api: ApiService;
  @service declare readonly flags: FlagsService;
  @service declare readonly version: VersionService;
  @service declare readonly capabilities: CapabilitiesService;

  async model() {
    const isActivated = this.flags.secretsSyncIsActivated;
    const capabilitiesReq = this.capabilities.for('syncActivate');
    const requests = isActivated
      ? [
          capabilitiesReq,
          this.api.sys.systemListSyncAssociations(SystemListSyncAssociationsListEnum.TRUE).catch(() => []),
          this.api.sys.systemListSyncDestinations(SystemListSyncDestinationsListEnum.TRUE).catch(() => []),
        ]
      : [capabilitiesReq, [], []];

    const [{ canCreate, canUpdate }, { totalSecrets }, destinations] = (await Promise.all(requests)) as [
      Capabilities,
      SystemListSyncAssociationsResponse,
      SystemListSyncDestinationsResponse,
    ];

    return {
      canActivateSecretsSync: canCreate || canUpdate,
      totalSecrets,
      destinations: listDestinationsTransform(destinations),
    };
  }

  redirect() {
    if (!this.version.hasSecretsSync) {
      this.router.replaceWith('vault.cluster.dashboard');
    }
  }
}
