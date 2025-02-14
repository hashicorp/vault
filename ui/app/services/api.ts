import Service, { service } from '@ember/service';
import createClient from 'openapi-fetch';
import config from 'vault/config/environment';

import type Owner from '@ember/owner';
import type { paths } from 'vault/api';
import type { Middleware } from 'openapi-fetch';
import type AuthService from 'vault/services/auth';
import type NamespaceService from 'vault/services/namespace';
import type ControlGroupService from 'vault/services/control-group';

const {
  APP: { POLLING_URLS, NAMESPACE_ROOT_URLS },
} = config;

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
  @service declare readonly auth: AuthService;
  @service declare readonly namespace: NamespaceService;
  @service declare readonly controlGroup: ControlGroupService;

  constructor(owner: Owner) {
    super(owner);
    this.client.use(this.middleware);
  }

  private client = createClient<paths>({
    baseUrl: '/v1',
  });

  get = this.client.GET;
  post = this.client.POST;
  put = this.client.PUT;
  delete = this.client.DELETE;
  patch = this.client.PATCH;

  middleware: Middleware = {
    onRequest: async ({ request }) => {
      const isPolling = POLLING_URLS.some((str) => request.url.includes(str));
      if (!isPolling) {
        this.auth.setLastFetch(Date.now());
      }

      // since the client methods don't accept random options, perhaps these can be set on the authService instead
      const token = /* options.clientToken || */ this.auth.currentToken;
      if (token /* && !options.unauthenticated */) {
        request.headers.set('X-Vault-Token', token);
      }
      // if (options.wrapTTL) {
      //   headers['X-Vault-Wrap-TTL'] = options.wrapTTL;
      // }
      if (request.method === 'PATCH') {
        request.headers.set('Content-Type', 'application/merge-patch+json');
      }
      // similarly to clientToken, perhaps this can always be set on the namespaceService?
      // const namespace =
      //   typeof options.namespace === 'undefined' ? this.namespaceService.path : options.namespace;
      const namespace = this.namespace.path;
      if (namespace && !NAMESPACE_ROOT_URLS.some((str) => request.url.includes(str))) {
        request.headers.set('X-Vault-Namespace', namespace);
      }

      return request;
    },
    onResponse: async ({ response }) => {
      const { status, headers, ...resOptions } = response;

      if (!headers.get('Content-Length')) {
        return undefined;
      }

      const json = await response.clone().json();

      if (status >= 400) {
        // backwards compatibility with Ember Data
        const error: ApiError = json || {};
        error.httpStatus = response.status;
        error.path = response.url;
        // Most of the time when the Vault API returns an error, the response looks like:
        // { errors: ['some error message']}
        // But sometimes (eg RespondWithStatusCode) it looks like this:
        // { data: { error: 'some error message' } }
        if (error?.data?.error && !error.errors) {
          // Normalize the errors from RespondWithStatusCode
          error.errors = [error.data.error];
        }

        return new Response(JSON.stringify(error), { status, headers, ...resOptions });
      } else {
        // the responses in the OpenAPI spec don't account for the return values to be under the 'data' key
        // extract the data from the response so that it is returned by the client
        if (json.data) {
          return new Response(JSON.stringify(json.data), { status, ...resOptions });
        }
      }

      return undefined;
    },
  };
}
