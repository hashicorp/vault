/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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
    openApiExplorer: {
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
          'wizard',
        ],
        externalRoutes: {
          replication: 'vault.cluster.replication.index',
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
          'wizard',
          'secret-mount-path',
        ],
        externalRoutes: {
          secrets: 'vault.cluster.secrets.backends',
        },
      },
    },
    kubernetes: {
      dependencies: {
        services: ['router', 'store', 'secret-mount-path', 'flashMessages'],
        externalRoutes: {
          secrets: 'vault.cluster.secrets.backends',
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
          'wizard',
        ],
        externalRoutes: {
          secrets: 'vault.cluster.secrets.backends',
          externalMountIssuer: 'vault.cluster.secrets.backend.pki.issuers.issuer.details',
          secretsListRoot: 'vault.cluster.secrets.backend.list-root',
          secretsListRootConfiguration: 'vault.cluster.secrets.backend.configuration',
        },
      },
    },
  };
}

loadInitializers(App, config.modulePrefix);
