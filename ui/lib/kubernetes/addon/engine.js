import Engine from '@ember/engine';

import loadInitializers from 'ember-load-initializers';
import Resolver from 'ember-resolver';

import config from './config/environment';

const { modulePrefix } = config;

export default class KubernetesEngine extends Engine {
  modulePrefix = modulePrefix;
  Resolver = Resolver;
  dependencies = {
    services: ['router', 'store', 'secret-mount-path', 'flashMessages'],
    externalRoutes: ['secrets'],
  };
}

loadInitializers(KubernetesEngine, modulePrefix);
