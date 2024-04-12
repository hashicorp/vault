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
  @service() router;
  @service('replication-mode') rm;

  get replicationMode() {
    return this.rm.mode;
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
  onEnable(replicationMode, mode) {
    if (replicationMode == 'dr' && mode === 'secondary') {
      this.router.transitionTo('vault.cluster');
    } else if (replicationMode === 'dr') {
      this.router.transitionTo('vault.cluster.replication.mode', replicationMode);
    } else {
      this.waitForNewClusterToInit.perform(replicationMode);
    }
  }
  @action
  onDisable() {
    this.router.transitionTo('vault.cluster.replication.index');
  }
}
