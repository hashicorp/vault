/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { or } from '@ember/object/computed';
import { isPresent } from '@ember/utils';
import Mixin from '@ember/object/mixin';
import { task } from 'ember-concurrency';

export default Mixin.create({
  api: service(),

  loading: or('save.isRunning', 'submitSuccess.isRunning'),

  onDisable() {},
  onPromote() {},

  replicationAction(action, replicationMode, clusterMode, data = {}) {
    switch (action) {
      case 'disable':
        if (replicationMode === 'dr' && clusterMode === 'primary') {
          return this.api.sys.systemWriteReplicationDrPrimaryDisable();
        }
        if (replicationMode === 'dr' && clusterMode === 'secondary') {
          return this.api.sys.systemWriteReplicationDrSecondaryDisable(data);
        }
        if (replicationMode === 'performance' && clusterMode === 'primary') {
          return this.api.sys.systemWriteReplicationPerformancePrimaryDisable();
        }
        if (replicationMode === 'performance' && clusterMode === 'secondary') {
          return this.api.sys.systemWriteReplicationPerformanceSecondaryDisable();
        }
        break;
      case 'demote':
        if (replicationMode === 'dr' && clusterMode === 'primary') {
          return this.api.sys.systemWriteReplicationDrPrimaryDemote();
        }
        if (replicationMode === 'performance' && clusterMode === 'primary') {
          return this.api.sys.systemWriteReplicationPerformancePrimaryDemote();
        }
        break;
      case 'promote':
        if (replicationMode === 'dr' && clusterMode === 'secondary') {
          return this.api.sys.systemWriteReplicationDrSecondaryPromote(data);
        }
        if (replicationMode === 'performance' && clusterMode === 'secondary') {
          return this.api.sys.systemWriteReplicationPerformanceSecondaryPromote(data);
        }
        break;
      case 'update-primary':
        if (replicationMode === 'dr' && clusterMode === 'secondary') {
          return this.api.sys.systemWriteReplicationDrSecondaryUpdatePrimary(data);
        }
        if (replicationMode === 'performance' && clusterMode === 'secondary') {
          return this.api.sys.systemWriteReplicationPerformanceSecondaryUpdatePrimary(data);
        }
        break;
      case 'recover':
        return this.api.sys.systemWriteReplicationRecover();
      case 'reindex':
        return this.api.sys.systemWriteReplicationReindex(data);
    }

    throw new Error(`Unsupported replication action: ${replicationMode}/${clusterMode}/${action}`);
  },

  submitHandler: task(function* (action, clusterMode, data, event) {
    const replicationMode = (data && data.replicationMode) || this.replicationMode;
    if (event && event.preventDefault) {
      event.preventDefault();
    }
    this.setProperties({
      errors: [],
    });
    if (data) {
      data = Object.keys(data).reduce((newData, key) => {
        var val = data[key];
        if (isPresent(val)) {
          if (key === 'dr_operation_token_primary' || key === 'dr_operation_token_promote') {
            newData['dr_operation_token'] = val;
          } else {
            newData[key] = val;
          }
        }
        return newData;
      }, {});
      delete data.replicationMode;
    }
    return yield this.save.perform(action, replicationMode, clusterMode, data);
  }),

  save: task(function* (action, replicationMode, clusterMode, data) {
    try {
      const response = yield this.replicationAction(action, replicationMode, clusterMode, data);
      return yield this.submitSuccess.perform(response, action, clusterMode);
    } catch (e) {
      const { response } = yield this.api.parseError(e);
      return this.submitError(response);
    }
  }).drop(),

  submitSuccess: task(function* (resp, action) {
    // enable action is handled separately in EnableReplicationForm component
    const cluster = this.cluster;
    if (!cluster) {
      return;
    }

    if (resp && resp.wrap_info) {
      this.set('token', resp.wrap_info.token);
    }
    if (action === 'secondary-token') {
      this.setProperties({
        loading: false,
        primary_api_addr: null,
        primary_cluster_addr: null,
      });
      return cluster;
    }
    if (this.reset) {
      this.reset();
    }
    try {
      yield cluster.reload();
    } catch (e) {
      // no error handling here
    }
    cluster.rollbackAttributes();
    if (action === 'disable') {
      yield this.onDisable();
    }
    if (action === 'promote') {
      yield this.onPromote();
    }
  }).drop(),

  submitError(e) {
    if (e.errors) {
      this.set('errors', e.errors);
    } else {
      throw e;
    }
  },
});
