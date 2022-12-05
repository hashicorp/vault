import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { getOwner } from '@ember/application';
import errorMessage from 'vault/utils/error-message';

export default class ConfigurePageComponent extends Component {
  @service flashMessages;

  get mountPoint() {
    return getOwner(this).mountPoint;
  }

  @task
  @waitFor
  *onDelete(model) {
    try {
      yield model.destroyRecord();
      this.args.model.roles.removeObject(model);
    } catch (error) {
      const message = errorMessage(error, 'Error deleting role. Please try again or contact support');
      this.flashMessages.danger(message);
    }
  }
}
