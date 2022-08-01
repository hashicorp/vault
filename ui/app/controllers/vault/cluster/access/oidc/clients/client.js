import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class OidcClientController extends Controller {
  @service router;
  @tracked showHeader = true;

  constructor() {
    super(...arguments);
    this.router.on('routeDidChange', (transition) => this.showOrHideHeader(transition));
  }

  showOrHideHeader({ targetName }) {
    // hide header when rendering the edit form (client-form component has separate header)
    this.showHeader = targetName.includes('edit') ? false : true;
  }
}
