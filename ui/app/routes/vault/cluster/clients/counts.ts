/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import timestamp from 'core/utils/timestamp';
import { getUnixTime } from 'date-fns';
import { setStartTimeQuery } from 'core/utils/client-count-utils';

import type FlagsService from 'vault/services/flags';
import type StoreService from 'vault/services/store';
import type VersionService from 'vault/services/version';
import type { ModelFrom } from 'vault/vault/route';
import type ClientsRoute from '../clients';
import type ClientsCountsController from 'vault/controllers/vault/cluster/clients/counts';

export interface ClientsCountsRouteParams {
  start_time?: string | number | undefined;
  end_time?: string | number | undefined;
  ns?: string | undefined;
  mountPath?: string | undefined;
}

export default class ClientsCountsRoute extends Route {
  @service declare readonly flags: FlagsService;
  @service declare readonly store: StoreService;
  @service declare readonly version: VersionService;

  queryParams = {
    start_time: { refreshModel: true, replace: true },
    end_time: { refreshModel: true, replace: true },
    ns: { refreshModel: true, replace: true },
    mountPath: { refreshModel: false, replace: true },
  };

  beforeModel() {
    return this.flags.fetchActivatedFlags();
  }

  /**
   * This method returns the query param timestamp if it exists. If not, it returns the activity timestamp value instead.
   */
  paramOrResponseTimestamp(
    qpMillisString: string | number | undefined,
    activityTimeStamp: string | undefined
  ) {
    let timestamp: string | undefined;
    const millis = Number(qpMillisString);
    if (!isNaN(millis)) {
      timestamp = fromUnixTime(millis).toISOString();
    }
    // fallback to activity timestamp only if there was no query param
    if (!timestamp && activityTimeStamp) {
      timestamp = activityTimeStamp;
    }
    return timestamp;
  }

  async getActivity(params: ClientsCountsRouteParams) {
    let activity, activityError;
    // if CE without start time we want to skip the activity call
    // so that the user is forced to choose a date range
    if (this.version.isEnterprise || params.start_time) {
      const query = {
        // start and end params are optional -- if not provided, will fallback to API default
        start_time: this.formatTimeQuery(params?.start_time),
        end_time: this.formatTimeQuery(params?.end_time),
      };
      try {
        activity = await this.store.queryRecord('clients/activity', query);
      } catch (error) {
        activityError = error;
      }
    }
    return {
      activity,
      activityError,
    };
  }

  // Takes the string URL param and formats it as the adapter expects it,
  // if it exists and is valid
  formatTimeQuery(param: string | number | undefined) {
    let timeParam: { timestamp: number } | undefined;
    const millis = Number(param);
    if (!isNaN(millis)) {
      timeParam = { timestamp: millis };
    }
    return timeParam;
  }

  async model(params: ClientsCountsRouteParams) {
    const { config, versionHistory } = this.modelFor('vault.cluster.clients') as ModelFrom<ClientsRoute>;
    const { activity, activityError } = await this.getActivity(params);
    return {
      activity,
      activityError,
      config,
      // activity.startTime corresponds to first month with data, but we want first month returned or requested
      startTimestamp: this.paramOrResponseTimestamp(params?.start_time, activity?.byMonth[0]?.timestamp),
      endTimestamp: this.paramOrResponseTimestamp(params?.end_time, activity?.endTime),
      versionHistory,
    };
  }

  resetController(controller: ClientsCountsController, isExiting: boolean) {
    if (isExiting) {
      controller.setProperties({
        start_time: undefined,
        end_time: undefined,
        ns: undefined,
        mountPath: undefined,
      });
    }
  }
}
