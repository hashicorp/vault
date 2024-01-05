/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import RESTAdapter from '@ember-data/adapter/rest';
import { inject as service } from '@ember/service';
import { assign } from '@ember/polyfills';
import { set } from '@ember/object';
import RSVP from 'rsvp';
import config from '../config/environment';
import fetch from 'fetch';

const { APP } = config;
const { POLLING_URLS, NAMESPACE_ROOT_URLS } = APP;

export default RESTAdapter.extend({
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

  addHeaders(url, options, method) {
    const token = options.clientToken || this.auth.currentToken;
    const headers = {};
    if (token && !options.unauthenticated) {
      headers['X-Vault-Token'] = token;
    }
    if (options.wrapTTL) {
      headers['X-Vault-Wrap-TTL'] = options.wrapTTL;
    }
    if (method === 'PATCH') {
      headers['Content-Type'] = 'application/merge-patch+json';
    }
    const namespace =
      typeof options.namespace === 'undefined' ? this.namespaceService.path : options.namespace;
    if (namespace && !NAMESPACE_ROOT_URLS.some((str) => url.includes(str))) {
      headers['X-Vault-Namespace'] = namespace;
    }
    options.headers = assign(options.headers || {}, headers);
  },

  _preRequest(url, options, method) {
    this.addHeaders(url, options, method);
    const isPolling = POLLING_URLS.some((str) => url.includes(str));
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
    const controlGroup = this.controlGroup;
    const controlGroupToken = controlGroup.tokenForUrl(url);
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
    const opts = this._preRequest(url, options, method);

    return this._super(url, type, opts).then((...args) => {
      if (controlGroupToken) {
        controlGroup.deleteControlGroupToken(controlGroupToken.accessor);
      }
      const [resp] = args;
      if (resp && resp.warnings) {
        const flash = this.flashMessages;
        resp.warnings.forEach((message) => {
          flash.info(message);
        });
      }
      return controlGroup.checkForControlGroup(args, resp, options.wrapTTL);
    });
  },

  // for use on endpoints that don't return JSON responses
  rawRequest(url, type, options = {}) {
    const opts = this._preRequest(url, options);
    return fetch(url, {
      method: type || 'GET',
      headers: opts.headers || {},
      body: opts.body,
      signal: opts.signal,
    }).then((response) => {
      if (response.status >= 200 && response.status < 300) {
        return RSVP.resolve(response);
      } else {
        return RSVP.reject(response);
      }
    });
  },

  handleResponse(status, headers, payload, requestData) {
    const returnVal = this._super(...arguments);
    if (returnVal instanceof AdapterError) {
      // ember data errors don't have the status code, so we add it here
      set(returnVal, 'httpStatus', status);
      set(returnVal, 'path', requestData.url);
      // Most of the time when the Vault API returns an error, the payload looks like:
      // { errors: ['some error message']}
      // But sometimes (eg RespondWithStatusCode) it looks like this:
      // { data: { error: 'some error message' } }
      if (payload?.data?.error && !payload.errors) {
        // Normalize the errors from RespondWithStatusCode
        set(returnVal, 'errors', [payload.data.error]);
      }
    }
    return returnVal;
  },
});
