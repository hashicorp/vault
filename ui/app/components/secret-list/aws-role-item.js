import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

export default class SecretListAwsRoleItemComponent extends Component {
  @tracked showConfirmModal = false;
}
