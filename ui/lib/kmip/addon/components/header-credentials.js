import Component from '@glimmer/component';
import { inject as service } from '@ember/service';

export default class HeaderCredentialsComponent extends Component {
  @service secretMountPath;

  get scope() {
    return this.args.scope || null;
  }
  get role() {
    return this.args.role || null;
  }
}
