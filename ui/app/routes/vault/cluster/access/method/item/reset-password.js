import Route from '@ember/routing/route';

export default class VaultClusterAccessMethodItemResetPasswordRoute extends Route {
  beforeModel() {
    // TODO: redirect if not allowed to reset?
  }

  fetchPasswordPolicies() {}

  model() {
    return {
      backend: 'userpass',
      username: 'bob',
      passwordPolicies: ['example'],
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'methods', route: 'vault.cluster.access' },
      { label: 'users', route: 'vault.cluster.access.method' },
      { label: 'reset password' },
    ];
  }
}
