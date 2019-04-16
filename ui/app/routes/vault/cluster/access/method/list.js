import { set } from '@ember/object';
import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import DS from 'ember-data';

export default Route.extend({
  wizard: service(),
  pathHelp: service('path-help'),

  model(params) {
    const { item_type: itemType } = params;
    const { path: methodType } = this.paramsFor('vault.cluster.access.method');
    let backend = this.modelFor('vault.cluster.access.method');
    let { apiPath } = backend;
    this.pathHelp.getPaths(apiPath, methodType).then(paths => {
      debugger; // eslint-disable-line
      backend.set('paths', paths);
    });
    return backend;
  },

  // setupController(controller) {
  //   const { section_name: section } = this.paramsFor(this.routeName);
  //   let backend = this.modelFor('vault.cluster.access.method');
  //   let { apiPath } = backend;
  //   let { path } = this.paramsFor('vault.cluster.access.method');
  //   this._super(...arguments);
  //   controller.set('section', section);
  //   this.pathHelp.getPaths(apiPath, path).then(paths => {
  //     controller.set('paths', paths);
  //   });
  // },
});
