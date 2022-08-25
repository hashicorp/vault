import Route from '@ember/routing/route';

export default Route.extend({
  templateName: 'vault/cluster/oidc-callback',
  model() {
    // left blank so we render the template immediately
  },
  afterModel() {
    let { auth_path: path, code, state } = this.paramsFor(this.routeName);
    let { namespaceQueryParam: namespace } = this.paramsFor('vault.cluster');
    // only replace namespace param from cluster if state has a namespace
    if (state?.includes(',ns=')) {
      let arrayParams = state.split(',ns=');
      state = arrayParams[0];
      namespace = arrayParams[1];
    }
    path = window.decodeURIComponent(path);
    const source = 'oidc-callback'; // required by event listener in auth-jwt component
    let queryParams = { source, namespace, path, code, state };
    window.opener.postMessage(queryParams, window.origin);
  },
  setupController(controller) {
    this._super(...arguments);
    controller.set('pageContainer', document.querySelector('.page-container'));
  },
});
