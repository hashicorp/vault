/**
 * Copyright (c) HashiCorp, Inc.
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
      'auth',
      'download',
      'flash-messages',
      'namespace',
      'path-help',
      'app-router',
      'store',
      'pagination',
      'version',
      'secret-mount-path',
    ],
    externalRoutes: ['secrets'],
  };
}

loadInitializers(KmipEngine, modulePrefix);
