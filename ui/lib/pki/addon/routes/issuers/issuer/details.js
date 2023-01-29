import PkiIssuerIndexRoute from './index';

export default class PkiIssuerDetailsRoute extends PkiIssuerIndexRoute {
  // Details route gets issuer data from PkiIssuerIndexRoute
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs.push({ label: resolvedModel.id });
  }
}
