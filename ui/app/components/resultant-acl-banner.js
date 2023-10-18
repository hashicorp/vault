import { service } from '@ember/service';
import Component from '@glimmer/component';

export default class ResultantAclBannerComponent extends Component {
  @service namespace;
  @service router;

  get ns() {
    return this.namespace.path || 'root';
  }

  get queryParams() {
    // Bring user back to current page after login
    return { redirect_to: this.router.currentURL };
  }
}
