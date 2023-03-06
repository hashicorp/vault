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
    // Ember paramsFor will strip out the value after the "=" sign.
    const urlParams = new URLSearchParams(queryString);
    let state = urlParams.get('state');
    const code = urlParams.get('code');

    let { auth_path: path } = this.paramsFor(this.routeName);
    let { namespaceQueryParam: namespace } = this.paramsFor('vault.cluster');
    // only replace namespace param from cluster if state has a namespace
    if (state?.includes(',ns=')) {
      [state, namespace] = state.split(',ns=');
    }
    path = window.decodeURIComponent(path);
    const source = 'oidc-callback'; // required by event listener in auth-jwt component
    const queryParams = { source, path: path || '', code: code || '', state: state || '' };
    if (namespace) {
      queryParams.namespace = namespace;
    }
    window.opener.postMessage(queryParams, window.origin);
  },
  setupController(controller) {
    this._super(...arguments);
    controller.set('pageContainer', document.querySelector('.page-container'));
  },
});
