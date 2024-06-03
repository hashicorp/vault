/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { or } from '@ember/object/computed';
import { isPresent } from '@ember/utils';
import Mixin from '@ember/object/mixin';
import { task } from 'ember-concurrency';

export default Mixin.create({
  store: service(),
  router: service(),
  loading: or('save.isRunning', 'submitSuccess.isRunning'),
  onDisable() {},
  onPromote() {},
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
    let resp;
    try {
      resp = yield this.store
        .adapterFor('cluster')
        .replicationAction(action, replicationMode, clusterMode, data);
    } catch (e) {
      return this.submitError(e);
    }
    return yield this.submitSuccess.perform(resp, action, clusterMode);
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
