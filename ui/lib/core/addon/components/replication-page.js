/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-page';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';

/**
 * @module ReplicationPage
 * The `ReplicationPage` component is the parent contextual component that holds the replication-dashboard, and various replication-<name>-card components.
 * It is the top level component on routes displaying replication dashboards.
 *
 * @example
 * <ReplicationPage
    @model={{cluster}}
    />
 *
 * @param {Object} cluster=null - An Ember data object that is pulled from the Ember Cluster Model.
 */

const MODE = {
  dr: 'Disaster Recovery',
  performance: 'Performance',
};

export default Component.extend({
  layout,
  store: service(),
  router: service(),
  reindexingDetails: null,
  didReceiveAttrs() {
    this._super(arguments);
    this.getReplicationModeStatus.perform();
  },
  getReplicationModeStatus: task(function* () {
    let resp;
    const { replicationMode } = this.model;

    if (this.isSummaryDashboard) {
      // the summary dashboard is not mode specific and will error
      // while running replication/null/status in the replication-mode adapter
      return;
    }

    try {
      resp = yield this.store.adapterFor('replication-mode').fetchStatus(replicationMode);
    } catch (e) {
      // do not handle error
    }
    this.set('reindexingDetails', resp);
  }),
  isSummaryDashboard: computed('model.{performance.mode,dr.mode}', function () {
    const router = this.router;
    const currentRoute = router.get('currentRouteName');

    // we only show the summary dashboard in the replication index route
    if (currentRoute === 'vault.cluster.replication.index') {
      const drMode = this.model.dr.mode;
      const performanceMode = this.model.performance.mode;
      return drMode === 'primary' && performanceMode === 'primary';
    }
    return '';
  }),
  formattedReplicationMode: computed('model.replicationMode', 'isSummaryDashboard', function () {
    // dr or performance 🤯
    const { isSummaryDashboard } = this;
    if (isSummaryDashboard) {
      return 'Disaster Recovery & Performance';
    }
    const mode = this.model.replicationMode;
    return MODE[mode];
  }),
  clusterMode: computed('model.replicationAttrs', 'isSummaryDashboard', function () {
    // primary or secondary
    const { model } = this;
    const { isSummaryDashboard } = this;
    if (isSummaryDashboard) {
      // replicationAttrs does not exist when summaryDashboard
      return 'primary';
    }
    return model.replicationAttrs.mode;
  }),
  isLoadingData: computed('clusterMode', 'model.replicationAttrs', function () {
    const { clusterMode } = this;
    const { model } = this;
    const { isSummaryDashboard } = this;
    if (isSummaryDashboard) {
      return false;
    }
    const clusterId = model.replicationAttrs.clusterId;
    const replicationDisabled = model.replicationAttrs.replicationDisabled;
    if (clusterMode === 'bootstrapping' || (!clusterId && !replicationDisabled)) {
      // if clusterMode is bootstrapping
      // if no clusterId, the data hasn't loaded yet, wait for another status endpoint to be called
      return true;
    }
    return false;
  }),
  isSecondary: computed('clusterMode', function () {
    const { clusterMode } = this;
    return clusterMode === 'secondary';
  }),
  replicationDetailsSummary: computed('isSummaryDashboard', function () {
    const { model } = this;
    const { isSummaryDashboard } = this;
    if (!isSummaryDashboard) {
      return;
    }
    if (isSummaryDashboard) {
      const combinedObject = {};
      combinedObject.dr = model['dr'];
      combinedObject.performance = model['performance'];
      return combinedObject;
    }
    return {};
  }),
  replicationDetails: computed('model.replicationMode', 'isSummaryDashboard', function () {
    const { model } = this;
    const { isSummaryDashboard } = this;
    if (isSummaryDashboard) {
      // Cannot return null
      return {};
    }
    const replicationMode = model.replicationMode;
    return model[replicationMode];
  }),
  isDisabled: computed('replicationDetails.mode', function () {
    if (this.replicationDetails.mode === 'disabled' || this.replicationDetails.mode === 'primary') {
      return true;
    }
    return false;
  }),
  message: computed('model.anyReplicationEnabled', 'formattedReplicationMode', function () {
    let msg;
    if (this.model.anyReplicationEnabled) {
      msg = `This ${this.formattedReplicationMode} secondary has not been enabled.  You can do so from the ${this.formattedReplicationMode} Primary.`;
    } else {
      msg = `This cluster has not been enabled as a ${this.formattedReplicationMode} Secondary. You can do so by enabling replication and adding a secondary from the ${this.formattedReplicationMode} Primary.`;
    }
    return msg;
  }),
});
