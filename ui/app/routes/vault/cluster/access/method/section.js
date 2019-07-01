/* eslint-disable prettier/prettier */
import { set } from '@ember/object';
import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import DS from 'ember-data';

export default Route.extend({
  wizard: service(),

  model(params) {
    const { section_name: section } = params;
    if (section !== 'configuration') {
      const error = new DS.AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }
    let backend = this.modelFor('vault.cluster.access.method');
    this.wizard.transitionFeatureMachine(this.wizard.featureState, 'DETAILS', backend.type);
    return backend;
  },

  setupController(controller) {
    const { section_name: section } = this.paramsFor(this.routeName);
    this._super(...arguments);
    controller.set('section', section);
    let method = this.modelFor('vault.cluster.access.method');
    let paths = method.paths.navPaths.map(pathInfo => pathInfo.path);
    controller.set('paths', paths);
  },
});
