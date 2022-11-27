import AdapterError from '@ember-data/adapter/error';
import { set } from '@ember/object';
import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';

export default Route.extend({
  wizard: service(),

  model(params) {
    const { section_name: section } = params;
    if (section !== 'configuration') {
      const error = new AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }
    const backend = this.modelFor('vault.cluster.access.method');
    this.wizard.transitionFeatureMachine(this.wizard.featureState, 'DETAILS', backend.type);
    return backend;
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
