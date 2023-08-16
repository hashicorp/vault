import { alias } from '@ember/object/computed';
import Component from '@ember/component';
import ReplicationActions from 'core/mixins/replication-actions';
import layout from '../templates/components/replication-actions';
import { A } from '@ember/array';

const DEFAULTS = {
  token: null,
  primary_api_addr: null,
  primary_cluster_addr: null,
  errors: A(),
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
