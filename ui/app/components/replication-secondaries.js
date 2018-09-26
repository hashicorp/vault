import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
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
