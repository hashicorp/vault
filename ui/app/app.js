import Application from '@ember/application';
import Resolver from 'ember-resolver';
import loadInitializers from 'ember-load-initializers';
import config from 'vault/config/environment';
import defineModifier from 'ember-concurrency-test-waiter/define-modifier';

defineModifier();
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
  };
}

loadInitializers(App, config.modulePrefix);
