import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';

export default class OidcConfigureController extends Controller {
  @tracked header = null;

  get isCta() {
    return this.header === 'cta';
  }
}
