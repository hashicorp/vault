import Ember from 'ember';
import ControlGroupError from 'vault/lib/control-group-error';
import getStorage from 'vault/lib/token-storage';

const CONTROL_GROUP_PREFIX = 'vault:cg-';
const TOKEN_SEPARATOR = '☃';
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

  storage() {
    return getStorage();
  },

  storageKey(accessor, path) {
    return `${CONTROL_GROUP_PREFIX}${accessor}${TOKEN_SEPARATOR}${path}`;
  },

  storeControlGroupToken(info) {
    let key = this.storageKey(info.accessor, info.creation_path);
    this.storage().setItem(key, info);
  },

  wrapInfoForAccessor(accessor) {
    let keys = this.storage().keys() || [];
    let returnKey = keys
      .filter(k => k.startsWith(CONTROL_GROUP_PREFIX))
      .find(key => key.replace(CONTROL_GROUP_PREFIX, '').startsWith(accessor));
    return returnKey ? this.storage().getItem(returnKey) : null;
  },

  hasTokenForUrl(url) {
    let pathForUrl = url.replace('/v1/', '');
    if (this.get('version.isOSS')) {
      return false;
    }
  },

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

  handleError(error, transition) {
    let { accessor, token, creation_path, creation_time, ttl } = error;
    let { name, contexts, queryParams } = transition.intent;
    let data = { accessor, token, creation_path, creation_time, ttl };
    data.uiParams = { name, contexts, queryParams };
    this.storeControlGroupToken(data);
    return this.get('router').transitionTo('vault.cluster.access.control-group-accessor', accessor);
  },

  logFromError(error) {
    let { accessor, token, creation_path, creation_time, ttl } = error;
    let data = { accessor, token, creation_path, creation_time, ttl };
    this.storeControlGroupToken(data);

    let href = this.get('router').urlFor('vault.cluster.access.control-group-accessor', accessor);
    let lines = [
      `A Control Group was encountered at ${error.creation_path}.`,
      `The Control Group Token is ${error.token}.`,
      `The Accessor is ${error.accessor}.`,
      `Visit <a href='${href}'>${href}</a> for more details.`,
    ];
    return {
      type: 'error-with-html',
      content: lines.join('\n'),
    };
  },
});
