import PkiKeysIndexRoute from '.';
import { inject as service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

@withConfirmLeave()
export default class PkiKeysCreateRoute extends PkiKeysIndexRoute {
  @service store;

  model() {
    return this.store.createRecord('pki/key');
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs.push({ label: 'generate' });
  }
}
