import { alias } from '@ember/object/computed';
import Component from '@ember/component';
import { computed } from '@ember/object';
import ReplicationActions from 'core/mixins/replication-actions';
import layout from '../templates/components/replication-actions';

const DEFAULTS = {
  token: null,
  primary_api_addr: null,
  primary_cluster_addr: null,
  errors: [],
  id: null,
  force: false,
};

export default Component.extend(ReplicationActions, DEFAULTS, {
  layout,
  replicationMode: null,
  model: null,
  cluster: alias('model'),
  // ARG this is the problem
  // Right now I'm not calling it anywhere
  reset() {
    if (!this || this.isDestroyed || this.isDestroying) {
      return;
    }
    this.setProperties(DEFAULTS);
  },

  replicationDisplayMode: computed('replicationMode', function() {
    const replicationMode = this.get('replicationMode');
    if (replicationMode === 'dr') {
      return 'DR';
    }
    if (replicationMode === 'performance') {
      return 'Performance';
    }
  }),

  actions: {
    onSubmit() {
      return this.submitHandler.perform(...arguments);
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
