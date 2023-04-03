/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Engine from '@ember/engine';
import loadInitializers from 'ember-load-initializers';
import Resolver from 'ember-resolver';
import config from './config/environment';

const { modulePrefix } = config;

export default class ReplicationEngine extends Engine {
  modulePrefix = modulePrefix;
  Resolver = Resolver;
  dependencies = {
    services: ['auth', 'namespace', 'replication-mode', 'router', 'store', 'version'],
    externalRoutes: ['replication'],
  };
}

loadInitializers(ReplicationEngine, modulePrefix);
