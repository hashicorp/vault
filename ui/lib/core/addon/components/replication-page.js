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
      console.log(e);
    }
    this.set('reindexingDetails', resp);
  }),

  replicationMode: computed('model.{replicationMode}', function() {
    // dr or performance 🤯
    let mode = this.model.replicationMode;
    return MODE[mode];
  }),
  clusterMode: computed('model.{replicationAttrs}', function() {
    // primary or secondary
    const { model } = this;
    return model.replicationAttrs.mode;
  }),
  isSecondary: computed('clusterMode', function() {
    const { clusterMode } = this;
    return clusterMode === 'secondary';
  }),
  replicationDetails: computed('model.{replicationMode}', function() {
    const { model } = this;
    const replicationMode = model.replicationMode;
    return model[replicationMode];
  }),
  isDisabled: computed('replicationDetails.{mode}', function() {
    if (this.replicationDetails.mode === 'disabled' || this.replicationDetails.mode === 'primary') {
      return true;
    }
    return false;
  }),
  message: computed('model.{anyReplicationEnabled}', 'replicationMode', function() {
    if (this.model.anyReplicationEnabled) {
      return `This ${this.replicationMode} secondary has not been enabled.  You can do so from the ${this.replicationMode} Primary.`;
    }
    return `This cluster has not been enabled as a ${this.replicationMode} Secondary. You can do so by enabling replication and adding a secondary from the ${this.replicationMode} Primary.`;
  }),
});
