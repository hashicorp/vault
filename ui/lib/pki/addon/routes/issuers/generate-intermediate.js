import PkiIssuersIndexRoute from '.';

export default class PkiIssuersGenerateIntermediateRoute extends PkiIssuersIndexRoute {
  model() {
    return this.store.createRecord('pki/action');
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs.push({ label: 'generate CSR' });
  }
}
