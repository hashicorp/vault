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
      let controlGroup = this.get('controlGroup');
      if (err instanceof ControlGroupError) {
        return controlGroup.handleError(err, transition);
      }
      if (err.path === '/v1/sys/wrapping/unwrap') {
        controlGroup.unmarkTokenForUnwrap();
      }
      return true;
    },
  },
});
