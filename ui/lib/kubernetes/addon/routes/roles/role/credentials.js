import Route from '@ember/routing/route';
export default class KubernetesRoleCredentialsRoute extends Route {
  model() {
    return this.modelFor('roles.role');
  }
}
