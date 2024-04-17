/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Controller from '@ember/controller';
import { task, timeout } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { action } from '@ember/object';

export default class ReplicationModeBaseController extends Controller {
  @service('replication-mode') rm;
  @service router;
  @service store;

  get replicationMode() {
    return this.rm.mode;
  }

  get replicationForMode() {
    if (!this.replicationMode || !this.model) return null;
    return this.model[this.replicationMode];
  }

  @task
  @waitFor
  *waitForNewClusterToInit(replicationMode) {
    // waiting for the newly enabled cluster to init
    // this ensures we don't hit a capabilities-self error, called in the model of the mode/index route
    yield timeout(1000);
    this.router.transitionTo('vault.cluster.replication.mode', replicationMode);
  }

  @action
  onDisable() {
    this.router.transitionTo('vault.cluster.replication.index');
  }

  @action
  async onEnableSuccess(resp, replicationMode, clusterMode, doTransition = false) {
    // this is extrapolated from the replication-actions mixin "submitSuccess"
    const cluster = this.model;
    if (!cluster) {
      return;
    }
    // do something to show model is pending
    cluster.set(
      replicationMode,
      this.store.createRecord('replication-attributes', {
        mode: 'bootstrapping',
      })
    );
    if (clusterMode === 'secondary' && replicationMode === 'performance') {
      // if we're enabing a secondary, there could be mount filtering,
      // so we should unload all of the backends
      this.store.unloadAll('secret-engine');
    }
    try {
      await cluster.reload();
    } catch (e) {
      // no error handling here
    }
    cluster.rollbackAttributes();
    // we should only do the transitions if called from vault.cluster.replication.index
    if (doTransition) {
      if (replicationMode == 'dr' && clusterMode === 'secondary') {
        this.router.transitionTo('vault.cluster');
      } else if (replicationMode === 'dr') {
        this.router.transitionTo('vault.cluster.replication.mode', replicationMode);
      } else {
        this.waitForNewClusterToInit.perform(replicationMode);
      }
    }
  }
}
