import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { getOwner } from '@ember/application';
import errorMessage from 'vault/utils/error-message';

export default class ConfigurePageComponent extends Component {
  @service flashMessages;

  get mountPoint() {
    return getOwner(this).mountPoint;
  }

  @action
  async onDelete(model) {
    try {
      const message = `Successfully deleted role ${model.name}`;
      await model.destroyRecord();
      this.args.roles.removeObject(model);
      this.flashMessages.success(message);
    } catch (error) {
      const message = errorMessage(error, 'Error deleting role. Please try again or contact support');
      this.flashMessages.danger(message);
    }
  }
}
