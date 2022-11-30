import FetchConfigRoute from './fetch-config';

export default class KubernetesConfigureRoute extends FetchConfigRoute {
  model() {
    return {
      backend: this.modelFor('application'),
      config: this.configModel,
    };
  }
}
