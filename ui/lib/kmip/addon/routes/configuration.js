import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default Route.extend({
  store: service(),
  secretMountPath: service(),
  pathHelp: service(),
  beforeModel() {
    return this.pathHelp.getNewModel('kmip/config', this.secretMountPath.currentPath);
  },
  model() {
    return this.store.findRecord('kmip/config', this.secretMountPath.currentPath).catch(err => {
      if (err.httpStatus === 404) {
        return;
      } else {
        throw err;
      }
    });
  },
});
