import Route from '@ember/routing/route';
import Swag from 'swagger-ui';
import { inject as service } from '@ember/service';
import parseURL from 'core/utils/parse-url';

const { SwaggerUIBundle, SwaggerUIStandalonePreset } = Swag;

export default Route.extend({
  auth: service(),
  model() {
    return {
      url: '/v1/sys/internal/specs/openapi',
      requestInterceptor: req => {
        req.headers['X-Vault-Token'] = this.auth.currentToken;
        // we want to link to the right JSON in swagger UI so
        // it's already been pre-pended
        if (!req.loadSpec) {
          let { protocol, host, pathname } = parseURL(req.url);
          //           http(s):  vlt.io:4200  /sys/mounts
          req.url = `${protocol}//${host}/v1${pathname}`;
        }
        return req;
      },
      deepLinking: false,
      presets: [SwaggerUIStandalonePreset, SwaggerUIBundle.presets.apis, SwaggerUIBundle.plugins.DownloadUrl],
      layout: 'StandaloneLayout',
      docExpansion: 'none',
      tagsSorter: 'alpha',
      operationsSorter: 'alpha',
      defaultModelsExpandDepth: -1,
      defaultModelExpandDepth: 1,
    };
  },
});
