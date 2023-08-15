/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import { inject as service } from '@ember/service';
import parseURL from 'core/utils/parse-url';
import config from 'open-api-explorer/config/environment';

const { APP } = config;

const SearchFilterPlugin = () => {
  return {
    fn: {
      opsFilter: (taggedOps, phrase) => {
        // map over the options and filter out operations where the path doesn't match what's typed
        return (
          taggedOps
            .map((tagObj) => {
              const operations = tagObj.get('operations').filter((operationObj) => {
                return operationObj.get('path').includes(phrase);
              });
              return tagObj.set('operations', operations);
            })
            // then traverse again and remove the top level item if there are no operations left after filtering
            .filter((tagObj) => !!tagObj.get('operations').size)
        );
      },
    },
  };
};

const CONFIG = (SwaggerUIBundle, componentInstance, initialFilter) => {
  return {
    dom_id: `#${componentInstance.elementId}-swagger`,
    url: '/v1/sys/internal/specs/openapi',
    deepLinking: false,
    presets: [SwaggerUIBundle.presets.apis],
    plugins: [SwaggerUIBundle.plugins.DownloadUrl, SearchFilterPlugin],
    // 'list' expands tags, but not operations
    docExpansion: 'list',
    operationsSorter: 'alpha',
    filter: initialFilter || true,
    // this makes sure we show the x-vault- options
    showExtensions: true,
    // we don't have any models defined currently
    defaultModelsExpandDepth: -1,
    defaultModelExpandDepth: 1,
    requestInterceptor: (req) => {
      // we need to add vault authorization header
      // and namepace headers for things to work properly
      req.headers['X-Vault-Token'] = componentInstance.auth.currentToken;

      const namespace = componentInstance.namespaceService.path;
      if (namespace && !APP.NAMESPACE_ROOT_URLS.some((str) => req.url.includes(str))) {
        req.headers['X-Vault-Namespace'] = namespace;
      }
      // we want to link to the right JSON in swagger UI so
      // it's already been pre-pended
      if (!req.loadSpec) {
        const { protocol, host, pathname, search } = parseURL(req.url);
        //paths in the spec don't have /v1 in them, so we need to add that here
        //           http(s):  vlt.io:4200  /sys/mounts
        req.url = `${protocol}//${host}/v1${pathname}${search}`;
      }
      return req;
    },
    onComplete: () => {
      componentInstance.set('swaggerLoading', false);
    },
  };
};

export default Component.extend({
  auth: service(),
  namespaceService: service('namespace'),
  initialFilter: null,
  onFilterChange() {},
  swaggerLoading: true,

  async didInsertElement() {
    const { default: SwaggerUIBundle } = await import('swagger-ui-dist/swagger-ui-bundle.js');
    this._super(...arguments);
    // trim any initial slashes
    const initialFilter = this.initialFilter.replace(/^(\/)+/, '');
    SwaggerUIBundle(CONFIG(SwaggerUIBundle, this, initialFilter));
  },

  actions: {
    // sets the filter so the query param is updated so we get sharable URLs
    updateFilter(e) {
      this.onFilterChange(e.target.value || '');
    },
    proxyEvent(e) {
      const swaggerInput = this.element.querySelector('.operation-filter-input');
      // if this breaks because of a react upgrade,
      // change this to
      //let originalSetter = Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, 'value').set;
      //originalSetter.call(swaggerInput, e.target.value);
      // see post on triggering react events externally for an explanation of
      // why this works: https://stackoverflow.com/a/46012210
      const evt = new Event('input', { bubbles: true });
      evt.simulated = true;
      swaggerInput.value = e.target.value.replace(/^(\/)+/, '');
      swaggerInput.dispatchEvent(evt);
    },
  },
});
