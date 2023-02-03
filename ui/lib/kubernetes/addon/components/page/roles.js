import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { getOwner } from '@ember/application';
import errorMessage from 'vault/utils/error-message';

/**
 * @module Roles
 * RolesPage component is a child component to show list of roles
 *
 * @param {array} roles - array of roles
 * @param {object} config - config model that contains kubernetes configuration
 * @param {array} pageFilter - array of filtered roles
 * @param {array} breadcrumbs - breadcrumbs as an array of objects that contain label and route
 */
export default class RolesPageComponent extends Component {
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
