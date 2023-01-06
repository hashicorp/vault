import PkiIssuersIndexRoute from '.';

export default class PkiIssuersImportRoute extends PkiIssuersIndexRoute {
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs.push({ label: 'import' });
  }
}
