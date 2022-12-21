import PkiIssuerIndexRoute from './index';

export default class PkiIssuerSignRoute extends PkiIssuerIndexRoute {
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs.push({ label: resolvedModel.id, route: 'issuers.issuer.details' });
    controller.breadcrumbs.push({ label: 'sign intermediate' });
  }
}
