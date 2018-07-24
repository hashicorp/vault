import Ember from 'ember';
import UnloadModel from 'vault/mixins/unload-model-route';

export default Ember.Route.extend(UnloadModel, {
  model() {
    return this.store.createRecord('namespace');
  },
});
