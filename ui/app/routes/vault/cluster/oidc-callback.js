import Route from '@ember/routing/route';

export default Route.extend({
  templateName: 'vault/cluster/oidc-callback',
  model() {
    // left blank so we render the template immediately
  },
  afterModel() {
    let { auth_path: path, code } = this.paramsFor(this.routeName);
    // some SSO providers do not return a url-encoded state param
    // parse state using URLSearchParams instead of paramsFor
    const queryString = decodeURIComponent(window.location.search);
    const urlParams = new URLSearchParams(queryString);
    let state = urlParams.get('state');
    // check cluster for a namespace, i.e. HCP namespace flag
    let { namespaceQueryParam: namespace } = this.paramsFor('vault.cluster');
    // namespace from state takes precedence over the cluster's ns
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
