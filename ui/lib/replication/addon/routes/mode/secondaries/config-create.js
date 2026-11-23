/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import PathFilterConfigForm from 'vault/forms/path-filter-config';
import Base from '../../replication-base';

/** Uses a fresh form for absent filters. */
export default class ConfigCreateRoute extends Base {
  @service flashMessages;

  async findOrCreate(id) {
    try {
      await this.api.sys.systemReadReplicationPerformancePrimaryPathsFilterId(id);
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return new PathFilterConfigForm({ id }, { isNew: true });
      }
      throw error;
    }
    await this.router
      .transitionTo('vault.cluster.replication.mode.secondaries.config-edit', id)
      .followRedirects();
    this.flashMessages.info(
      `${id} already had a path filter config, so we loaded the config edit screen for you.`
    );
  }

  redirect(model) {
    const cluster = model.cluster;
    const replicationMode = this.replicationMode;
    if (
      !this.version.hasPerfReplication ||
      replicationMode !== 'performance' ||
      !cluster[replicationMode].isPrimary ||
      !cluster.canAddSecondary
    ) {
      return this.router.transitionTo('vault.cluster.replication.mode', replicationMode);
    }
  }

  async model(params) {
    return {
      cluster: this.modelFor('mode'),
      config: await this.findOrCreate(params.secondary_id),
    };
  }
}
