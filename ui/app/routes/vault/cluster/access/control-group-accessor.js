import Ember from 'ember';
import UnloadModel from 'vault/mixins/unload-model-route';
const { inject } = Ember;

export default Ember.Route.extend(UnloadModel, {
  version: inject.service(),

  beforeModel() {
    return this.get('version').fetchFeatures().then(() => {
      return this._super(...arguments);
    });
  },

  model(params) {
    return this.get('version').hasFeature('Control Groups')
      ? this.store.findRecord('control-group', params.accessor)
      : null;
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
