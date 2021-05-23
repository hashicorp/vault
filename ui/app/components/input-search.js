import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class inputSelect extends Component {
  /*
   * @public
   * @param Function
   *
   * Function called when any of the inputs change
   *
   */
  @tracked searchInput = '';
  @action
  inputChanged() {
    this.args.onChange(this.searchInput);
  }
}
