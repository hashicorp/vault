import Component from '@glimmer/component';
import { action } from '@ember/object';

export default class HmacComponent extends Component {
  @action
  onSubmit(...args) {
    this.args.doSubmit(...args);
  }

  @action
  toggleModal(data) {
    this.args.toggleModal(data);
  }
}
