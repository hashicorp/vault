/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Engine from 'ember-engines/engine';
import loadInitializers from 'ember-load-initializers';
import Resolver from './resolver';
import config from './config/environment';

const { modulePrefix } = config;
/* eslint-disable ember/avoid-leaking-state-in-ember-objects */
const Eng = Engine.extend({
  modulePrefix,
  Resolver,
  dependencies: {
    services: ['auth', 'namespace', 'replication-mode', 'router', 'store', 'version'],
    externalRoutes: ['replication'],
  },
});

loadInitializers(Eng, modulePrefix);

export default Eng;
