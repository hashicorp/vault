import Ember from 'ember';
import { task } from 'ember-concurrency';
import { underscore } from 'vault/helpers/underscore';

const { inject } = Ember;

export default Ember.Component.extend({
  store: inject.service(),
  flashMessages: inject.service(),
  routing: inject.service('-routing'),

  // Public API - either 'entity' or 'group'
  // this will determine which adapter is used to make the lookup call
  type: 'entity',

  param: 'alias name',
  paramValue: null,
  aliasMountAccessor: null,

  authMethods: null,

  init() {
    this._super(...arguments);
    this.get('store').findAll('auth-method').then(methods => {
      this.set('authMethods', methods);
      this.set('aliasMountAccessor', methods.get('firstObject.accessor'));
    });
  },

  adapter() {
    let type = this.get('type');
    let store = this.get('store');
    return store.adapterFor(`identity/${type}`);
  },

  data() {
    let { param, paramValue, aliasMountAccessor } = this.getProperties(
      'param',
      'paramValue',
      'aliasMountAccessor'
    );
    let data = {};

    data[underscore([param])] = paramValue;
    if (param === 'alias name') {
      data.alias_mount_accessor = aliasMountAccessor;
    }
    return data;
  },

  lookup: task(function*() {
    let flash = this.get('flashMessages');
    let type = this.get('type');
    let store = this.get('store');
    let { param, paramValue } = this.getProperties('param', 'paramValue');
    let response;
    try {
      response = yield this.adapter().lookup(store, this.data());
    } catch (err) {
      flash.danger(
        `We encountered an error attempting the ${type} lookup: ${err.message || err.errors.join('')}.`
      );
      return;
    }
    if (response) {
      return this.get('routing.router').transitionTo(
        'vault.cluster.access.identity.show',
        response.id,
        'details'
      );
    } else {
      flash.danger(`We were unable to find an identity ${type} with a "${param}" of "${paramValue}".`);
    }
  }),
});
