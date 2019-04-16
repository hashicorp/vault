import { set } from '@ember/object';
import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import DS from 'ember-data';

export default Route.extend({
  wizard: service(),
  pathHelp: service('path-help'),

  model(params) {
    const { section_name: section } = params;
    if (section !== 'configuration') {
      const error = new DS.AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }
    let backend = this.modelFor('vault.cluster.access.method');
    this.get('wizard').transitionFeatureMachine(
      this.get('wizard.featureState'),
      'DETAILS',
      backend.get('type')
    );
    return backend;
  },

  setupController(controller) {
    const { section_name: section } = this.paramsFor(this.routeName);
    let backend = this.modelFor('vault.cluster.access.method');
    let { apiPath } = backend;
    let { path } = this.paramsFor('vault.cluster.access.method');
    this._super(...arguments);
    controller.set('section', section);
    this.pathHelp.getPaths(apiPath, path).then(paths => {
      controller.set('paths', paths);
    });
  },
});
