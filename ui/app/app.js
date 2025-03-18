/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Application from '@ember/application';
import Resolver from 'ember-resolver';
import loadInitializers from 'ember-load-initializers';
import config from 'vault/config/environment';
import { ROUTES } from 'vault/utils/routes';

export default class App extends Application {
  modulePrefix = config.modulePrefix;
  podModulePrefix = config.podModulePrefix;
  Resolver = Resolver;
  engines = {
    'config-ui': {
      dependencies: {
        services: [
          'auth',
          'flash-messages',
          'namespace',
          { 'app-router': 'router' },
          'store',
          'pagination',
          'version',
          'custom-messages',
        ],
      },
    },
    'open-api-explorer': {
      dependencies: {
        services: ['auth', 'flash-messages', 'namespace', { 'app-router': 'router' }, 'version'],
      },
    },
    replication: {
      dependencies: {
        services: [
          'auth',
          'capabilities',
          'flash-messages',
          'namespace',
          'replication-mode',
          { 'app-router': 'router' },
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
          { 'app-router': 'router' },
          'store',
          'pagination',
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
        services: [{ 'app-router': 'router' }, 'store', 'secret-mount-path', 'flash-messages'],
        externalRoutes: {
          secrets: 'vault.cluster.secrets.backends',
        },
      },
    },
    ldap: {
      dependencies: {
        services: [
          { 'app-router': 'router' },
          'store',
          'pagination',
          'secret-mount-path',
          'flash-messages',
          'auth',
        ],
        externalRoutes: {
          secrets: 'vault.cluster.secrets.backends',
        },
      },
    },
    kv: {
      dependencies: {
        services: [
          'capabilities',
          'control-group',
          'download',
          'flash-messages',
          'namespace',
          { 'app-router': 'router' },
          'secret-mount-path',
          'store',
          'pagination',
          'version',
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
          { 'app-router': 'router' },
          'secret-mount-path',
          'store',
          'pagination',
          'version',
        ],
        externalRoutes: {
          secrets: 'vault.cluster.secrets.backends',
          externalMountIssuer: ROUTES.VAULT_CLUSTER_SECRETS_BACKEND_PKI_ISSUERS_ISSUER_DETAILS,
          secretsListRootConfiguration: 'vault.cluster.secrets.backend.configuration',
        },
      },
    },
    sync: {
      dependencies: {
        services: ['flash-messages', 'flags', { 'app-router': 'router' }, 'store', 'pagination', 'version'],
        externalRoutes: {
          kvSecretOverview: ROUTES.VAULT_CLUSTER_SECRETS_BACKEND_KV_SECRET_INDEX,
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
