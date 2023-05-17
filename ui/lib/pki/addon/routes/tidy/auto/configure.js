import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
// import { withConfirmLeave } from 'core/decorators/confirm-leave';

// @withConfirmLeave()
export default class PkiTidyAutoConfigureRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const tidyRouteModel = this.modelFor('tidy');
    return tidyRouteModel.autoTidyConfig;
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'configuration', route: 'configuration.index' },
      { label: 'tidy', route: 'tidy' },
      { label: 'auto-tidy configuration', route: 'tidy.auto' },
      { label: 'configure auto-tidy' },
    ];
  }
}
