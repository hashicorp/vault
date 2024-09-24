/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import parseURL from 'core/utils/parse-url';
import config from 'open-api-explorer/config/environment';
import { guidFor } from '@ember/object/internals';
import SwaggerUIBundle from 'swagger-ui-dist/swagger-ui-bundle.js';

const { APP } = config;

export default class SwaggerUiComponent extends Component {
  @service auth;
  @service namespace;

  @tracked swaggerLoading = true;

  inputId = `${guidFor(this)}-swagger`;

  SearchFilterPlugin() {
    return {
      fn: {
        opsFilter: (taggedOps, phrase) => {
          const filteredOperations = taggedOps.reduce((acc, tagObj) => {
            const operations = tagObj.get('operations');

            // filter out operations where the path doesn't match search phrase
            const operationsWithMatchingPath = operations.filter((operationObj) => {
              const path = operationObj.get('path');
              return path.includes(phrase);
            });

            // if there are any operations left after filtering, add the tagObj to the accumulator
            if (operationsWithMatchingPath.size > 0) {
              acc.push(tagObj.set('operations', operationsWithMatchingPath));
            }

            return acc;
          }, []);

          return filteredOperations;
        },
      },
    };
  }

  CONFIG = (SwaggerUIBundle, componentInstance) => {
    return {
      dom_id: `#${componentInstance.inputId}`,
      url: '/v1/sys/internal/specs/openapi',
      deepLinking: false,
      presets: [SwaggerUIBundle.presets.apis],
      plugins: [SwaggerUIBundle.plugins.DownloadUrl, this.SearchFilterPlugin],
      // 'list' expands tags, but not operations
      docExpansion: 'list',
      operationsSorter: 'alpha',
      filter: true,
      // this makes sure we show the x-vault- options
      showExtensions: true,
      // we don't have any models defined currently
      defaultModelsExpandDepth: -1,
      defaultModelExpandDepth: 1,
      requestInterceptor: (req) => {
        // we need to add vault authorization header
        // and namespace headers for things to work properly
        req.headers['X-Vault-Token'] = componentInstance.auth.currentToken;
        const namespace = componentInstance.namespace.path;
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
        componentInstance.swaggerLoading = false;
      },
    };
  };

  // using an action to bind the correct "this" context
  @action async swaggerInit() {
    const configSettings = this.CONFIG(SwaggerUIBundle, this);
    SwaggerUIBundle(configSettings);
  }
}
