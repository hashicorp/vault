import Component from '@ember/component';
import { computed } from '@ember/object';

import layout from '../templates/components/replication-secondaries';
export default Component.extend({
  layout,
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
