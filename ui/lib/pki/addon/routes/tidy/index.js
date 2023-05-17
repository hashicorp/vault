import Route from '@ember/routing/route';
import { PKI_DEFAULT_EMPTY_STATE_MSG } from '../overview';
import { hash } from 'rsvp';
import { inject as service } from '@ember/service';

export default class PkiTidyIndexRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const adapter = this.store.adapterFor('application');
    const { hasConfig, autoTidyConfig, engine } = this.modelFor('tidy');

    return hash({
      tidyStatus: adapter.ajax(`/v1/${this.secretMountPath.currentPath}/tidy-status`, 'GET'),
      hasConfig,
      autoTidyConfig,
      engine,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.notConfiguredMessage = PKI_DEFAULT_EMPTY_STATE_MSG;
  }
}
