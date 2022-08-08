import Controller from '@ember/controller';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
export default class OidcAssignmentDetailsController extends Controller {
  @service router;
  @service flashMessages;

  queryParams = ['listEntities', 'listGroups'];

  @tracked listEntities = null;
  @tracked listGroups = null;
  @tracked model; // ARG TODO following example, but try without once have working

  @action
  async delete() {
    try {
      await this.model.destroyRecord();
      this.flashMessages.success('Assignment deleted successfully');
      this.router.transitionTo('vault.cluster.access.oidc.scopes');
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.flashMessages.danger(message);
    }
  }
}
