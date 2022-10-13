import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class PkiIssuersIndexRoute extends Route {
  @service store;
  @service secretMountPath;
  @service pathHelp;

  beforeModel() {
    // Must call this promise before the model hook otherwise it doesn't add OpenApi to record.
    return this.pathHelp.getNewModel('pki/pki-issuer-engine', 'pki');
  }

  model() {
    // the pathHelp service is needed for adding openAPI to the model
    this.pathHelp.getNewModel('pki/pki-issuer-engine', 'pki');

    return this.store
      .query('pki/pki-issuer-engine', { backend: this.secretMountPath.currentPath })
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
}
