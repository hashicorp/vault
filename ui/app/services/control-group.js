import Ember from 'ember';
import ControlGroupError from 'vault/lib/control-group-error';

const { Service, assign, inject, RSVP } = Ember;

export default Service.extend({
  version: inject.service(),
  router: inject.service(),

  checkForControlGroup(callbackArgs, response, wasWrapTTLRequested) {
    if (this.get('version.isOSS') || wasWrapTTLRequested || !response.wrap_info) {
      return RSVP.resolve(...callbackArgs);
    }
    let error = new ControlGroupError();
    error = assign(error, response.wrap_info);
    return RSVP.reject(error);
  },

  handleError(error, transition) {

    //console requests won't have a transition
    if (transition) {}
    debugger;
  }

});
