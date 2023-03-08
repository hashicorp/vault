import Component from '@glimmer/component';
import { action } from '@ember/object';

export default class InfoTooltip extends Component {
  @action
  preventSubmit(e) {
    e.preventDefault();
  }
}
