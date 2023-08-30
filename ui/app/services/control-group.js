/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { inject as service } from '@ember/service';
import RSVP from 'rsvp';
import ControlGroupError from 'vault/lib/control-group-error';
import getStorage from 'vault/lib/token-storage';
import parseURL from 'core/utils/parse-url';

const CONTROL_GROUP_PREFIX = 'vault:cg-';
const TOKEN_SEPARATOR = '☃';

// list of endpoints that return wrapped responses
// without `wrap-ttl`
const WRAPPED_RESPONSE_PATHS = [
  'sys/wrapping/rewrap',
  'sys/wrapping/wrap',
  'sys/replication/performance/primary/secondary-token',
  'sys/replication/dr/primary/secondary-token',
];

const storageKey = (accessor, path) => {
  return `${CONTROL_GROUP_PREFIX}${accessor}${TOKEN_SEPARATOR}${path}`;
};

export { storageKey, CONTROL_GROUP_PREFIX, TOKEN_SEPARATOR };
export default Service.extend({
  version: service(),
  router: service(),

  storage() {
    return getStorage();
  },

  keyFromAccessor(accessor) {
    const keys = this.storage().keys() || [];
    const returnKey = keys
      .filter((k) => k.startsWith(CONTROL_GROUP_PREFIX))
      .find((key) => key.replace(CONTROL_GROUP_PREFIX, '').startsWith(accessor));
    return returnKey ? returnKey : null;
  },

  storeControlGroupToken(info) {
    const key = storageKey(info.accessor, info.creation_path);
    this.storage().setItem(key, info);
  },

  deleteControlGroupToken(accessor) {
    this.unmarkTokenForUnwrap();
    const key = this.keyFromAccessor(accessor);
    this.storage().removeItem(key);
  },

  deleteTokens() {
    const keys = this.storage().keys() || [];
    keys.filter((k) => k.startsWith(CONTROL_GROUP_PREFIX)).forEach((key) => this.storage().removeItem(key));
  },

  wrapInfoForAccessor(accessor) {
    const key = this.keyFromAccessor(accessor);
    return key ? this.storage().getItem(key) : null;
  },

  tokenToUnwrap: null,
  markTokenForUnwrap(accessor) {
    this.set('tokenToUnwrap', this.wrapInfoForAccessor(accessor));
  },

  unmarkTokenForUnwrap() {
    this.set('tokenToUnwrap', null);
  },

  tokenForUrl(url) {
    if (this.version.isOSS) {
      return null;
    }
    let pathForUrl = parseURL(url).pathname;
    pathForUrl = pathForUrl.replace('/v1/', '');
    const tokenInfo = this.tokenToUnwrap;
    if (tokenInfo && tokenInfo.creation_path === pathForUrl) {
      const { token, accessor, creation_time } = tokenInfo;
      return { token, accessor, creationTime: creation_time };
    }
    return null;
  },

  checkForControlGroup(callbackArgs, response, wasWrapTTLRequested) {
    const creationPath = response && response?.wrap_info?.creation_path;
    if (
      this.version.isOSS ||
      wasWrapTTLRequested ||
      !response ||
      (creationPath && WRAPPED_RESPONSE_PATHS.includes(creationPath)) ||
      !response.wrap_info
    ) {
      return RSVP.resolve(...callbackArgs);
    }
    const error = new ControlGroupError(response.wrap_info);
    return RSVP.reject(error);
  },

  handleError(error) {
    const { accessor, token, creation_path, creation_time, ttl } = error;
    const data = { accessor, token, creation_path, creation_time, ttl };
    data.uiParams = { url: this.router.currentURL };
    this.storeControlGroupToken(data);
    return this.router.transitionTo('vault.cluster.access.control-group-accessor', accessor);
  },

  // Handle error from non-read request (eg. POST or UPDATE) so it can be retried
  saveTokenFromError(error) {
    const { accessor, token, creation_path, creation_time, ttl } = error;
    const data = { accessor, token, creation_path, creation_time, ttl };
    this.storeControlGroupToken(data);
    // In the read flow the accessor is marked once the user clicks "Visit" from the control group page
    // On a POST/UPDATE flow we don't redirect, so we need to mark automatically so that on the next try
    // the request will attempt unwrap.
    this.markTokenForUnwrap(accessor);
  },

  logFromError(error) {
    const { accessor, token, creation_path, creation_time, ttl } = error;
    const data = { accessor, token, creation_path, creation_time, ttl };
    this.storeControlGroupToken(data);

    const href = this.router.urlFor('vault.cluster.access.control-group-accessor', accessor);
    const lines = [
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
