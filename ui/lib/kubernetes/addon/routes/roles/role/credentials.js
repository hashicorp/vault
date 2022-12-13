import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
export default class KubernetesRoleCredentialsRoute extends Route {
  @service secretMountPath;

  model() {
    return {
      roleName: this.paramsFor('roles.role').name,
      backend: this.secretMountPath.get(),
    };
  }
}
