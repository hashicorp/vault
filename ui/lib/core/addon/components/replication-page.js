/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

/**
 * @module ReplicationPage
 * The `ReplicationPage` component is the parent contextual component that holds the replication-dashboard, and various replication-<name>-card components.
 * It is the top level component on routes displaying replication dashboards.
 *
 * @example
 * ```js
 * <ReplicationPage
    @model={{cluster}}
    />
 * ```
 * @param {Object} model=null - An Ember data object that is pulled from the Ember Cluster Model.
 */

const MODE = {
  dr: 'Disaster Recovery',
  performance: 'Performance',
};

export default class ReplicationPage extends Component {
  @service store;
  @service router;
  @tracked reindexingDetails = null;

  @action onModeUpdate(evt, replicationMode) {
    // Called on did-insert and did-update
    this.getReplicationModeStatus.perform(replicationMode);
  }

  @task
  @waitFor
  *getReplicationModeStatus(replicationMode) {
    let resp;
    if (this.isSummaryDashboard) {
      // the summary dashboard is not mode specific and will error
      // while running replication/null/status in the replication-mode adapter
      return;
    }

    try {
      resp = yield this.store.adapterFor('replication-mode').fetchStatus(replicationMode);
    } catch (e) {
      // do not handle error
    } finally {
      this.reindexingDetails = resp;
    }
  }
  get isSummaryDashboard() {
    const currentRoute = this.router.currentRouteName;

    // we only show the summary dashboard in the replication index route
    if (currentRoute === 'vault.cluster.replication.index') {
      const drMode = this.args.model.dr.mode;
      const performanceMode = this.args.model.performance.mode;
      return drMode === 'primary' && performanceMode === 'primary';
    }
    return '';
  }
  get formattedReplicationMode() {
    // dr or performance 🤯
    if (this.isSummaryDashboard) {
      return 'Disaster Recovery & Performance';
    }
    const mode = this.args.model.replicationMode;
    return MODE[mode];
  }
  get clusterMode() {
    // primary or secondary
    if (this.isSummaryDashboard) {
      // replicationAttrs does not exist when summaryDashboard
      return 'primary';
    }
    return this.args.model.replicationAttrs.mode;
  }
  get isLoadingData() {
    if (this.isSummaryDashboard) {
      return false;
    }
    const { clusterId, replicationDisabled } = this.args.model.replicationAttrs;
    if (this.clusterMode === 'bootstrapping' || (!clusterId && !replicationDisabled)) {
      // if clusterMode is bootstrapping
      // if no clusterId, the data hasn't loaded yet, wait for another status endpoint to be called
      return true;
    }
    return false;
  }
  get isSecondary() {
    return this.clusterMode === 'secondary';
  }
  get replicationDetailsSummary() {
    if (this.isSummaryDashboard) {
      const combinedObject = {};
      combinedObject.dr = this.args.model['dr'];
      combinedObject.performance = this.args.model['performance'];
      return combinedObject;
    }
    return {};
  }
  get replicationDetails() {
    if (this.isSummaryDashboard) {
      // Cannot return null
      return {};
    }
    const { replicationMode } = this.args.model;
    return this.args.model[replicationMode];
  }
  get isDisabled() {
    if (this.replicationDetails.mode === 'disabled' || this.replicationDetails.mode === 'primary') {
      return true;
    }
    return false;
  }
  get message() {
    let msg;
    if (this.args.model.anyReplicationEnabled) {
      msg = `This ${this.formattedReplicationMode} secondary has not been enabled.  You can do so from the ${this.formattedReplicationMode} Primary.`;
    } else {
      msg = `This cluster has not been enabled as a ${this.formattedReplicationMode} Secondary. You can do so by enabling replication and adding a secondary from the ${this.formattedReplicationMode} Primary.`;
    }
    return msg;
  }
}
