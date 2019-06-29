import Component from '@ember/component';
import { inject as service } from '@ember/service';
import parseURL from 'core/utils/parse-url';
import layout from '../templates/components/swagger-ui';

export default Component.extend({
  layout,
  tagName: '',
  auth: service(),

  didInsertElement() {
    this._super(...arguments);
    import('swagger-ui-dist/swagger-ui-bundle.js').then(module => {
      let SwaggerUIBundle = module.default;
      const CONFIG = {
        dom_id: '#swagger-container',
        url: '/v1/sys/internal/specs/openapi',
        deepLinking: false,
        presets: [SwaggerUIBundle.presets.apis],
        plugins: [SwaggerUIBundle.plugins.DownloadUrl],
        docExpansion: 'none',
        operationsSorter: 'alpha',
        filter: true,
        showExtensions: true,
        defaultModelsExpandDepth: -1,
        defaultModelExpandDepth: 1,
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
      };

      SwaggerUIBundle(CONFIG);
    });
  },
});
