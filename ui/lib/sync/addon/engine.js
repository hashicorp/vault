/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Engine from 'ember-engines/engine';
import loadInitializers from 'ember-load-initializers';
import Resolver from 'ember-resolver';
import config from './config/environment';

const { modulePrefix } = config;

export default class SyncEngine extends Engine {
  modulePrefix = modulePrefix;
  Resolver = Resolver;
  dependencies = {
    services: ['flash-messages', 'flags', 'app-router', 'store', 'version'],
    externalRoutes: ['kvSecretOverview', 'clientCountOverview'],
  };
}

loadInitializers(SyncEngine, modulePrefix);
