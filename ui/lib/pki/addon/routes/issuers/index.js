import PkiOverviewRoute from '../overview';
import { inject as service } from '@ember/service';

export default class PkiIssuersListRoute extends PkiOverviewRoute {
  @service store;
  @service secretMountPath;
  @service pathHelp;

  beforeModel() {
    // Must call this promise before the model hook otherwise it doesn't add OpenApi to record.
    return this.pathHelp.getNewModel('pki/issuer', this.secretMountPath.currentPath);
  }

  model() {
    return this.store
      .query('pki/issuer', { backend: this.secretMountPath.currentPath })
      .then((issuersModel) => {
        return { issuersModel, parentModel: this.modelFor('issuers') };
      })
      .catch((err) => {
        if (err.httpStatus === 404) {
          return { parentModel: this.modelFor('issuers') };
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
      { label: 'issuers', route: 'issuers.index' },
    ];
  }
}
