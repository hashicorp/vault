import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class PkiKeysIndexRoute extends Route {
  @service store;
  @service secretMountPath;
  @service pathHelp;

  beforeModel() {
    // Must call this promise before the model hook otherwise it doesn't add OpenApi to record.
    return this.pathHelp.getNewModel('pki/key', this.secretMountPath.currentPath);
  }

  model() {
    return this.store
      .query('pki/key', { backend: this.secretMountPath.currentPath })
      .then((keyModel) => {
        return { keyModel, parentModel: this.modelFor('keys') };
      })
      .catch((err) => {
        if (err.httpStatus === 404) {
          return { parentModel: this.modelFor('keys') };
        } else {
          throw err;
        }
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
