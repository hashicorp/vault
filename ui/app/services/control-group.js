import Ember from 'ember';
import ControlGroupError from 'vault/lib/control-group-error';

const { Service, inject, RSVP } = Ember;

// list of endpoints that return wrapped responses
// without `wrap-ttl`
const WRAPPED_RESPONSE_PATHS = [
  'sys/wrapping/rewrap',
  'sys/wrapping/wrap',
  'sys/replication/performance/primary/secondary-token',
  'sys/replication/dr/primary/secondary-token',
];

export default Service.extend({
  version: inject.service(),
  router: inject.service(),

  checkForControlGroup(callbackArgs, response, wasWrapTTLRequested) {
    let creationPath = response && Ember.get(response, 'wrap_info.creation_path');
    if (
      this.get('version.isOSS') ||
      wasWrapTTLRequested ||
      !response ||
      (creationPath && WRAPPED_RESPONSE_PATHS.includes(creationPath)) ||
      !response.wrap_info
    ) {
      return RSVP.resolve(...callbackArgs);
    }
    let error = new ControlGroupError(response.wrap_info);
    return RSVP.reject(error);
  },

  handleError(error) {
    let {accessor} = error;
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
