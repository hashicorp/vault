import Ember from 'ember';
import DS from 'ember-data';

import { methods } from 'vault/helpers/mountable-auth-methods';

const METHODS = methods();

export default Ember.Route.extend({
  model() {
    const { method } = this.paramsFor(this.routeName);
    return this.store.findAll('auth-method').then(() => {
      const model = this.store.peekRecord('auth-method', method);
      const modelType = model && model.get('type');
      if (!model || (modelType !== 'token' && !METHODS.findBy('type', modelType))) {
        const error = new DS.AdapterError();
        Ember.set(error, 'httpStatus', 404);
        throw error;
      }
      return model;
    });
  },
});
