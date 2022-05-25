import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';

export default class MfaLoginEnforcementIndexController extends Controller {
  @service router;
  @service flashMessages;

  queryParams = ['tab'];
  tab = 'targets';

  @tracked showDeleteConfirmation = false;
  @tracked deleteError;

  @action
  async delete() {
    try {
      await this.model.destroyRecord();
      this.showDeleteConfirmation = false;
      this.flashMessages.success('MFA login enforcement deleted successfully');
      this.router.transitionTo('vault.cluster.access.mfa.enforcements');
    } catch (error) {
      this.deleteError = error;
    }
  }
}
