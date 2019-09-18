import Application from '@ember/application';
import Resolver from './resolver';
import loadInitializers from 'ember-load-initializers';
import config from './config/environment';
import defineModifier from 'ember-concurrency-test-waiter/define-modifier';

defineModifier();

let App;

/* eslint-disable ember/avoid-leaking-state-in-ember-objects */
App = Application.extend({
  modulePrefix: config.modulePrefix,
  podModulePrefix: config.podModulePrefix,
  Resolver,
  engines: {
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
  },
});

loadInitializers(App, config.modulePrefix);

export default App;
