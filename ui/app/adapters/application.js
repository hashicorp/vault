import Ember from 'ember';
import DS from 'ember-data';
import fetch from 'fetch';

const POLLING_URL_PATTERNS = ['sys/seal-status', 'sys/health', 'sys/replication/status'];
const { inject, assign, set, RSVP } = Ember;

export default DS.RESTAdapter.extend({
  auth: inject.service(),
  controlGroup: inject.service(),

  flashMessages: inject.service(),

  namespace: 'v1/sys',

  shouldReloadAll() {
    return true;
  },

  shouldReloadRecord() {
    return true;
  },

  shouldBackgroundReloadRecord() {
    return false;
  },

  _preRequest(url, options) {
    const token = options.clientToken || this.get('auth.currentToken');
    if (token && !options.unauthenticated) {
      options.headers = assign(options.headers || {}, {
        'X-Vault-Token': token,
      });
      if (options.wrapTTL) {
        assign(options.headers, { 'X-Vault-Wrap-TTL': options.wrapTTL });
      }
    }
    const isPolling = POLLING_URL_PATTERNS.some(str => url.includes(str));
    if (!isPolling) {
      this.get('auth').setLastFetch(Date.now());
    }
    if (this.get('auth.shouldRenew')) {
      this.get('auth').renew();
    }
    options.timeout = 60000;
    return options;
  },

  ajax(intendedUrl, method, passedOptions = {}) {
    let url = intendedUrl;
    let type = method;
    let options = passedOptions;
    let controlGroup = this.get('controlGroup');
    let controlGroupToken = controlGroup.tokenForUrl(url);
    // if we have a control group token that matches the intendedUrl,
    // then we want to unwrap it and return the unwrapped response as
    // if it were the initial request
    // To do this, we rewrite the function args
    if (controlGroupToken) {
      url = '/v1/sys/wrapping/unwrap';
      type = 'POST';
      options = {
        clientToken: controlGroupToken.token,
        data: {
          token: controlGroupToken.token,
        },
      };
    }
    let opts = this._preRequest(url, options);

    return this._super(url, type, opts).then((...args) => {
      if (controlGroupToken) {
        controlGroup.deleteControlGroupToken(controlGroupToken.accessor);
      }
      const [resp] = args;
      if (resp && resp.warnings) {
        let flash = this.get('flashMessages');
        resp.warnings.forEach(message => {
          flash.info(message);
        });
      }
      return controlGroup.checkForControlGroup(args, resp, options.wrapTTL);
    });
  },

  // for use on endpoints that don't return JSON responses
  rawRequest(url, type, options = {}) {
    let opts = this._preRequest(url, options);
    return fetch(url, {
      method: type | 'GET',
      headers: opts.headers | {},
    }).then(response => {
      if (response.status >= 200 && response.status < 300) {
        return RSVP.resolve(response);
      } else {
        return RSVP.reject();
      }
    });
  },

  handleResponse(status, headers, payload, requestData) {
    const returnVal = this._super(...arguments);
    // ember data errors don't have the status code, so we add it here
    if (returnVal instanceof DS.AdapterError) {
      set(returnVal, 'httpStatus', status);
      set(returnVal, 'path', requestData.url);
    }
    return returnVal;
  },
});
