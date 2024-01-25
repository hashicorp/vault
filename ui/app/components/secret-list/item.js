import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

export default class SecretListItemComponent extends Component {
  @tracked showConfirmModal = false;
}
