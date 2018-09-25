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
    let type = 'control-group-config';
    return this.get('version').hasFeature('Control Groups')
      ? this.store.findRecord(type, 'config').catch(e => {
          // if you haven't saved a config, the API 404s, so create one here to edit and return it
          if (e.httpStatus === 404) {
            return this.store.createRecord(type, {
              id: 'config',
            });
          }
          throw e;
        })
      : null;
  },

  actions: {
    reload() {
      this.refresh();
    },
  },
});
