import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class KubernetesFetchConfigRoute extends Route {
  @service store;
  @service secretMountPath;

  configModel = null;

  async beforeModel() {
    const backend = this.secretMountPath.get();
    // check the store for record first
    this.configModel = this.store.peekRecord('kubernetes/config', backend);
    if (!this.configModel) {
      return this.store
        .queryRecord('kubernetes/config', { backend })
        .then((record) => {
          this.configModel = record;
        })
        .catch(() => {
          // it's ok! we don't need to transition to the error route
        });
    }
  }
}
