import Route from '@ember/routing/route';

export default Route.extend({
  templateName: 'vault/cluster/oidc-callback',
  model() {
    // left blank so we render the template immediately
  },
  // const queryString = decodeURIComponent(window.location.search);
  // Since state param can also contain namespace, fetch the values using native url api.
  // For instance, state params value can be state=st_123456,ns=d4fq
  afterModel() {
    let { auth_path: path, code, state } = this.paramsFor(this.routeName);
    let { namespaceQueryParam: namespace } = this.paramsFor('vault.cluster');
    console.log(namespace, 'ANAMESPCE');
    if (namespace === '' && state?.includes(',ns=')) {
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
