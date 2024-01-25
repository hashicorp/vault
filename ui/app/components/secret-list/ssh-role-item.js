import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

export default class SecretListSshRoleItemComponent extends Component {
  @tracked showConfirmModal = false;
}
