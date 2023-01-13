import PkiIssuerIndexRoute from './index';
// TODO comment in before merging
// import { withConfirmLeave } from 'core/decorators/confirm-leave';

// @withConfirmLeave()
export default class PkiIssuerCrossSignRoute extends PkiIssuerIndexRoute {
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs.push({ label: 'cross-sign' });
  }
}
