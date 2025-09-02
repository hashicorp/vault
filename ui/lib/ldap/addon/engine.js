/**
 * Copyright (c) HashiCorp, Inc.
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
    services: ['app-router', 'store', 'pagination', 'secret-mount-path', 'flash-messages', 'auth'],
    externalRoutes: ['secrets'],
  };
}

loadInitializers(LdapEngine, modulePrefix);
