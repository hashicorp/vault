import { set } from '@ember/object';
import Route from '@ember/routing/route';
import DS from 'ember-data';
import { inject as service } from '@ember/service';

export default Route.extend({
  pathHelp: service('path-help'),
  model(params) {
    const { path } = params;
    return this.store.findAll('auth-method').then(modelArray => {
      const model = modelArray.findBy('id', path);
      if (!model) {
        const error = new DS.AdapterError();
        set(error, 'httpStatus', 404);
        throw error;
      }
      return this.pathHelp.getPaths(model.apiPath, path).then(paths => {
        model.set('paths', paths);
        return model;
      });
    });
  },
});
