import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class PkiKeysIndexRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    return this.store
      .query('pki/pki-key-engine', { backend: this.secretMountPath.currentPath })
      .catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
  }
}
