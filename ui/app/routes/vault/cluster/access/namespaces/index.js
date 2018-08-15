import Ember from 'ember';
import UnloadModel from 'vault/mixins/unload-model-route';

const { inject } = Ember;

export default Ember.Route.extend(UnloadModel, {
  version: inject.service(),
  beforeModel() {
    this.store.unloadAll('namespace');
    return this.get('version').fetchFeatures().then(() => {
      return this._super(...arguments);
    });
  },
  model() {
    return this.get('version.hasNamespaces')
      ? this.store.findAll('namespace').catch(e => {
          if (e.httpStatus === 404) {
            return [];
          }
          throw e;
        })
      : null;
  },
});
