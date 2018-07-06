import Ember from 'ember';

export default Ember.Route.extend({
  resetController(controller) {
    controller.set('accessor', null);
  },
});
