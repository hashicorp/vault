import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
export default class KubernetesRoleCredentialsRoute extends Route {
  @service secretMountPath;

  model() {
    return {
      roleModel: this.modelFor('roles.role'),
      backend: this.secretMountPath.get(),
    };
  }
}
