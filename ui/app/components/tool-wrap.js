import Component from '@glimmer/component';
import { action } from '@ember/object';

export default class HashTool extends Component {
  @action
  onClear() {
    this.args.onClear();
  }
  @action
  codemirrorUpdated() {
    this.args.codemirrorUpdated();
  }
  @action
  updateTtl() {
    this.args.updateTtl();
  }
}
