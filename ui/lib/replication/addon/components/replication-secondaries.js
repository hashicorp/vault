import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  cluster: null,
  replicationMode: null,
  secondaries: null,
  onRevoke: Function.prototype,
  // TODO: added eslint disable during upgrade come back and fix.
  /* eslint-disable ember/require-return-from-computed */
  addRoute: computed('replicationMode', function() {}),
  revokeRoute: computed('replicationMode', function() {}),

  actions: {
    onConfirmRevoke() {
      this.onRevoke(...arguments);
    },
  },
});
