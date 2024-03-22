/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { DEBUG } from '@glimmer/env';
import timestamp from 'core/utils/timestamp';
import { getUnixTime } from 'date-fns';

import type StoreService from 'vault/services/store';
import type VersionService from 'vault/services/version';
import type { ClientsRouteModel } from '../clients';
import type ClientsConfigModel from 'vault/models/clients/config';
import type ClientsVersionHistoryModel from 'vault/models/clients/version-history';
import type ClientsActivityModel from 'vault/models/clients/activity';
import type Controller from '@ember/controller';
import type AdapterError from 'ember-data/adapter'; // eslint-disable-line ember/use-ember-data-rfc-395-imports

export interface ClientsCountsRouteParams {
  start_time?: string | number | undefined;
  end_time?: string | number | undefined;
  ns?: string | undefined;
  mountPath?: string | undefined;
}

export interface ClientsCountsRouteModel {
  config: ClientsConfigModel;
  versionHistory: ClientsVersionHistoryModel;
  activity?: ClientsActivityModel;
  activityError?: AdapterError;
  isSecretsSyncActivated: boolean;
  startTimestamp: number;
  endTimestamp: number;
}
interface ClientsCountsController extends Controller {
  model: ClientsCountsRouteModel;
  start_time: number | undefined;
  end_time: number | undefined;
  ns: string | undefined;
  mountPath: string | undefined;
}

interface ActivationFlagsResponse {
  data: {
    activated: Array<string>;
    unactivated: Array<string>;
  };
}

export default class ClientsCountsRoute extends Route {
  @service declare readonly store: StoreService;
  @service declare readonly version: VersionService;

  queryParams = {
    start_time: { refreshModel: true, replace: true },
    end_time: { refreshModel: true, replace: true },
    ns: { refreshModel: false, replace: true },
    mountPath: { refreshModel: false, replace: true },
  };

  async getActivity(
    start_time: number,
    end_time: number
  ): Promise<[ClientsActivityModel | undefined, AdapterError | unknown] | [Record<string, never>, null]> {
    let activity, activityError;
    // if there is no billingStartTimestamp or selected start date initially we allow the user to manually choose a date
    // in that case bypass the query so that the user isn't stuck viewing the activity error
    if (start_time) {
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
    return [{}, null];
  }

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

  async isSecretsSyncActivated(activity: ClientsActivityModel | Record<string, never> | undefined) {
    // if there are secrets, the feature is activated
    if (activity && activity.total?.secret_syncs > 0) return true;

    // if feature is not in license, it's definitely not activated
    if (!this.version.hasSecretsSync) return false;

    // otherwise check explicitly if the feature has been activated
    const activatedFeatures = await this.getActivatedFeatures();
    return activatedFeatures.includes('secrets-sync');
  }

  async model(params: ClientsCountsRouteParams) {
    const { config, versionHistory } = this.modelFor('vault.cluster.clients') as ClientsRouteModel;
    // we could potentially make an additional request to fetch the license and get the start date from there if the config request fails
    const startTimestamp = Number(params.start_time) || getUnixTime(config.billingStartTimestamp);
    const endTimestamp = Number(params.end_time) || getUnixTime(timestamp.now());
    const [activity, activityError] = await this.getActivity(startTimestamp, endTimestamp);

    const isSecretsSyncActivated = await this.isSecretsSyncActivated(activity);

    return {
      activity,
      activityError,
      config,
      endTimestamp,
      isSecretsSyncActivated,
      startTimestamp,
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
