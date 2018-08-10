import Ember from 'ember';
import UnloadModel from 'vault/mixins/unload-model-route';

export default Ember.Route.extend(UnloadModel, {
  beforeModel() {
    return this.store.unloadAll('namespace');
  },
  model() {
    return this.store.findAll('namespace').catch(e => {
      if (e.httpStatus === 404) {
        return [];
      }
      throw e;
    });
  },
});
