import FetchConfigRoute from '../../fetch-config';
import { hash } from 'ember-concurrency';
export default class KubernetesRoleCredentialsRoute extends FetchConfigRoute {
  model() {
    const { name } = this.paramsFor('roles.role');

    return hash({
      roleName: name,
      kubernetesBackend: this.secretMountPath.get(),
      backend: this.modelFor('application'),
    });
  }
}
