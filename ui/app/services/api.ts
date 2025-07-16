/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { service } from '@ember/service';
import {
  Configuration,
  RequestContext,
  ResponseContext,
  AuthApi,
  IdentityApi,
  SecretsApi,
  SystemApi,
  HTTPQuery,
  HTTPRequestInit,
  RequestOpts,
  ResponseError,
} from '@hashicorp/vault-client-typescript';
import config from 'vault/config/environment';
import { waitForPromise } from '@ember/test-waiters';
import { underscore } from '@ember/string';

import type AuthService from 'vault/services/auth';
import type NamespaceService from 'vault/services/namespace';
import type ControlGroupService from 'vault/services/control-group';
import type FlashMessageService from 'vault/services/flash-messages';
import type { HeaderMap, XVaultHeaders } from 'vault/api';

export default class ApiService extends Service {
  @service('auth') declare readonly authService: AuthService;
  @service('namespace') declare readonly namespaceService: NamespaceService;
  @service declare readonly controlGroup: ControlGroupService;
  @service declare readonly flashMessages: FlashMessageService;

  // -- Pre Request Middleware --
  setLastFetch = async (context: RequestContext) => {
    const { url } = context;
    const isPolling = config.APP.POLLING_URLS.some((str) => url.includes(str));

    if (!isPolling) {
      this.authService.setLastFetch(Date.now());
    }
  };

  getControlGroupToken = async (context: RequestContext) => {
    const { url, init } = context;
    const controlGroupToken = this.controlGroup.tokenForUrl(url);
    let newUrl = url;
    // if we have a Control Group token that matches the url,
    // unwrap it and return the unwrapped response as if it were the initial request
    // to do this, we rewrite the request
    if (controlGroupToken) {
      const { token } = controlGroupToken;
      const { headers } = this.buildHeaders({ token });
      newUrl = '/v1/sys/wrapping/unwrap';
      init.method = 'POST';
      init.headers = headers;
      init.body = JSON.stringify({ token });
    }
    return { url: newUrl, init };
  };

  setHeaders = async (context: RequestContext) => {
    const { url, init } = context;
    const headers = new Headers(init.headers);
    // unauthenticated or clientToken requests should set the header in initOverrides
    // unauthenticated value should be empty string, not undefined or null
    const { currentToken } = this.authService;
    if (!headers.has('X-Vault-Token') && currentToken) {
      headers.set('X-Vault-Token', currentToken);
    }
    if (init.method === 'PATCH') {
      headers.set('Content-Type', 'application/merge-patch+json');
    }
    // use initOverrides to set the namespace header to something other than path set in the namespace service
    // for requests that must be made to root namespace pass empty string as value
    const namespace = this.namespaceService.path;
    if (!headers.has('X-Vault-Namespace') && namespace) {
      headers.set('X-Vault-Namespace', namespace);
    }

    init.headers = headers;
    return { url, init };
  };

  normalizeRequestBodyKeys = async (context: RequestContext) => {
    const { url, init } = context;
    if (init.body) {
      const convertKeys = (value: unknown): unknown => {
        const notAnObject = (obj: unknown) =>
          !['[object Object]', '[object Array]'].includes(Object.prototype.toString.call(obj));

        if (notAnObject(value)) {
          return value;
        }
        const json = value as Record<string, unknown>;
        // object could be an array, in which case convert keys if it contains objects
        if (Array.isArray(json)) {
          return json.map(convertKeys);
        }
        // convert object keys to snake_case
        return Object.keys(json).reduce((convertedJson: Record<string, unknown>, key) => {
          const value = json[key];
          // if the value is an object, convert those keys too
          const convertedValue = notAnObject(value) ? value : convertKeys(value);
          convertedJson[underscore(key)] = convertedValue;
          return convertedJson;
        }, {});
      };

      const requestBody = JSON.parse(init.body as string);
      const convertedBody = convertKeys(requestBody);
      init.body = JSON.stringify(convertedBody);
    }

    return { url, init };
  };

  // -- Post Request Middleware --
  showWarnings = async (context: ResponseContext) => {
    const response = context.response.clone();
    // if the response is empty, don't try to parse it
    if (response.headers.get('Content-Length')) {
      const json = await response.json();

      if (json?.warnings) {
        json.warnings.forEach((message: string) => {
          this.flashMessages.info(message);
        });
      }
    }
  };

