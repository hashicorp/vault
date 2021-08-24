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
    return this.store
      .queryRecord('secret-v2', {
        backend: this.backend,
        id: secret,
      })
      .then(record => {
        // ARG TODO remove when API is done
        record.customMetadata = {
          // foo: 'abc',
          // bar: '123',
          // baz: '5c07d823-3810-48f6-a147-4c06b5219e84',
        };
        return record;
      });
  }

  setupController(controller, model) {
    controller.set('backend', this.backend); // for backendCrumb
    controller.set('model', model);
  }
}
