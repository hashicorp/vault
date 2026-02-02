/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Engine from '@ember/engine';

import loadInitializers from 'ember-load-initializers';
import Resolver from 'ember-resolver';

import config from './config/environment';

const { modulePrefix } = config;

export default class ConfigUiEngine extends Engine {
  modulePrefix = modulePrefix;
  Resolver = Resolver;
  dependencies = {
    services: [
      'auth',
      'flash-messages',
      'namespace',
      'app-router',
      'version',
      'custom-messages',
      'api',
      'capabilities',
      // services needed for tools sidebar component
      'permissions',
      'current-cluster',
      '-portal',
    ],
    // 'vault', 'tool', 'messages', 'openApiExplorer', 'loginSettings' external routes are used in tools sidebar component
    externalRoutes: ['vault', 'tool', 'messages', 'openApiExplorer', 'loginSettings'],
  };
}

loadInitializers(ConfigUiEngine, modulePrefix);
