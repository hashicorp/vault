import AdapterError from '@ember-data/adapter/error';
import { set } from '@ember/object';
import Route from '@ember/routing/route';

export default Route.extend({
  model(params) {
    const { section_name: section } = params;
    if (section !== 'configuration') {
      const error = new AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }
    return this.modelFor('vault.cluster.access.method');
  },

  setupController(controller) {
    const { section_name: section } = this.paramsFor(this.routeName);
    this._super(...arguments);
    controller.set('section', section);
    const method = this.modelFor('vault.cluster.access.method');
    controller.set(
      'paths',
      method.paths.paths.filter((path) => path.navigation)
    );
  },
});
