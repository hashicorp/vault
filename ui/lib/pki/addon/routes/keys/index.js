import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class PkiKeysIndexRoute extends Route {
  @service store;
  @service secretMountPath;
  @service pathHelp;

  beforeModel() {
    // Must call this promise before the model hook otherwise it doesn't add OpenApi to record.
    return this.pathHelp.getNewModel('pki/key', 'pki');
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
}
