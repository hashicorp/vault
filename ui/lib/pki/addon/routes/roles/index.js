import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class RolesIndexRoute extends Route {
  @service store;
  @service secretMountPath;

  async model() {
    let response = await this.store
      .query('pki/pki-role-engine', { backend: this.secretMountPath.currentPath })
      .catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
    return response;
  }
}
