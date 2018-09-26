import { set } from '@ember/object';
import Route from '@ember/routing/route';
import DS from 'ember-data';

const MODEL_FROM_PARAM = {
  entities: 'entity',
  groups: 'group',
};

export default Route.extend({
  model(params) {
    let model = MODEL_FROM_PARAM[params.item_type];
    if (!model) {
      const error = new DS.AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }
    return model;
  },
});
