import Ember from 'ember';

const MODEL_FROM_PARAM = {
  entities: 'entity',
  groups: 'group',
};

export default Ember.Route.extend({
  model(params) {
    let model = MODEL_FROM_PARAM[params.item_type];

    //TODO 404 behavior;
    return model;
  },
});
