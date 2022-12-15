import PkiOverviewRoute from '../overview';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
export default class PkiKeysIndexRoute extends PkiOverviewRoute {
  @service store;
  @service secretMountPath;
  @service pathHelp;

  beforeModel() {
    // Must call this promise before the model hook otherwise it doesn't add OpenApi to record.
    return this.pathHelp.getNewModel('pki/key', this.secretMountPath.currentPath);
  }

  model() {
    return hash({
      hasConfig: this.hasConfig(),
      parentModel: this.modelFor('keys'),
      keyModels: this.store.query('pki/key', { backend: this.secretMountPath.currentPath }).catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      }),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const backend = this.secretMountPath.currentPath || 'pki';
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: backend, route: 'overview' },
      { label: 'keys', route: 'keys.index' },
    ];
  }
}
