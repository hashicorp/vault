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
    replication: {
      dependencies: {
        services: ['auth', 'replication-mode', 'router', 'store', 'version', 'wizard'],
      },
    },
  },
});

loadInitializers(App, config.modulePrefix);

export default App;
