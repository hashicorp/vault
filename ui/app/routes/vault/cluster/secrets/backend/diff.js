import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
export default class diff extends Route {
  @service store;
  noReadAccess = false; // ARG TODO keep or remove but put in controller
  secretMetadata;

  beforeModel() {
    let { backend } = this.paramsFor('vault.cluster.secrets.backend');
    this.backend = backend; // coming in undefined on totally
  }

  model(params) {
    let { id } = params;
    return this.store
      .queryRecord('secret-v2', {
        backend: this.backend,
        id,
      })
      .catch(error => {
        // there was an error likely in read metadata.
        // still load the page and handle what you show by filtering for this property
        if (error.httpStatus === 403) {
          this.noReadAccess = true;
        }
      });
  }

  setupController(controller, model) {
    controller.set('backend', this.backend); // for backendCrumb
    controller.set('id', model.id); // for navigation on tabs
    controller.set('model', model);
    controller.set('noReadAccess', model.noReadAccess);
    controller.set('currentVersion', model.currentVersion);
  }
}
