import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import Component from '@ember/component';
import decodeConfigFromJWT from 'replication/utils/decode-config-from-jwt';
import ReplicationActions from 'core/mixins/replication-actions';
import { task } from 'ember-concurrency';
import { A } from '@ember/array';

const DEFAULTS = {
  token: null,
  id: null,
  loading: false,
  errors: A(),
  primary_api_addr: null,
  primary_cluster_addr: null,
  ca_file: null,
  ca_path: null,
};

export default Component.extend(ReplicationActions, DEFAULTS, {
  replicationMode: 'dr',
  mode: 'primary',
  wizard: service(),
  version: service(),
  didReceiveAttrs() {
    this._super(...arguments);
    const initialReplicationMode = this.initialReplicationMode;
    if (initialReplicationMode) {
      this.set('replicationMode', initialReplicationMode);
    }
  },
  showModeSummary: false,
  initialReplicationMode: null,
  cluster: null,

  replicationAttrs: alias('cluster.replicationAttrs'),

  tokenIncludesAPIAddr: computed('token', function () {
    const config = decodeConfigFromJWT(this.token);
    return config && config.addr ? true : false;
  }),

  disallowEnable: computed(
    'replicationMode',
    'version.hasPerfReplication',
    'mode',
    'tokenIncludesAPIAddr',
    'primary_api_addr',
    function () {
      const inculdesAPIAddr = this.tokenIncludesAPIAddr;
      if (this.replicationMode === 'performance' && this.version.hasPerfReplication === false) {
        return true;
      }
      if (this.mode !== 'secondary' || inculdesAPIAddr || (!inculdesAPIAddr && this.primary_api_addr)) {
        return false;
      }
      return true;
    }
  ),

  reset() {
    this.setProperties(DEFAULTS);
  },

  submit: task(function* () {
    try {
      yield this.submitHandler.perform(...arguments);
    } catch (e) {
      // do not handle error
    }
  }),
  actions: {
    onSubmit(/*action, mode, data, event*/) {
      this.submit.perform(...arguments);
    },

    clear() {
      this.reset();
      this.setProperties({
        token: null,
        id: null,
      });
    },
  },
});
