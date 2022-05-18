import Route from '@ember/routing/route';

export default Route.extend({
  templateName: 'vault/cluster/oidc-callback',
  model() {
    // left blank so we render the template immediately
  },
  afterModel() {
    const queryString = decodeURIComponent(window.location.search);
    // Since state param can also contain namespace, fetch the values using native url api.
    // For instance, state params value can be state=st_123456,ns=d4fq
    // Ember paramsFor used to strip out the value after the "=" sign. In short ns value was not being passed along.
    let urlParams = new URLSearchParams(queryString);
    let state = urlParams.get('state'),
      code = urlParams.get('code'),
      ns;
    if (state.includes(',ns=')) {
      let arrayParams = state.split(',ns=');
      state = arrayParams[0];
      ns = arrayParams[1];
    }
    let { auth_path: path } = this.paramsFor(this.routeName);
    let { namespaceQueryParam: namespace } = this.paramsFor('vault.cluster');
    path = window.decodeURIComponent(path);
    const source = 'oidc-callback'; // required by event listener in auth-jwt component
    let queryParams = { source, namespace, path, code, state };
    // If state had ns value, send it as part of namespace param
    if (ns) {
      queryParams.namespace = ns;
    }
    window.opener.postMessage(queryParams, window.origin);
  },
  setupController(controller) {
    this._super(...arguments);
    controller.set('pageContainer', document.querySelector('.page-container'));
  },
});
