import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { task } from 'ember-concurrency';

export default Component.extend({
  store: service(),

  // Public API
  //value for the external mount selector
  value: null,
  onChange: () => {},

  init() {
    this._super(...arguments);
    this.authMethods.perform();
  },

  authMethods: task(function*() {
    let methods = yield this.store.findAll('auth-method');
    if (!this.value) {
      this.set('value', methods.get('firstObject.accessor'));
    }
    return methods;
  }).drop(),

  actions: {
    change(value) {
      this.onChange(value);
    },
  },
});
