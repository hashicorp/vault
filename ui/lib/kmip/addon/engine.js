/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Engine from 'ember-engines/engine';
import loadInitializers from 'ember-load-initializers';
import Resolver from 'ember-resolver';
import config from './config/environment';

const { modulePrefix } = config;
export default class KmipEngine extends Engine {
  modulePrefix = modulePrefix;
  Resolver = Resolver;
  dependencies = {
    services: [
      'api',
      'auth',
      'capabilities',
      'download',
      'flash-messages',
      'namespace',
      'path-help',
      'app-router',
      'version',
      'secret-mount-path',
    ],
    externalRoutes: ['secrets'],
  };
}

loadInitializers(KmipEngine, modulePrefix);
