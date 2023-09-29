/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ENV from 'vault/config/environment';
import handlers from './handlers';

// remember to export handler name from mirage/handlers/index.js file

export default function () {
  this.namespace = 'v1';

  // start ember in development running mirage -> yarn start:mirage handlerName
  // if handler is not provided, general config will be used
  // this is useful for feature development when a specific and limited config is required
  const { handler } = ENV['ember-cli-mirage'];
  const handlerName = handler in handlers ? handler : 'base';
  handlers[handlerName](this);
  this.logging = false; // disables passthrough logging which spams the console
  console.log(`⚙ Using ${handlerName} Mirage request handlers ⚙`); // eslint-disable-line
  // passthrough all unhandled requests
  this.passthrough();
}
