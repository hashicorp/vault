import PkiIssuersIndexRoute from '../index';

export default class PkiIssuerDetailsRoute extends PkiIssuersIndexRoute {
  model() {
    const { issuer_ref } = this.paramsFor('issuers/issuer');
    return this.store.queryRecord('pki/issuer', {
      backend: this.secretMountPath.currentPath,
      id: issuer_ref,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const { id } = resolvedModel;
    const backend = this.secretMountPath.currentPath || 'pki';
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: backend, route: 'overview' },
      { label: 'issuers', route: 'issuers.index' },
      { label: id },
    ];
  }
}
