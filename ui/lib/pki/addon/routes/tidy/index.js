import Route from '@ember/routing/route';
import { PKI_DEFAULT_EMPTY_STATE_MSG } from '../overview';

export default class PkiTidyIndexRoute extends Route {
  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.notConfiguredMessage = PKI_DEFAULT_EMPTY_STATE_MSG;
  }
}
