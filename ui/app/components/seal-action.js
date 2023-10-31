import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';

export default class SealActionComponent extends Component {
  @tracked error;

  @action
  async handleSeal() {
    try {
      await this.args.onSeal();
    } catch (e) {
      this.error = errorMessage(e, 'Seal attempt failed. Check Vault logs for details.');
    }
  }
}
