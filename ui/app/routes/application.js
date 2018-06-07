import Ember from 'ember';
import ControlGroupError from 'vault/lib/control-group-error';

const { inject } = Ember;
export default Ember.Route.extend({
  controlGroup: inject.service(),

  actions: {
    willTransition() {
      window.scrollTo(0, 0);
    },
    error(err, transition) {
      if (err instanceof ControlGroupError) {
        return this.get('controlGroup').handleError(err, transition);
      }
      return true;
    }
  },
});
