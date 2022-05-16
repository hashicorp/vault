import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';

export default class MfaMethodForm extends Component {
  @service store;
  @service flashMessages;

  @task
  *save() {
    try {
      yield this.args.model.save();
      this.args.onSave();
    } catch (e) {
      this.flashMessages.danger(e.errors?.join('. ') || e.message);
    }
  }
  @action
  cancel() {
    // revert model changes
    this.args.model.rollbackAttributes();
    this.args.onClose();
  }
}
