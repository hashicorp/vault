import { set } from '@ember/object';
import Route from '@ember/routing/route';
import DS from 'ember-data';

export default Route.extend({
  model(params) {
    const { path } = params;
    return this.store.findAll('auth-method').then(modelArray => {
      const model = modelArray.findBy('id', path);
      if (!model) {
        const error = new DS.AdapterError();
        set(error, 'httpStatus', 404);
        throw error;
      }
      return model;
    });
  },
});
