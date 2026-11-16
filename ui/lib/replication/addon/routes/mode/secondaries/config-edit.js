/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import PathFilterConfigForm from 'vault/forms/path-filter-config';
import Base from '../../replication-base';

/** Keeps unsaved filter edits separate from subsequent reads. */
export default class ConfigEditRoute extends Base {
  async model(params) {
    const id = params.secondary_id;
    const { data } = await this.api.sys.systemReadReplicationPerformancePrimaryPathsFilterId(id);
    return {
      cluster: this.modelFor('mode.secondaries'),
      config: new PathFilterConfigForm({ ...data, id }),
    };
  }

  redirect(model) {
    const cluster = model.cluster;
    const replicationMode = this.paramsFor('mode').replication_mode;
    if (
      !this.version.hasPerfReplication ||
      replicationMode !== 'performance' ||
      !cluster[replicationMode].isPrimary ||
      !cluster.canAddSecondary
    ) {
      return this.router.transitionTo('vault.cluster.replication.mode', replicationMode);
    }
  }
}
