import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default Route.extend({
  store: service(),
  secretMountPath: service(),
  beforeModel() {
    this.store.unloadAll('kmip/scope');
  },
  model() {
    let model = this.store.createRecord('kmip/scope', {
      backend: this.secretMountPath.currentPath,
    });
    return model;
  },
});
