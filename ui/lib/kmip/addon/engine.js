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
      'flash-messages',
      'namespace',
      'path-help',
      'router',
      'store',
      'version',
      'wizard',
      'secret-mount-path',
    ],
    externalRoutes: ['secrets'],
  },
});

loadInitializers(Eng, modulePrefix);

export default Eng;
