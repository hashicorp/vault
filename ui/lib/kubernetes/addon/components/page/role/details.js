import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';

export default class RoleDetailsPageComponent extends Component {
  @service router;
  @service flashMessages;

  get extraFields() {
    const fields = [];
    if (this.args.model.extraAnnotations) {
      fields.push({ label: 'Annotations', key: 'extraAnnotations' });
    }
    if (this.args.model.extraLabels) {
      fields.push({ label: 'Labels', key: 'extraLabels' });
    }
    return fields;
  }

  @action
  async delete() {
    try {
      await this.args.model.destroyRecord();
      this.router.transitionTo('vault.cluster.secrets.backend.kubernetes.roles');
    } catch (error) {
      const message = errorMessage(error, 'Unable to delete role. Please try again or contact support');
      this.flashMessages.danger(message);
    }
  }
}
