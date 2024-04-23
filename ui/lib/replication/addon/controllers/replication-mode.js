/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { alias } from '@ember/object/computed';
import { service } from '@ember/service';
import Controller from '@ember/controller';
import { task, timeout } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

export default Controller.extend({
  router: service(),
  rm: service('replication-mode'),
  replicationMode: alias('rm.mode'),
  waitForNewClusterToInit: task(
    waitFor(function* (replicationMode) {
      // waiting for the newly enabled cluster to init
      // this ensures we don't hit a capabilities-self error, called in the model of the mode/index route
      yield timeout(1000);
      this.router.transitionTo('vault.cluster.replication.mode', replicationMode);
    })
  ),
  actions: {
    onEnable(replicationMode, mode) {
      if (replicationMode == 'dr' && mode === 'secondary') {
        this.router.transitionTo('vault.cluster');
      } else if (replicationMode === 'dr') {
        this.router.transitionTo('vault.cluster.replication.mode', replicationMode);
      } else {
        this.waitForNewClusterToInit.perform(replicationMode);
      }
    },
    onDisable() {
      this.router.transitionTo('vault.cluster.replication.index');
    },
  },
});
