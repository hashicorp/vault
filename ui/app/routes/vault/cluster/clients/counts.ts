/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import timestamp from 'core/utils/timestamp';
import { getUnixTime } from 'date-fns';

import type StoreService from 'vault/services/store';
import type { ClientsRouteModel } from '../clients';

export default class ClientCountTokenRoute extends Route {
  @service declare readonly store: StoreService;

  queryParams = {
    start_time: {
      refreshModel: true,
    },
    end_time: {
      refreshModel: true,
    },
  };

  currentTimestamp = getUnixTime(timestamp.now());

  async getActivity(start_time: number, end_time: number) {
    let activity, activityError;
    try {
      activity = await this.store.queryRecord('clients/activity', {
        start_time: { timestamp: start_time },
        end_time: { timestamp: end_time },
      });
    } catch (error) {
      activityError = error;
    }
    return [activity, activityError];
  }

  async model(params: { start_time: string; end_time: string }) {
    const { config } = this.modelFor('vault.cluster.clients') as ClientsRouteModel;
    const startTimestamp = Number(params.start_time) || getUnixTime(config.billingStartTimestamp);
    const endTimestamp = Number(params.end_time) || this.currentTimestamp;
    // only make network request if we have a start_time
    const [activity, activityError] = startTimestamp
      ? await this.getActivity(startTimestamp, endTimestamp)
      : [];
    return {
      config,
      activity,
      activityError,
      startTimestamp,
      endTimestamp,
      currentTimestamp: this.currentTimestamp,
    };
  }
}
