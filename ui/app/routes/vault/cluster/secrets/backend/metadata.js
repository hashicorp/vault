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
    this.secret = secret; // set for access below in setupController
    return this.store
      .queryRecord('secret-v2', {
        backend: this.backend,
        id: secret,
      })
      .catch(() => {
        // there was an error likely in read metadata.
        // still load the page and handle what you show by filtering for this property
        this.noReadAccess = true;
      });
  }

  setupController(controller, model) {
    controller.set('backend', this.backend); // for backendCrumb
    controller.set('model', model);
    controller.set('noReadAccess', this.noReadAccess);
  }
}
