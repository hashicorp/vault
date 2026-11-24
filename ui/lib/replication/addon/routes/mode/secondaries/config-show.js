/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Base from '../../replication-base';

/** Treats an absent filter as the unfiltered state, while surfacing other read errors. */
export default class ConfigShowRoute extends Base {
  async model(params) {
    const id = params.secondary_id;
    let config;
    try {
      const { data } = await this.api.sys.systemReadReplicationPerformancePrimaryPathsFilterId(id);
      config = { ...data, id };
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status !== 404) {
        throw error;
      }
      config = { id };
    }
    return {
      cluster: this.modelFor('application'),
      config,
    };
  }
  redirect(model) {
    const cluster = model.cluster;
    const replicationMode = this.paramsFor('mode').replication_mode;
    if (
      !this.version.hasPerfReplication ||
      replicationMode !== 'performance' ||
      !cluster[replicationMode].isPrimary
    ) {
      return this.router.transitionTo('vault.cluster.replication.mode', cluster.name, replicationMode);
    }
  }
}
