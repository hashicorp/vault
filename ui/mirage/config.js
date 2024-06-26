/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { createServer } from 'miragejs';
import { discoverEmberDataModels } from 'ember-cli-mirage';
import ENV from 'vault/config/environment';
import handlers from './handlers';
// remember to export handler name from mirage/handlers/index.js file

export default function (config) {
  const finalConfig = {
    ...config,
    logging: false,
    models: {
      ...discoverEmberDataModels(config.store),
      ...config.models,
    },
    routes,
  };

  return createServer(finalConfig);
}

function routes() {
  this.namespace = 'v1';

  const { handler } = ENV['ember-cli-mirage'];
  const handlerName = handler in handlers ? handler : 'base';
  handlers[handlerName](this);
  this.logging = false; // disables passthrough logging which spams the console
  console.log(`⚙ Using ${handlerName} Mirage request handlers ⚙`); // eslint-disable-line

  this.passthrough();
}
