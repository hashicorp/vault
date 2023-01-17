import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class PkiCertificatesIndexRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    return this.store
      .query('pki/certificate/base', { backend: this.secretMountPath.currentPath })
      .then((certificates) => {
        return { certificates, parentModel: this.modelFor('certificates') };
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
