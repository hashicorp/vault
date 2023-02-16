import PkiKeyRoute from '../key';

export default class PkiKeyDetailsRoute extends PkiKeyRoute {
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs.push({ label: resolvedModel.id });
  }
}
