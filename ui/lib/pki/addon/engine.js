/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Engine from '@ember/engine';

import loadInitializers from 'ember-load-initializers';
import Resolver from 'ember-resolver';

import config from './config/environment';

const { modulePrefix } = config;

export default class PkiEngine extends Engine {
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
      'secret-mount-path',
      'version',
    ],
    externalRoutes: [
      'secrets',
      'secretsListRootConfiguration',
      'externalMountIssuer',
      'secretsGeneralSettingsConfiguration',
    ],
  };
}

loadInitializers(PkiEngine, modulePrefix);
