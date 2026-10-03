/**
 * Copyright IBM Corp. 2016, 2025
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
import { waitForPromise, waitFor } from '@ember/test-waiters';

import type AuthService from 'vault/services/auth';
import type NamespaceService from 'vault/services/namespace';
import type ControlGroupService from 'vault/services/control-group';
import type FlashMessageService from 'vault/services/flash-messages';
import type { HeaderMap, XVaultHeaders } from 'vault/api';
import type { HTTPMethod } from '@hashicorp/vault-client-typescript';

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

  // -- Post Request Middleware --
  showWarnings = waitFor(async (context: ResponseContext) => {
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
  });

  checkControlGroup = waitFor(async (context: ResponseContext) => {
    const response = context.response.clone();
    const { headers } = response;

    // since control group requests are forwarded to /v1/sys/wrapping/unwrap we cannot use controlGroup.tokenForUrl here
    // instead, we can check if tokenToUnwrap exists on the service and compare the token value with the request header value
    if (this.controlGroup.tokenToUnwrap) {
      const { token, accessor } = this.controlGroup.tokenToUnwrap || {};
      const requestHeaders = context.init.headers as Headers;

      if (requestHeaders.get('X-Vault-Token') === token) {
        this.controlGroup.deleteControlGroupToken(accessor);
      }
    }
    // if the requested path is locked by a control group we need to create a new error response
    if (headers.get('Content-Length')) {
      const json = await response.json();
      const wrapTtl = headers.get('X-Vault-Wrap-TTL');
      const isLockedByControlGroup = this.controlGroup.isRequestedPathLocked(json, wrapTtl);

      if (isLockedByControlGroup) {
        const error = {
          message: 'Control Group encountered',
          isControlGroupError: true,
          ...json.wrap_info,
        };
        return new Response(JSON.stringify(error), { headers, status: 403, statusText: 'Forbidden' });
      }
    }

    return;
  });
  // --- End Middleware ---

  configuration = new Configuration({
    basePath: '/v1',
    middleware: [
      { pre: this.setLastFetch },
      { pre: this.getControlGroupToken },
      { pre: this.setHeaders },
      { post: this.showWarnings },
      { post: this.checkControlGroup },
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
        recoverSnapshotId: 'X-Vault-Recover-Snapshot-Id',
        recoverSourcePath: 'X-Vault-Recover-Source-Path',
      }[key] as keyof XVaultHeaders;

      headers[headerKey] = headerMap[key as keyof HeaderMap];
    }

    return { headers };
  }

  // convenience method for updating the query params object on the request context
  // eg. this.api.sys.uiConfigListCustomMessages(true, ({ context: { query } }) => { query.authenticated = true });
  // -> this.api.sys.uiConfigListCustomMessages(true, (context) => this.api.addQueryParams(context, { authenticated: true }));
  async addQueryParams(
    requestContext: { init: HTTPRequestInit; context: RequestOpts },
    params: HTTPQuery = {}
  ) {
    const { context, init } = requestContext;
    context.query = { ...context.query, ...params };
    return init;
  }

  // accepts an error response and returns { status, message, response, path }
  // message is built as error.errors joined with a comma, error.message or a fallback message
  // path is the url of the request, minus the origin -> /v1/sys/wrapping/unwrap
  parseError = waitFor(async (e: unknown, fallbackMessage = 'An error occurred, please try again') => {
    if (e instanceof ResponseError) {
      const { status, url } = e.response;
      // instances where an error is thrown multiple times could result in the body already being read
      // this will result in a readable stream failure and we can't parse the body
      // to avoid this, clone the response so we can access the body consistently
      const error = await e.response.clone().json();
      // typically the Vault API error response looks like { errors: ['some error message'] }
      // but sometimes (eg RespondWithStatusCode) it's { data: { error: 'some error message' } }
      const errors = error.data?.error && !error.errors ? [error.data.error] : error.errors;
      const message = errors && typeof errors[0] === 'string' ? errors.join(', ') : error.message;

      return {
        message: message || fallbackMessage,
        status,
        path: decodeURIComponent(url.replace(document.location.origin, '')),
        response: error,
      };
    }

    // log out generic error for ease of debugging in dev env
    if (config.environment === 'development') {
      console.error('API Error:', e);
    }

    return {
      message: (e as Error)?.message || fallbackMessage,
    };
  });

  // accepts a list response as { key_info, keys } and returns a flat array of the key_info datum
  // to preserve the keys (unique identifiers) the value will be set on the datum as the provided uuidKey or id
  keyInfoToArray(response: unknown = {}, uuidKey = 'id') {
    const { key_info, keys } = response as { key_info?: Record<string, unknown>; keys?: string[] };
    if (!key_info || !keys) {
      return [];
    }
    return keys.reduce(
      (arr, key) => {
        const datum = key_info[key];
        if (datum) {
          arr.push({ [uuidKey]: key, ...datum });
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

  // interface for making raw fetch requests outside of the generated API methods
  // this should only be used in cases where it's not possible to use the client
  private async rawRequest(path: string, method: HTTPMethod, body?: unknown) {
    const context = {
      url: `${this.configuration.basePath}${path}`,
      init: { method } as RequestInit,
    };

    if (body) {
      context.init.body = JSON.stringify(body);
    }

    const { url, init } = await this.setHeaders(context as RequestContext);

    const response = await this.configuration.fetchApi?.(new Request(url, init));
    if (!response?.ok) {
      throw response;
    }
    // with various content types like application/pem-certificate-chain or application/pkix-cert for example,
    // return the response so the caller can read the body with the appropriate method (blob, text, json etc.)
    return response;
  }
  request = {
    get: (path: string) => this.rawRequest(path, 'GET'),
    post: (path: string, body?: unknown) => this.rawRequest(path, 'POST', body),
    put: (path: string, body?: unknown) => this.rawRequest(path, 'PUT', body),
    patch: (path: string, body?: unknown) => this.rawRequest(path, 'PATCH', body),
    delete: (path: string, body?: unknown) => this.rawRequest(path, 'DELETE', body),
  };
}