  deleteControlGroupToken = async (context: ResponseContext) => {
    const { url } = context;
    const controlGroupToken = this.controlGroup.tokenForUrl(url);
    if (controlGroupToken) {
      this.controlGroup.deleteControlGroupToken(controlGroupToken.accessor);
    }
  };
  // --- End Middleware ---

  configuration = new Configuration({
    basePath: '/v1',
    middleware: [
      { pre: this.setLastFetch },
      { pre: this.getControlGroupToken },
      { pre: this.setHeaders },
      { pre: this.normalizeRequestBodyKeys },
      { post: this.showWarnings },
      { post: this.deleteControlGroupToken },
    ],
    fetchApi: (...args: [Request]) => {
      return waitForPromise(window.fetch(...args));
    },
  });

  auth = new AuthApi(this.configuration);
  identity = new IdentityApi(this.configuration);
  secrets = new SecretsApi(this.configuration);
  sys = new SystemApi(this.configuration);

  // convenience method for overriding headers for given requests to ensure consistency
  // eg. this.api.sys.wrap(data, { headers: { 'X-Vault-Wrap-TTL': wrap } });
  // -> this.api.sys.wrap(data, this.api.buildHeaders({ wrap }));
  buildHeaders(headerMap: HeaderMap) {
    const headers = {} as XVaultHeaders;

    for (const key in headerMap) {
      const headerKey = {
        namespace: 'X-Vault-Namespace',
        token: 'X-Vault-Token',
        wrap: 'X-Vault-Wrap-TTL',
      }[key] as keyof XVaultHeaders;

      headers[headerKey] = headerMap[key as keyof HeaderMap];
    }

    return { headers };
  }

  // convenience method for updating the query params object on the request context
  // eg. this.api.sys.uiConfigListCustomMessages(true, ({ context: { query } }) => { query.authenticated = true });
  // -> this.api.sys.uiConfigListCustomMessages(true, (context) => this.api.addQueryParams(context, { authenticated: true }));
  addQueryParams(requestContext: { init: HTTPRequestInit; context: RequestOpts }, params: HTTPQuery = {}) {
    const { context } = requestContext;
    context.query = { ...context.query, ...params };
  }

  // accepts an error response and returns { status, message, response, path }
  // message is built as error.errors joined with a comma, error.message or a fallback message
  // path is the url of the request, minus the origin -> /v1/sys/wrapping/unwrap
  async parseError(e: unknown, fallbackMessage = 'An error occurred, please try again') {
    if (e instanceof ResponseError) {
      const { status, url } = e.response;
      const error = await e.response.json();
      // typically the Vault API error response looks like { errors: ['some error message'] }
      // but sometimes (eg RespondWithStatusCode) it's { data: { error: 'some error message' } }
      const errors = error.data?.error && !error.errors ? [error.data.error] : error.errors;
      const message = errors && typeof errors[0] === 'string' ? errors.join(', ') : error.message;

      return {
        message: message || fallbackMessage,
        status,
        path: url.replace(document.location.origin, ''),
        response: error,
      };
    }

    // log out generic error for ease of debugging in dev env
    if (config.environment === 'development') {
      console.log('API Error:', e); // eslint-disable-line no-console
    }

    return {
      message: (e as Error)?.message || fallbackMessage,
    };
  }

  // accepts a list response as { keyInfo, keys } and returns a flat array of the keyInfo datum
  // to preserve the keys (unique identifiers) the value will be set on the datum as id
  keyInfoToArray(response: unknown = {}) {
    const { keyInfo, keys } = response as { keyInfo?: Record<string, unknown>; keys?: string[] };
    if (!keyInfo || !keys) {
      return [];
    }
    return keys.reduce(
      (arr, key) => {
        const datum = keyInfo[key];
        if (datum) {
          arr.push({ id: key, ...datum });
        }
        return arr;
      },
      [] as Record<string, unknown>[]
    );
  }

  // some responses return an object with a uuid as the key rather than an array
  // in most cases it is easier to work with an array with the uuid set as a property on the object
  // for example, internalUiListEnabledVisibleMounts returns an object like: { secret: { '/path/to/secret': { ... } } }
  // usage for above example -> this.api.objectToArray(response.secret, 'path');
  // this would return an array of objects like: [{ path: '/path/to/secret', ... }]
  responseObjectToArray<T extends object>(obj?: T, uuidKey?: string) {
    if (obj) {
      return Object.entries(obj).map(([key, value]) => {
        if (uuidKey) {
          return { [uuidKey]: key, ...value };
        }
        return { ...value };
      });
    }
    return [];
  }
}
