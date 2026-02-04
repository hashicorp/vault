/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
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
    services: [
      'auth',
      'capabilities',
      'flash-messages',
      'namespace',
      'replication-mode',
      'app-router',
      'store',
      'version',
      // services needed for tools sidebar component
      'permissions',
      'current-cluster',
      'flags',
      '-portal',
      'control-group',
    ],
    externalRoutes: ['replication', 'vault', 'recoverySnapshots', 'settingsSeal', 'replicationMode'],
  },
});

loadInitializers(Eng, modulePrefix);

export default Eng;
