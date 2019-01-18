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

  model() {
    return this.get('version').hasFeature('Control Groups') ? this.store.createRecord('control-group') : null;
  },
});
