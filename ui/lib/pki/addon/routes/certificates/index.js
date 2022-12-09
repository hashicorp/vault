import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class PkiCertificatesIndexRoute extends Route {
  @service store;
  @service secretMountPath;
  @service pathHelp;

  beforeModel() {
    // Must call this promise before the model hook otherwise it doesn't add OpenApi to record.
    return this.pathHelp.getNewModel('pki/certificate', this.secretMountPath.currentPath);
  }

  model() {
    return this.store
      .query('pki/certificate', { backend: this.secretMountPath.currentPath })
      .then((certificateModel) => {
        return { certificateModel, parentModel: this.modelFor('certificates') };
      })
      .catch((err) => {
        if (err.httpStatus === 404) {
          return { parentModel: this.modelFor('certificates') };
        } else {
          throw err;
        }
      });
  }
}
