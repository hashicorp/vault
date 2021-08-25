import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class MetadataShow extends Route {
  @service store;

  beforeModel() {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    this.backend = backend;
  }

  model(params) {
    let { secret } = params;
    return this.store.queryRecord('secret-v2', {
      backend: this.backend,
      id: secret,
    });
  }

  setupController(controller, model) {
    controller.set('backend', this.backend); // for backendCrumb
    controller.set('model', model);
  }
}
