import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import RSVP from 'rsvp';

export default class RolesIndexRoute extends Route {
  @service store;
  @service secretMountPath;
  @service pathHelp;

  beforeModel() {
    // Must call this promise before the model hook otherwise it doesn't add OpenApi to record.
    return this.pathHelp.getNewModel('pki/pki-role-engine', 'pki');
  }

  model() {
    let parentModel = this.modelFor('roles');
    let model = this.store
      .query('pki/pki-role-engine', { backend: this.secretMountPath.currentPath })
      .catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
    return RSVP.hash({
      parentModel,
      model,
    });
  }
}
