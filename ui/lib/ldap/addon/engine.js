/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Engine from 'ember-engines/engine';
import loadInitializers from 'ember-load-initializers';
import Resolver from 'ember-resolver';
import config from './config/environment';

const { modulePrefix } = config;

export default class LdapEngine extends Engine {
  modulePrefix = modulePrefix;
  Resolver = Resolver;
  dependencies = {
    services: ['app-router', 'store', 'pagination', 'secret-mount-path', 'flash-messages', 'auth', 'api'],
    externalRoutes: ['secrets', 'secretsGeneralSettingsConfiguration'],
  };
}

loadInitializers(LdapEngine, modulePrefix);
