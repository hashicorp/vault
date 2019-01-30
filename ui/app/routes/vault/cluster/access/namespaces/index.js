import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import UnloadModel from 'vault/mixins/unload-model-route';

export default Route.extend(UnloadModel, {
  version: service(),
  beforeModel() {
    this.store.unloadAll('namespace');
    return this.get('version')
      .fetchFeatures()
      .then(() => {
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
