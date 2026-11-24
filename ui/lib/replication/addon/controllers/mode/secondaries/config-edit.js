/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Controller from '@ember/controller';

/** Shares filter submission between create and edit without cached model rollback. */
export default class ConfigEditController extends Controller {
  @service api;
  @service flashMessages;
  @service('app-router') router;
  @service('replication-mode') rm;

  get replicationMode() {
    return this.rm.mode;
  }

  @action
  async saveConfig(config) {
    const isDelete = config.mode === null;
    const id = config.id;
    const redirectArgs = isDelete
      ? ['vault.cluster.replication.mode.secondaries', this.replicationMode]
      : ['vault.cluster.replication.mode.secondaries.config-show', id];

    try {
      if (isDelete) {
        if (!config.isNew) {
          await this.api.sys.systemDeleteReplicationPerformancePrimaryPathsFilterId(id);
        }
      } else {
        const { data } = config.toJSON();
        await this.api.sys.systemWriteReplicationPerformancePrimaryPathsFilterId(id, data);
      }
      await this.router.transitionTo(...redirectArgs).followRedirects();
      this.flashMessages.success(
        `The performance mount filter config for the secondary ${id} was successfully ${
          isDelete ? 'deleted' : 'saved'
        }.`
      );
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.flashMessages.danger(
        `There was an error ${isDelete ? 'deleting' : 'saving'} the config for ${id}: ${message}`
      );
    }
  }
}
