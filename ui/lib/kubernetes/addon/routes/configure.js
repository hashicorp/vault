import FetchConfigRoute from './fetch-config';

export default class KubernetesConfigureRoute extends FetchConfigRoute {
  async model() {
    const backend = this.secretMountPath.get();
    return this.configModel || this.store.createRecord('kubernetes/config', { backend });
  }
}
