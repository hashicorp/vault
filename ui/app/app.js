/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Application from '@ember/application';
import Resolver from 'ember-resolver';
import loadInitializers from 'ember-load-initializers';
import config from 'vault/config/environment';

export default class App extends Application {
  modulePrefix = config.modulePrefix;
  podModulePrefix = config.podModulePrefix;
  Resolver = Resolver;
  engines = {
    'config-ui': {
      dependencies: {
        services: ['auth', 'flash-messages', 'namespace', 'router', 'store', 'version', 'custom-messages'],
      },
    },
    'open-api-explorer': {
      dependencies: {
        services: ['auth', 'flash-messages', 'namespace', 'router', 'version'],
      },
    },
    replication: {
      dependencies: {
        services: [
          'auth',
          'flash-messages',
          'namespace',
          'replication-mode',
          'router',
          'store',
          'version',
          '-portal',
        ],
        externalRoutes: {
          replication: 'vault.cluster.replication.index',
          vault: 'vault.cluster',
        },
      },
    },
    kmip: {
      dependencies: {
        services: [
          'auth',
          'download',
          'flash-messages',
          'namespace',
          'path-help',
          'router',
          'store',
          'version',
          'secret-mount-path',
        ],
        externalRoutes: {
          secrets: 'vault.cluster.secrets.backends',
        },
      },
    },
    kubernetes: {
      dependencies: {
        services: ['router', 'store', 'secret-mount-path', 'flash-messages'],
        externalRoutes: {
          secrets: 'vault.cluster.secrets.backends',
        },
      },
    },
    ldap: {
      dependencies: {
        services: ['router', 'store', 'secret-mount-path', 'flash-messages', 'auth'],
        externalRoutes: {
          secrets: 'vault.cluster.secrets.backends',
        },
      },
    },
    kv: {
      dependencies: {
        services: [
          'download',
          'namespace',
          'router',
          'store',
          'secret-mount-path',
          'flash-messages',
          'control-group',
        ],
        externalRoutes: {
          secrets: 'vault.cluster.secrets.backends',
          syncDestination: 'vault.cluster.sync.secrets.destinations.destination',
        },
      },
    },
    pki: {
      dependencies: {
        services: [
          'auth',
          'download',
          'flash-messages',
          'namespace',
          'path-help',
          'router',
          'secret-mount-path',
          'store',
          'version',
        ],
        externalRoutes: {
          secrets: 'vault.cluster.secrets.backends',
          externalMountIssuer: 'vault.cluster.secrets.backend.pki.issuers.issuer.details',
          secretsListRootConfiguration: 'vault.cluster.secrets.backend.configuration',
        },
      },
    },
    sync: {
      dependencies: {
        services: ['flash-messages', 'flags', 'router', 'store', 'version'],
        externalRoutes: {
          kvSecretDetails: 'vault.cluster.secrets.backend.kv.secret.details',
          clientCountOverview: 'vault.cluster.clients',
        },
      },
    },
  };
}

loadInitializers(App, config.modulePrefix);

/**
 * @typedef {import('ember-source/types')} EmberTypes
 */
