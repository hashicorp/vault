import Component from '@glimmer/component';
import { action } from '@ember/object';

export default class VerifyComponent extends Component {
  @action
  onSubmit(...args) {
    this.args.doSubmit(...args);
  }

  @action
  clearParams(...args) {
    this.args.clearParams(...args);
  }
}
