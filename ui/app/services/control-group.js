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

  keyFromAccessor(accessor) {
    let keys = this.storage().keys() || [];
    let returnKey = keys
      .filter(k => k.startsWith(CONTROL_GROUP_PREFIX))
      .find(key => key.replace(CONTROL_GROUP_PREFIX, '').startsWith(accessor));
    return returnKey ? returnKey : null;
  },

  storeControlGroupToken(info) {
    let key = this.storageKey(info.accessor, info.creation_path);
    this.storage().setItem(key, info);
  },

  deleteControlGroupToken(accessor) {
    this.set('tokenToUnwrap', null);
    let key = this.keyFromAccessor(accessor);
    return key ? this.storage().removeItem(key) : null;
  },

  wrapInfoForAccessor(accessor) {
    let key = this.keyFromAccessor(accessor);
    return key ? this.storage().getItem(key) : null;
  },

  tokenToUnwrap: null,
  markTokenForUnwrap(accessor) {
    this.set('tokenToUnwrap', this.wrapInfoForAccessor(accessor));
  },
  tokenForUrl(url) {
    if (this.get('version.isOSS')) {
      return null;
    }
    let pathForUrl = url.replace('/v1/', '');
    let tokenInfo = this.get('tokenToUnwrap');
    if (tokenInfo && tokenInfo.creation_path === pathForUrl) {
      let { token, accessor } = tokenInfo;
      return { token, accessor };
    }
    return null;
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
    let { url, name, contexts, queryParams } = transition.intent;
    let data = { accessor, token, creation_path, creation_time, ttl };
    data.uiParams = { url, name, contexts, queryParams };
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
