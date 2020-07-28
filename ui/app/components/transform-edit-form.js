import Component from '@ember/component';
import { inject as service } from '@ember/service';

export default Component.extend({
  init() {
    this._super(...arguments);
    // TODO: don't need confirm?
    this.set('backendType', 'transform');
  },
  store: service(),
  actions: {
    // TODO modify the parameters and potentially rename
    createOrUpdate(type, event) {
      const adapter = this.get('store').adapterFor('transform');
      adapter.createOrUpdate(store, type, snapshot, requestType);
    },
  },
});
