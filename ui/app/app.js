/**
 * Copyright IBM Corp. 2016, 2025
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
        services: [
          'auth',
          'flash-messages',
          'namespace',
          { 'app-router': 'router' },
          'version',
          'custom-messages',
          'api',
          'capabilities',
          // services needed for tools sidebar component
          'permissions',
          'current-cluster',
          '-portal',
        ],
        externalRoutes: {
          vault: 'vault.cluster',
          tool: 'vault.cluster.tools.tool',
          messages: 'vault.cluster.config-ui.messages',
          openApiExplorer: 'vault.cluster.tools.open-api-explorer',
          loginSettings: 'vault.cluster.config-ui.login-settings',
        },
      },
    },
    'open-api-explorer': {
      dependencies: {
        services: ['auth', 'flash-messages', 'namespace', { 'app-router': 'router' }, 'version'],
        externalRoutes: { vault: 'vault.cluster' },
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
          'api',
          'auth',
          'capabilities',
          'download',
          'flash-messages',
          'namespace',
          'path-help',
          { 'app-router': 'router' },
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
        services: [{ 'app-router': 'router' }, 'secret-mount-path', 'flash-messages', 'api', 'capabilities'],
        externalRoutes: {
          secrets: 'vault.cluster.secrets.backends',
          secretsGeneralSettingsConfiguration: 'vault.cluster.secrets.backend.configuration.general-settings',
        },
      },
    },
    ldap: {
      dependencies: {
        services: [
          { 'app-router': 'router' },
          'secret-mount-path',
          'flash-messages',
          'auth',
          'api',
          'capabilities',
        ],
        externalRoutes: {
          secrets: 'vault.cluster.secrets.backends',
          secretsGeneralSettingsConfiguration: 'vault.cluster.secrets.backend.configuration.general-settings',
          secretsPluginSettingsConfiguration: 'vault.cluster.secrets.backend.configuration.plugin-settings',
        },
      },
    },
    kv: {
      dependencies: {
        services: [
          'api',
          'capabilities',
          'control-group',
          'download',
          'flash-messages',
          'namespace',
          { 'app-router': 'router' },
          'secret-mount-path',
          'version',
        ],
        externalRoutes: {
          secrets: 'vault.cluster.secrets.backends',
          syncDestination: 'vault.cluster.sync.secrets.destinations.destination',
          secretsGeneralSettingsConfiguration: 'vault.cluster.secrets.backend.configuration.general-settings',
        },
      },
    },
    pki: {
      dependencies: {
        services: [
          'api',
          'auth',
          'capabilities',
          'download',
          'flash-messages',
          'namespace',
          'path-help',
          { 'app-router': 'router' },
          'secret-mount-path',
          'version',
        ],
        externalRoutes: {
          secrets: 'vault.cluster.secrets.backends',
          externalMountIssuer: 'vault.cluster.secrets.backend.pki.issuers.issuer.details',
          secretsListRootConfiguration: 'vault.cluster.secrets.backend.configuration',
          secretsGeneralSettingsConfiguration: 'vault.cluster.secrets.backend.configuration.general-settings',
        },
      },
    },
    sync: {
      dependencies: {
        services: [
          'flash-messages',
          'flags',
          { 'app-router': 'router' },
          'store',
          'api',
          'capabilities',
          'version',
        ],
        externalRoutes: {
          kvSecretOverview: 'vault.cluster.secrets.backend.kv.secret.index',
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
