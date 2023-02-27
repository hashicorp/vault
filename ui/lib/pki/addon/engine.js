import Engine from '@ember/engine';

import loadInitializers from 'ember-load-initializers';
import Resolver from 'ember-resolver';

import config from './config/environment';

const { modulePrefix } = config;

export default class PkiEngine extends Engine {
  modulePrefix = modulePrefix;
  Resolver = Resolver;
  dependencies = {
    services: [
      'auth',
      'flash-messages',
      'namespace',
      'path-help',
      'router',
      'store',
      'version',
      'secret-mount-path',
    ],
    externalRoutes: ['secrets'],
  };
}

loadInitializers(PkiEngine, modulePrefix);
