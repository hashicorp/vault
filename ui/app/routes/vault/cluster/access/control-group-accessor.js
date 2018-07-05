import Ember from 'ember';
import DS from 'ember-data';

export default Ember.Route.extend({
  model(params) {
    let model = this.store.createRecord('control-group', {
      id: params.accessor,
    });
    return model.request();
  },
});
