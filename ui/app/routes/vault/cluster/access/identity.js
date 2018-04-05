import Ember from 'ember';
import DS from 'ember-data';

const MODEL_FROM_PARAM = {
  entities: 'entity',
  groups: 'group',
};

export default Ember.Route.extend({
  model(params) {
    let model = MODEL_FROM_PARAM[params.item_type];
    if (!model) {
      const error = new DS.AdapterError();
      Ember.set(error, 'httpStatus', 404);
      throw error;
    }
    return model;
  },
});
