import Controller from '@ember/controller';
import { inject as service } from '@ember/service';

export default class ConfigurationController extends Controller {
  @service router;

  get currentRouteName() {
    return this.router.currentRouteName;
  }
}
