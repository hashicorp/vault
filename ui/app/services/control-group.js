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
    let {accessor} = error;
    if (transition) {
      transition.intent
    }
    return this.get('router')
      .transitionTo('vault.cluster.access.control-group-accessor', accessor);
  },

  logFromError(error) {
    let {accessor} = error;
    let href = this.get('router').urlFor('vault.cluster.access.control-group-accessor', accessor);
    let lines = [
      `A Control Group was encountered at ${error.creation_path}.`,
      `The Control Group Token is ${error.token}.`,
      `The Accessor is ${error.accessor}.`,
      `Visit <a href='${href}'>${href}</a> for more details.`
    ];
    return {
      type: 'error-with-html',
      content: lines.join('\n')
    };
  }

});
