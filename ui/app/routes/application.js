import Ember from 'ember';

export default Ember.Route.extend({
  actions: {
    willTransition() {
      window.scrollTo(0, 0);
    },
  },
});
