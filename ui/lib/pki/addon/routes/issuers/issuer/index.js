import PkiIssuersListRoute from '../index';

// Single issuer index route extends issuers list route
export default class PkiIssuerIndexRoute extends PkiIssuersListRoute {
  model() {
    const { issuer_ref } = this.paramsFor('issuers/issuer');
    return this.store.queryRecord('pki/issuer', {
      backend: this.secretMountPath.currentPath,
      id: issuer_ref,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'issuers', route: 'issuers.index' },
    ];
  }
}
