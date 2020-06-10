import { alias } from '@ember/object/computed';
import Component from '@ember/component';
import ReplicationActions from 'core/mixins/replication-actions';
import layout from '../templates/components/replication-actions';
import keys from 'vault/lib/keycodes';

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
  reset() {
    if (!this || this.isDestroyed || this.isDestroying) {
      return;
    }
    this.setProperties(DEFAULTS);
  },

  actions: {
    onSubmit() {
      // do not submit if a key is pressed and it is not the enter key
      if (event.keyCode && event.keyCode !== keys.ENTER) {
        return;
      }
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
