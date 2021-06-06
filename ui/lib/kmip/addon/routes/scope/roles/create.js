import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default Route.extend({
  store: service(),
  secretMountPath: service(),
  pathHelp: service(),
  scope() {
    return this.paramsFor('scope').scope_name;
  },
  beforeModel() {
    this.store.unloadAll('kmip/role');
    return this.pathHelp.getNewModel('kmip/role', this.secretMountPath.currentPath);
  },
  model() {
    let model = this.store.createRecord('kmip/role', {
      backend: this.secretMountPath.currentPath,
      scope: this.scope(),
    });
    return model;
  },
  setupController(controller) {
    this._super(...arguments);
    controller.set('scope', this.scope());
  },
});
