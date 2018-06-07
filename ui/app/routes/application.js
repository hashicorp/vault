import Ember from 'ember';

const { inject } = Ember;
export default Ember.Route.extend({
  controlGroup: inject.service(),

  actions: {
    willTransition() {
      window.scrollTo(0, 0);
    },
    error(err, transition) {
      this.get('controlGroup').handleError(err, transition);
    }
  },
});
