import Ember from 'ember';

const { computed } = Ember;

export default Ember.Component.extend({
  cluster: null,
  replicationMode: null,
  secondaries: null,
  onRevoke: Function.prototype,

  addRoute: computed('replicationMode', function() {}),
  revokeRoute: computed('replicationMode', function() {}),

  actions: {
    onConfirmRevoke() {
      this.get('onRevoke')(...arguments);
    },
  },
});
