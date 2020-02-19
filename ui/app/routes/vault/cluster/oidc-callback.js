import Route from '@ember/routing/route';

export default Route.extend({
  templateName: 'vault/cluster/oidc-callback',
  model() {
    // left blank so we render the template immediately
  },
  afterModel() {
    let { auth_path: path, code, state } = this.paramsFor(this.routeName);
    let { namespaceQueryParam: namespace } = this.paramsFor('vault.cluster');
    path = window.decodeURIComponent(path);
    let queryParams = { namespace, path, code, state };
    window.localStorage.setItem('oidcState', JSON.stringify(queryParams));
  },
  renderTemplate() {
    this.render(this.templateName, {
      into: 'application',
      outlet: 'main',
    });
  },
});
