import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class PkiIssuersIndexRoute extends Route {
  @service store;
  @service secretMountPath;
  @service pathHelp;

  async model() {
    // the pathHelp service is needed for adding openAPI to the model
    await this.pathHelp.getNewModel('pki/pki-issuer-engine', 'pki');

    return await this.store
      .query('pki/pki-issuer-engine', { backend: this.secretMountPath.currentPath })
      .catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
  }
}
