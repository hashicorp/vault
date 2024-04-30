/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import timestamp from 'core/utils/timestamp';
import { getUnixTime } from 'date-fns';

import type FlagsService from 'vault/services/flags';
import type StoreService from 'vault/services/store';
import type VersionService from 'vault/services/version';
import type { ModelFrom } from 'vault/vault/route';
import type ClientsRoute from '../clients';
import type ClientsActivityModel from 'vault/models/clients/activity';
import type ClientsConfigModel from 'vault/models/clients/config';
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
    ns: { refreshModel: false, replace: true },
    mountPath: { refreshModel: false, replace: true },
  };

  beforeModel() {
    return this.flags.fetchActivatedFlags();
  }

  async getActivity(start_time: number | null, end_time: number) {
    let activity, activityError;
    // if there is no start_time we want the user to manually choose a date
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
    }
    return { activity, activityError };
  }

  showSecretsSync(activity: ClientsActivityModel | undefined): boolean {
    // if there is any sync client data, show it
    if (activity && activity?.total?.secret_syncs > 0) return true;

    // otherwise, show the tab based on the cluster type and license
    if (this.version.isCommunity) return false;

    const isManaged = this.flags.isManaged;
    const onLicense = this.version.hasSecretsSync;

    // we can't tell if HVD clusters have the feature or not, so we show it by default
    // if the cluster is not HVD, show the tab if the feature is on the license
    return isManaged || onLicense;
  }

  async model(params: ClientsCountsRouteParams) {
    const { config, versionHistory } = this.modelFor('vault.cluster.clients') as ModelFrom<ClientsRoute>;
    // only enterprise versions will have a relevant billing start date, if null users must select initial start time
    let startTime = null;
    if (this.version.isEnterprise && this._hasConfig(config)) {
      startTime = getUnixTime(config.billingStartTimestamp);
    }

    const startTimestamp = Number(params.start_time) || startTime;
    const endTimestamp = Number(params.end_time) || getUnixTime(timestamp.now());
    const { activity, activityError } = await this.getActivity(startTimestamp, endTimestamp);

    const showSecretsSync = this.showSecretsSync(activity);

    return {
      activity,
      activityError,
      config,
      endTimestamp,
      showSecretsSync,
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

  _hasConfig(model: ClientsConfigModel | object): model is ClientsConfigModel {
    return 'billingStartTimestamp' in model;
  }
}
