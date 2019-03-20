import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import UnloadModel from 'vault/mixins/unload-model-route';

export default Route.extend(UnloadModel, {
  version: service(),

  beforeModel() {
    return this.get('version')
      .fetchFeatures()
      .then(() => {
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
