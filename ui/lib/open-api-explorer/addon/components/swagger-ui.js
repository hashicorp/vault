import Component from '@ember/component';
import { inject as service } from '@ember/service';
import parseURL from 'core/utils/parse-url';
import layout from '../templates/components/swagger-ui';

const CONFIG = (SwaggerUIBundle, componentInstance, initialFilter) => {
  return {
    dom_id: '#swagger-container',
    url: '/v1/sys/internal/specs/openapi',
    deepLinking: false,
    presets: [SwaggerUIBundle.presets.apis],
    plugins: [SwaggerUIBundle.plugins.DownloadUrl, componentInstance.searchFilterPlugin],
    // 'list' expands tags, but not operations
    docExpansion: 'list',
    operationsSorter: 'alpha',
    filter: initialFilter || true,
    // this makes sure we show the x-vault- options
    showExtensions: true,
    // we don't have any models defined currently
    defaultModelsExpandDepth: -1,
    defaultModelExpandDepth: 1,
    requestInterceptor: req => {
      // we need to add vault authorization header
      // and namepace headers for things to work properly
      req.headers['X-Vault-Token'] = componentInstance.auth.currentToken;

      let namespace = componentInstance.namespaceService.path;
      if (namespace && !NAMESPACE_ROOT_URLS.some(str => req.url.includes(str))) {
        req.headers['X-Vault-Namespace'] = namespace;
      }
      // we want to link to the right JSON in swagger UI so
      // it's already been pre-pended
      if (!req.loadSpec) {
        let { protocol, host, pathname } = parseURL(req.url);
        //paths in the spec don't have /v1 in them, so we need to add that here
        //           http(s):  vlt.io:4200  /sys/mounts
        req.url = `${protocol}//${host}/v1${pathname}`;
      }
      return req;
    },
  };
};

export default Component.extend({
  layout,
  tagName: '',
  auth: service(),
  namespaceService: service('namespace'),
  initialFilter: null,
  onFilterChange() {},

  // sets the filter so the query param is updated so we get sharable URLs
  updatedFilter(val) {
    this.onFilterChange(val || '');
  },

  init() {
    this._super(...arguments);
    // we need to rebind here because the react app is calling the opsFilter function - rebinding here lets us
    // have a reference to the component's updateFilter so that we can track the react app's state out to the
    // ember app's url 🙃
    this.searchFilterPlugin = this.searchFilterPlugin.bind(this);
  },

  searchFilterPlugin() {
    return {
      fn: {
        // apparently this doesn't fire if `phrase` is empty so we can't zero out our query param :-/
        opsFilter: (taggedOps, phrase) => {
          // we don't want the initial slash in the query param
          this.updatedFilter(phrase);
          // but we do want it for
          let path = '/' + phrase;
          return taggedOps
            .map((tagObj, tag) => {
              let operations = tagObj.get('operations').filter(operationObj => {
                return operationObj.get('path').startsWith(path);
              });
              return tagObj.set('operations', operations);
            })
            .filter(tagObj => !!tagObj.get('operations').size);
        },
      },
    };
  },

  async didInsertElement() {
    this._super(...arguments);
    // trim any initial slashes
    let initialFilter = this.initialFilter.replace(/^(\/)+/, '');
    let module = await import('swagger-ui-dist/swagger-ui-bundle.js');
    let SwaggerUIBundle = module.default;
    SwaggerUIBundle(CONFIG(SwaggerUIBundle, this, initialFilter));
  },
});
