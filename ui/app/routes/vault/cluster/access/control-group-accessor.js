import Ember from 'ember';
import DS from 'ember-data';

import UnloadModel from 'vault/mixins/unload-model-route';

export default Ember.Route.extend(UnloadModel, {
  model(params) {
    return this.store.findRecord('control-group', params.accessor);
  },
});
