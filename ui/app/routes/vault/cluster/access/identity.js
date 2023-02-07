import AdapterError from '@ember-data/adapter/error';
import { set } from '@ember/object';
import Route from '@ember/routing/route';

const MODEL_FROM_PARAM = {
  entities: 'entity',
  groups: 'group',
};

export default Route.extend({
  model(params) {
    const model = MODEL_FROM_PARAM[params.item_type];
    if (!model) {
      const error = new AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }
    return model;
  },
});
