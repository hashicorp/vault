import { service } from '@ember/service';
import Component from '@glimmer/component';

export default class ResultantAclBannerComponent extends Component {
  @service namespace;

  get ns() {
    return this.namespace.path;
  }
}
