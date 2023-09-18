import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';

export default class IdentityPageIndexComponent extends Component {
  @service flashMessages;
  @service store;
  @service router;

  @action
  delete(model) {
    const type = model.get('identityType');
    const id = model.id;
    return model
      .destroyRecord()
      .then(() => {
        this.store.clearDataset('identity/entity');
        this.flashMessages.success(`Successfully deleted ${type}: ${id}`);
        // Transition to refresh the model
        this.router.transitionTo('vault.cluster.access.identity.index');
      })
      .catch((e) => {
        this.flashMessages.danger(
          `There was a problem deleting ${type}: ${id} - ${e.errors.join(' ') || e.message}`
        );
      });
  }

  @action
  toggleDisabled(model) {
    const action = model.disabled ? ['enabled', 'enabling'] : ['disabled', 'disabling'];
    const type = model.identityType;
    const id = model.id;
    model.toggleProperty('disabled');

    model
      .save()
      .then(() => {
        this.flashMessages.success(`Successfully ${action[0]} ${type}: ${id}`);
      })
      .catch((e) => {
        this.flashMessages.danger(
          `There was a problem ${action[1]} ${type}: ${id} - ${e.errors.join(' ') || e.message}`
        );
      });
  }

  @action reloadRecord(model) {
    model.reload();
  }
}
