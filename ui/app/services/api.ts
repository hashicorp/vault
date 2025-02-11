import Service, { service } from '@ember/service';
import {
  Configuration,
  HTTPHeaders,
  RequestContext,
  ResponseContext,
  ErrorContext,
  AuthApi,
  IdentityApi,
  SecretsApi,
  SystemApi,
} from '../api-client';
import config from '../config/environment';

import type AuthService from 'vault/vault/services/auth';
import type NamespaceService from 'vault/vault/services/namespace';

const {
  APP: { POLLING_URLS, NAMESPACE_ROOT_URLS },
} = config;

export { ErrorContext };

export interface ApiError {
  httpStatus: number;
  path: string;
  message: string;
  errors: Array<string | { [key: string]: unknown; title?: string; message?: string }>;
  data?: {
    [key: string]: unknown;
    error?: string;
  };
}

export default class ApiService extends Service {
  @service('auth') declare readonly authService: AuthService;
  @service('namespace') declare readonly namespaceService: NamespaceService;

  preRequest = async (context: RequestContext) => {
    const { url, init } = context;

    const isPolling = POLLING_URLS.some((str) => url.includes(str));
    if (!isPolling) {
      this.authService.setLastFetch(Date.now());
    }

    // since the client methods don't accept random options, perhaps these can be set on the authService instead
    const token = /* options.clientToken || */ this.authService.currentToken;
    const headers: HTTPHeaders = {};
    if (token /* && !options.unauthenticated */) {
      headers['X-Vault-Token'] = token;
    }
    // if (options.wrapTTL) {
    //   headers['X-Vault-Wrap-TTL'] = options.wrapTTL;
    // }
    if (init.method === 'PATCH') {
      headers['Content-Type'] = 'application/merge-patch+json';
    }
    // similarly to clientToken, perhaps this can always be set on the namespaceService?
    // const namespace =
    //   typeof options.namespace === 'undefined' ? this.namespaceService.path : options.namespace;
    const namespace = this.namespaceService.path;
    if (namespace && !NAMESPACE_ROOT_URLS.some((str) => url.includes(str))) {
      headers['X-Vault-Namespace'] = namespace;
    }

    Object.assign(init.headers || {}, headers);

    return { url, init };
  };

  postRequest = async (context: ResponseContext) => {
    const response: ApiError = (await context.response?.json()) || {};
    const { status } = context.response;

    // backwards compatibility with Ember Data
    if (status >= 400) {
      response.httpStatus = context.response?.status;
      response.path = context.url;
      // Most of the time when the Vault API returns an error, the response looks like:
      // { errors: ['some error message']}
      // But sometimes (eg RespondWithStatusCode) it looks like this:
      // { data: { error: 'some error message' } }
      if (response?.data?.error && !response.errors) {
        // Normalize the errors from RespondWithStatusCode
        response.errors = [response.data.error];
      }

      const blob = new Blob([JSON.stringify(response, null, 2)], { type: 'application/json' });
      const { headers, status, statusText } = context.response || {};
      return new Response(blob, { headers, status, statusText });
    }

    return context.response;
  };

  configuration = new Configuration({
    basePath: '/v1',
    middleware: [{ pre: this.preRequest }, { post: this.postRequest }],
  });

  auth = new AuthApi(this.configuration);
  identity = new IdentityApi(this.configuration);
  secrets = new SecretsApi(this.configuration);
  sys = new SystemApi(this.configuration);
}
