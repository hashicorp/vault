import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class PkiCertificatesIndexRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    return this.store
      .query('pki/pki-certificate-engine', { backend: this.secretMountPath.currentPath })
      .catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
  }
}
