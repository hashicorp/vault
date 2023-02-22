import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class diff extends Route {
  @service store;

  beforeModel() {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    this.backend = backend;
  }

  model(params) {
    const { id } = params;
    return this.store.queryRecord('secret-v2', {
      backend: this.backend,
      id,
    });
  }

  setupController(controller, model) {
    controller.set('backend', this.backend); // for backendCrumb
    controller.set('id', model.id); // for navigation on tabs
    controller.set('model', model);
  }
}
