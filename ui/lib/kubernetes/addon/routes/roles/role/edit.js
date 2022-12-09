import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class KubernetesRoleEditRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const backend = this.secretMountPath.get();
    const { name } = this.paramsFor('roles.role');
    return this.store.queryRecord('kubernetes/role', { backend, name });
  }
}
