/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { DEBUG } from '@glimmer/env';

import type { ClientsCountsRouteModel } from '../counts';
import type StoreService from 'vault/services/store';
import ClientsActivityModel from 'vault/vault/models/clients/activity';

interface ActivationFlagsResponse {
  data: {
    activated: Array<string>;
    unactivated: Array<string>;
  };
}

export default class ClientsCountsSyncRoute extends Route {
  @service declare readonly store: StoreService;

  async getActivatedFeatures() {
    try {
      const resp: ActivationFlagsResponse = await this.store
        .adapterFor('application')
        .ajax('/v1/sys/activation-flags', 'GET', { unauthenticated: true, namespace: null });
      return resp.data?.activated;
    } catch (error) {
      if (DEBUG) console.error(error); // eslint-disable-line no-console
      return [];
    }
  }

  async secretsSyncActivated(activity: ClientsActivityModel | undefined) {
    // if there are secrets, the feature is activated
    if (activity && activity.total?.secret_syncs > 0) return true;

    // otherwise check explicitly if the feature has been activated
    const activatedFeatures = await this.getActivatedFeatures();
    return activatedFeatures.includes('secret-sync');
  }

  async model() {
    const { activity, versionHistory, startTimestamp, endTimestamp } = this.modelFor(
      'vault.cluster.clients.counts'
    ) as ClientsCountsRouteModel;

    const secretsSyncActivated = await this.secretsSyncActivated(activity);

    return {
      activity,
      secretsSyncActivated,
      versionHistory,
      startTimestamp,
      endTimestamp,
    };
  }
}
