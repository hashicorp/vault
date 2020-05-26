import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/replication-page';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';

const MODE = {
  dr: 'Disaster Recovery',
  performance: 'Performance',
};

export default Component.extend({
  layout,
  store: service(),
  reindexingDetails: null,
  didReceiveAttrs() {
    this._super(arguments);
    this.getReplicationModeStatus.perform();
  },
  getReplicationModeStatus: task(function*() {
    let resp;
    const { replicationMode } = this.model;
    try {
      resp = yield this.get('store')
        .adapterFor('replication-mode')
        .fetchStatus(replicationMode);
    } catch (e) {
      // do not handle error
    }
    this.set('reindexingDetails', resp);
  }),
  isSummaryDashboard: computed('model.dr.{mode}', 'model.performance.{mode}', function() {
    const drMode = this.model.dr.mode;
    const performanceMode = this.model.performance.mode;
    return drMode && performanceMode === 'primary';
  }),
  formattedReplicationMode: computed('model.{replicationMode}', function() {
    // dr or performance 🤯
    const { isSummaryDashboard } = this;
    if (isSummaryDashboard) {
      return 'Summary Dashboard';
    }
    const mode = this.model.replicationMode;
    return MODE[mode];
  }),
  clusterMode: computed('model.{replicationAttrs}', 'isSummaryDashboard', function() {
    // primary or secondary
    const { model } = this;
    const { isSummaryDashboard } = this;
    if (isSummaryDashboard) {
      console.log('1. summaryDashboard exit clusterMode with primary');
      // replicationAttrs does not exist when summaryDashboard
      return 'primary';
    }
    return model.replicationAttrs.mode;
  }),
  isLoadingData: computed('clusterMode', 'model.{replicationAttrs}', function() {
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
  isSecondary: computed('clusterMode', 'isSummaryDashboard', function() {
    const { clusterMode } = this;
    const { isSummaryDashboard } = this;
    return clusterMode === 'secondary';
  }),
  replicationDetails: computed('model.{replicationMode}', 'isSummaryDashboard', function() {
    const { model } = this;
    const { isSummaryDashboard } = this;
    if (isSummaryDashboard) {
      // ARG revisit, for now just returning dr, but will need to either combine or do something else
      return model['dr'];
    }
    const replicationMode = model.replicationMode;
    return model[replicationMode];
  }),
  isDisabled: computed('replicationDetails.{mode}', function() {
    if (this.replicationDetails.mode === 'disabled' || this.replicationDetails.mode === 'primary') {
      return true;
    }
    return false;
  }),
  message: computed('model.{anyReplicationEnabled}', 'formattedReplicationMode', function() {
    if (this.model.anyReplicationEnabled) {
      return `This ${this.formattedReplicationMode} secondary has not been enabled.  You can do so from the ${this.formattedReplicationMode} Primary.`;
    }
    return `This cluster has not been enabled as a ${this.formattedReplicationMode} Secondary. You can do so by enabling replication and adding a secondary from the ${this.formattedReplicationMode} Primary.`;
  }),
});
