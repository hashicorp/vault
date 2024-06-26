/**
 * Copyright (c) HashiCorp, Inc.
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
    services: ['auth', 'store', 'flash-messages', 'namespace', 'router', 'version', 'custom-messages'],
  };
}

loadInitializers(ConfigUiEngine, modulePrefix);
