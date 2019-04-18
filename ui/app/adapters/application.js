import { inject as service } from '@ember/service';
import { assign } from '@ember/polyfills';
import { set } from '@ember/object';
import RSVP from 'rsvp';
import DS from 'ember-data';
import fetch from 'fetch';
import config from '../config/environment';

const { APP } = config;
const { POLLING_URLS, NAMESPACE_ROOT_URLS } = APP;

export default DS.RESTAdapter.extend({
  auth: service(),
  namespaceService: service('namespace'),
  controlGroup: service(),

  flashMessages: service(),

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

  addHeaders(url, options) {
    let token = options.clientToken || this.get('auth.currentToken');
    let headers = {};
    if (token && !options.unauthenticated) {
      headers['X-Vault-Token'] = token;
    }
    if (options.wrapTTL) {
      headers['X-Vault-Wrap-TTL'] = options.wrapTTL;
    }
    let namespace =
      typeof options.namespace === 'undefined' ? this.get('namespaceService.path') : options.namespace;
    if (namespace && !NAMESPACE_ROOT_URLS.some(str => url.includes(str))) {
      headers['X-Vault-Namespace'] = namespace;
    }
    options.headers = assign(options.headers || {}, headers);
  },

  _preRequest(url, options) {
    this.addHeaders(url, options);
    const isPolling = POLLING_URLS.some(str => url.includes(str));
    if (!isPolling) {
      this.auth.setLastFetch(Date.now());
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
    // if we have a Control Group token that matches the intendedUrl,
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
      method: type || 'GET',
      headers: opts.headers || {},
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
