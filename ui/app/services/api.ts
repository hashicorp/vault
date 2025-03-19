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
} from '@hashicorp/vault-client-typescript';
import config from '../config/environment';

import type AuthService from 'vault/services/auth';
import type NamespaceService from 'vault/services/namespace';
import type ControlGroupService from 'vault/services/control-group';
import type FlashMessageService from 'vault/services/flash-messages';
import type { ApiError, ApiResponse, HeaderMap, XVaultHeaders } from 'vault/api';

export default class ApiService extends Service {
  @service('auth') declare readonly authService: AuthService;
  @service('namespace') declare readonly namespaceService: NamespaceService;
  @service declare readonly controlGroup: ControlGroupService;
  @service declare readonly flashMessages: FlashMessageService;

  #responseCache = new Map<string, ApiResponse>();

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
    // if we have a Control Group token that matches the intendedUrl,
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
    if (!headers.has('X-Vault-Token')) {
      headers.set('X-Vault-Token', this.authService.currentToken);
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
  showWarnings = async (context: ResponseContext) => {
    const response = context.response.clone();
    const json = await response?.json();
    // currently there is only 1 endpoint that intentionally hides warnings
    // handle that here for now, but if more endpoints are added, look for a more scalable pattern
    const hideWarnings = context.url === '/v1/sys/internal/counters/activity';

    if (json?.warnings && !hideWarnings) {
      json.warnings.forEach((message: string) => {
        this.flashMessages.info(message);
      });
    }
  };

  deleteControlGroupToken = async (context: ResponseContext) => {
    const { url } = context;
    const controlGroupToken = this.controlGroup.tokenForUrl(url);
    if (controlGroupToken) {
      this.controlGroup.deleteControlGroupToken(controlGroupToken.accessor);
    }
  };

  formatErrorResponse = async (context: ResponseContext) => {
    const response = context.response.clone();
    const { headers, status, statusText } = response;

    // backwards compatibility with Ember Data
    if (status >= 400) {
      const error: ApiError = (await response?.json()) || {};
      error.httpStatus = response?.status;
      error.path = context.url;
      // typically the Vault API error response looks like { errors: ['some error message'] }
      // but sometimes (eg RespondWithStatusCode) it's { data: { error: 'some error message' } }
      if (error?.data?.error && !error.errors) {
        // normalize the errors from RespondWithStatusCode
        error.errors = [error.data.error];
      }
      return new Response(JSON.stringify(error), { headers, status, statusText });
    }

    return;
  };

  // the responses in the OpenAPI spec don't account for the return values to be under the 'data' key
  // return the data rather than the entire response
  // furthermore, some requests require the full response to access things like wrap_info for example
  // if the response type in the OpenAPI spec is void then undefined is returned regardless of what is returned here
  // cache the response to be accessed in overridden handlers
  extractData = async (context: ResponseContext) => {
    const response = context.response.clone();
    const { headers, status, statusText } = response;

    if (status >= 200 && status < 300) {
      const json = await response?.json();
      // roughing this in for now more as a placeholder
      // thinking about how to make the outer level response data like wrap_info accessible
      if (!json.data) {
        this.#responseCache.set(context.url, json);
      }
      return new Response(JSON.stringify(json.data), { headers, status, statusText });
    }

    return;
  };
  // --- End Middleware ---

  configuration = new Configuration({
    basePath: '/v1',
    middleware: [
      { pre: this.setLastFetch },
      { pre: this.getControlGroupToken },
      { pre: this.setHeaders },
      { post: this.showWarnings },
      { post: this.deleteControlGroupToken },
      { post: this.formatErrorResponse },
      { post: this.extractData },
    ],
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
}
