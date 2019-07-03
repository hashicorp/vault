import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';

export default Route.extend({
  wizard: service(),
  pathHelp: service('path-help'),

  modelType(backendType, section) {
    const MODELS = {
      'aws-client': 'auth-config/aws/client',
      'aws-identity-whitelist': 'auth-config/aws/identity-whitelist',
      'aws-roletag-blacklist': 'auth-config/aws/roletag-blacklist',
      'azure-configuration': 'auth-config/azure',
      'github-configuration': 'auth-config/github',
      'gcp-configuration': 'auth-config/gcp',
      'jwt-configuration': 'auth-config/jwt',
      'oidc-configuration': 'auth-config/oidc',
      'kubernetes-configuration': 'auth-config/kubernetes',
      'ldap-configuration': 'auth-config/ldap',
      'okta-configuration': 'auth-config/okta',
      'radius-configuration': 'auth-config/radius',
    };
    return MODELS[`${backendType}-${section}`];
  },

  beforeModel() {
    const { apiPath, method, type } = this.getMethodAndModelInfo();
    let modelType = this.modelType(type, 'configuration');
    if (modelType) {
      return this.pathHelp.getNewModel(modelType, method, apiPath);
    }
  },

  getMethodAndModelInfo() {
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    const methodModel = this.modelFor('vault.cluster.access.method');
    const { apiPath, type } = methodModel;
    return { apiPath, type, method };
  },

  model() {
    const backend = this.modelFor('vault.cluster.access.method');
    const modelType = this.modelType(backend.type, 'configuration');
    this.wizard.transitionFeatureMachine(this.wizard.featureState, 'DETAILS', backend.type);

    if (!modelType) {
      return backend; //tune options
    }

    return this.store
      .findRecord(modelType, backend.id)
      .then(methodConfig => {
        methodConfig.set('backend', backend);
        return methodConfig;
      })
      .catch(e => {
        // if you haven't saved a config, the API 404s
        // we still have tune options
        if (e.httpStatus === 404) {
          return backend;
        }
        throw e;
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
