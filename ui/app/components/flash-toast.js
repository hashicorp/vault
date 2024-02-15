import { capitalize } from '@ember/string';
import Component from '@glimmer/component';

export default class FlashToastComponent extends Component {
  get color() {
    switch (this.args.flash.type) {
      case 'info':
        return 'highlight';
      case 'danger':
        return 'critical';
      default:
        return this.args.flash.type;
    }
  }

  get title() {
    if (this.args.title) return this.args.title;
    switch (this.args.flash.type) {
      case 'danger':
        return 'Error';
      default:
        return capitalize(this.args.flash.type);
    }
  }
}
