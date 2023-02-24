import PkiIssuersIndexRoute from '.';
import { inject as service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';
import { hash } from 'rsvp';

@withConfirmLeave()
export default class PkiIssuersCreateRoute extends PkiIssuersIndexRoute {
  @service store;

  queryParams = {
    formType: {
      refreshModel: true,
    },
  };

  beforeModel() {
    // pki/urls uses openApi to hydrate model
    return this.pathHelp.getNewModel('pki/urls', this.secretMountPath.currentPath);
  }

  async model(params) {
    return hash({
      config: this.store.createRecord('pki/action'),
      urls: this.getOrCreateUrls(this.secretMountPath.currentPath),
      formType: params.formType,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs.push({ label: 'create' });
  }

  async getOrCreateUrls(backend) {
    try {
      return this.store.findRecord('pki/urls', backend);
    } catch (e) {
      return this.store.createRecord('pki/urls', { id: backend });
    }
  }
}
