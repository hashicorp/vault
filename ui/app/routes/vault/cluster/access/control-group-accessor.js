import Ember from 'ember';

import UnloadModel from 'vault/mixins/unload-model-route';

export default Ember.Route.extend(UnloadModel, {
  model(params) {
    return this.store.findRecord('control-group', params.accessor);
  },
  actions: {
    willTransition() {
      return true;
    },
    // deactivate happens later than willTransition,
    // so since we're using the model to render links
    // we don't want the UI blinking
    deactivate() {
      this.unloadModel();
      return true;
    },
  },
});
