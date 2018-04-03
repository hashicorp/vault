import Ember from 'ember';
import ReplicationActions from 'vault/mixins/replication-actions';

const { computed } = Ember;

const DEFAULTS = {
  token: null,
  primary_api_addr: null,
  primary_cluster_addr: null,
  errors: [],
  id: null,
  replicationMode: null,
};

export default Ember.Component.extend(ReplicationActions, DEFAULTS, {
  replicationMode: null,
  selectedAction: null,
  tagName: 'form',

  didReceiveAttrs() {
    this._super(...arguments);
  },

  model: null,
  cluster: computed.alias('model'),
  loading: false,
  onSubmit: null,

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
      return this.submitHandler(...arguments);
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
