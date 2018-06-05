import Ember from 'ember';
import { task } from 'ember-concurrency';

const { inject } = Ember;

export default Ember.Component.extend({
  store: inject.service(),

  // Public API
  //value for the external mount selector
  value: null,
  onChange: () => {},

  init() {
    this._super(...arguments);
    this.get('authMethods').perform();
  },

  authMethods: task(function*() {
    let methods = yield this.get('store').findAll('auth-method');
    if (!this.get('value')) {
      this.set('value', methods.get('firstObject.accessor'));
    }
    return methods;
  }).drop(),

  actions: {
    change(value) {
      this.get('onChange')(value);
    },
  },
});
