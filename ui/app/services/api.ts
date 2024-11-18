import Service, { service } from '@ember/service';
import createClient from 'openapi-fetch';
import config from 'vault/config/environment';

import type Owner from '@ember/owner';
import type { paths } from 'vault/vault-openapi-schema';
import type { Middleware } from 'openapi-fetch';
import type AuthService from 'vault/services/auth';
import type NamespaceService from 'vault/services/namespace';
import type ControlGroupService from 'vault/services/control-group';

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
      const { path: namespace } = this.namespace;
      const { currentToken } = this.auth;
      const rootUrls = config.APP['NAMESPACE_ROOT_URLS'] as string[];
      const pollingUrls = config.APP['NAMESPACE_ROOT_URLS'] as string[];

      if (currentToken) {
        request.headers.set('X-Vault-Token', currentToken);
      }
      if (request.method === 'PATCH') {
        request.headers.set('Content-Type', 'application/merge-patch+json');
      }
      if (namespace && !rootUrls.some((str) => request.url.includes(str))) {
        request.headers.set('X-Vault-Namespace', namespace);
      }
      if (!pollingUrls.some((str) => request.url.includes(str))) {
        this.auth.setLastFetch(Date.now());
      }

      return request;
    },
    onResponse: async ({ response }) => {
      const { status, ...resOptions } = response;
      try {
        const json = await response.clone().json();

        if (status >= 200 && status < 400) {
          // if a data key is present in the response, return it directly
          if (json.data) {
            return new Response(JSON.stringify(json.data), { status, ...resOptions });
          }
        } else {
          /**
           * setting response properties on the body is unnecessary since the response is returned from the client
           * shimming this for now to avoid breaking changes while moving away from Ember Data
           */
          // ember data errors don't have the status code, so we add it here
          json.httpStatus = response.status;
          json.path = response.url;
          // Most of the time when the Vault API returns an error, the payload looks like:
          // { errors: ['some error message']}
          // But sometimes (eg RespondWithStatusCode) it looks like this:
          // { data: { error: 'some error message' } }
          if (json?.data?.error && !json.errors) {
            // Normalize the errors from RespondWithStatusCode
            json.errors = [json.data.error];
          }

          return new Response(JSON.stringify(json), { status, ...resOptions });
        }
      } catch (e) {
        // ignore errors from parsing an empty response body
      }

      // skip and return the original response
      return undefined;
    },
  };
}
