import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';

/**
 * MfaMethodForm component
 *
 * @example
 * ```js
 * <Mfa::MethodForm @model={{this.model}} @hasActions={{true}} @onSave={{this.onSave}} @onClose={{this.onClose}} />
 * ```
 * @param {Object} model - MFA method model
 * @param {boolean} [hasActions] - whether the action buttons will be rendered or not
 * @param {onSave} [onSave] - callback when save is successful
 * @param {onClose} [onClose] - callback when cancel is triggered
 */
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
    this.args.model.rollbackAttributes();
    this.args.onClose();
  }
}
