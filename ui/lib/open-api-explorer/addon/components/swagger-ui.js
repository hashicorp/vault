import Component from '@ember/component';
import { inject as service } from '@ember/service';
import parseURL from 'core/utils/parse-url';
import config from 'open-api-explorer/config/environment';
import Swag from 'swagger-ui-dist';

const { SwaggerUIBundle } = Swag;
const { APP } = config;

const CONFIG = (componentInstance, initialFilter) => {
  return {
    dom_id: `#${componentInstance.elementId}-swagger`,
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
      if (namespace && !APP.NAMESPACE_ROOT_URLS.some(str => req.url.includes(str))) {
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
  auth: service(),
  namespaceService: service('namespace'),
  initialFilter: null,
  onFilterChange() {},

  init() {
    this._super(...arguments);
    // the react app (SwaggerUI) is calling the opsFilter function - rebinding here lets us have a reference
    // to the component's updateFilter so that we can track the react app's state out to the ember app's url 🙃
    this.searchFilterPlugin = this.searchFilterPlugin.bind(this);
  },

  searchFilterPlugin() {
    return {
      fn: {
        // apparently this doesn't fire if `phrase` is empty so we can't zero out our query param :-/
        opsFilter: (taggedOps, phrase) => {
          // we don't want the initial slash in the query param, so call the component fn first
          this.updatedFilter(phrase);
          // but we do want it for
          let path = '/' + phrase;
          // map over the options and filter out operations where the path doesn't match what's typed
          return (
            taggedOps
              .map((tagObj, tag) => {
                let operations = tagObj.get('operations').filter(operationObj => {
                  // TODO: should this be includes instead of startsWith? I'm thinking yes
                  return operationObj.get('path').startsWith(path);
                });
                return tagObj.set('operations', operations);
              })
              // then traverse again and remove the top level item if there are no operations left after filtering
              .filter(tagObj => !!tagObj.get('operations').size)
          );
        },
      },
    };
  },

  didInsertElement() {
    this._super(...arguments);
    // trim any initial slashes
    let initialFilter = this.initialFilter.replace(/^(\/)+/, '');
    SwaggerUIBundle(CONFIG(this, initialFilter));
  },

  // sets the filter so the query param is updated so we get sharable URLs
  updatedFilter(val) {
    this.onFilterChange(val || '');
  },
});
