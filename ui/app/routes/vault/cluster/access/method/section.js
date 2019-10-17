import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import { allSettled } from 'rsvp';

export default Route.extend({
  wizard: service(),
  pathHelp: service('path-help'),

  configModels(backendType) {
    const backends = ['azure', 'github', 'gcp', 'jwt', 'oidc', 'kubernetes', 'ldap', 'okta', 'radius'];
    if (backendType === 'aws') {
      return [
        'auth-config/aws/client',
        'auth-config/aws/identity-whitelist',
        'auth-config/aws/roletag-blacklist',
      ];
    }
    if (backends.includes(backendType)) {
      return [`auth-config/${backendType}`];
    }
    return [];
  },

  beforeModel() {
    const { apiPath, method, type } = this.getMethodAndModelInfo();
    let configModelTypes = this.configModels(type);
    let configModels = configModelTypes.map(config => {
      return this.pathHelp.getNewModel(config, method, apiPath);
    });
    return allSettled(configModels);
  },

  getMethodAndModelInfo() {
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    const methodModel = this.modelFor('vault.cluster.access.method');
    const { apiPath, type } = methodModel;
    return { apiPath, type, method };
  },

  model() {
    const backend = this.modelFor('vault.cluster.access.method');
    this.wizard.transitionFeatureMachine(this.wizard.featureState, 'DETAILS', backend.type);
    let configModelTypes = this.configModels(backend.type);
    let authConfigs = configModelTypes.map(config => this.store.findRecord(config, backend.id));

    return allSettled(authConfigs).then(configs => {
      backend.authConfigs.pushObjects(configs.filter(config => config.state === 'fulfilled').mapBy('value'));
      return backend;
    });
  },

  setupController(controller) {
    const { section_name: section } = this.paramsFor(this.routeName);
    this._super(...arguments);
    controller.set('section', section);
    let method = this.modelFor('vault.cluster.access.method');
    let paths = method.paths.navPaths.map(pathInfo => pathInfo.path);
    controller.set('paths', paths);
  },
});
