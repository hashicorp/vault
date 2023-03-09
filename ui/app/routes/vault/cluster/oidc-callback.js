import Route from '@ember/routing/route';

export default Route.extend({
  templateName: 'vault/cluster/oidc-callback',
  model() {
    // left blank so we render the template immediately
  },
  afterModel() {
    let { auth_path: path, code, state } = this.paramsFor(this.routeName);
    let { namespaceQueryParam: namespace } = this.paramsFor('vault.cluster');
    // namespace from state takes precedence over the cluster's ns
    if (state?.includes(',ns=')) {
      [state, namespace] = state.split(',ns=');
    }
    // some SSO providers do not return a url-encoded string, check for namespace using URLSearchParams
    const queryString = decodeURIComponent(window.location.search);
    const urlParams = new URLSearchParams(queryString);
    const checkState = urlParams.get('state');
    if (checkState?.includes(',ns=')) {
      [state, namespace] = checkState.split(',ns=');
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
