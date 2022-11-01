import Controller from '@ember/controller';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';

export default class OidcAssignmentDetailsController extends Controller {
  @service router;
  @service flashMessages;

  @action
  async delete() {
    try {
      await this.model.destroyRecord();
      this.flashMessages.success('Assignment deleted successfully');
      this.router.transitionTo('vault.cluster.access.oidc.assignments');
    } catch (error) {
      this.model.rollbackAttributes();
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.flashMessages.danger(message);
    }
  }
}
