import Ember from 'ember';
import DS from 'ember-data';

export default Ember.Route.extend({
  model(params) {
    const { path } = params;
    return this.store.findAll('auth-method').then(modelArray => {
      const model = modelArray.findBy('id', path);
      if (!model) {
        const error = new DS.AdapterError();
        Ember.set(error, 'httpStatus', 404);
        throw error;
      }
      return model;
    });
  },
});
