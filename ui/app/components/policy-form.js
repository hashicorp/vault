import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import trimRight from 'vault/utils/trim-right';

/**
 * @module PolicyForm
 * PolicyForm components are used to display the create and edit forms for all types of policies
 *
 * @example
 *  <PolicyForm
 *    @model={{this.model}}
 *    @onSave={{transition-to "vault.cluster.policy.show" this.model.policyType this.model.name}}
 *    @onCancel={{transition-to "vault.cluster.policies.index"}}
 *  />
 * ```
 * @callback onCancel
 * @callback onSave
 * @param {object} model - The parent's model
 * @param {string} onCancel - callback triggered when cancel button is clicked
 * @param {string} onSave - callback triggered when save button is clicked
 */

export default class PolicyFormComponent extends Component {
  @service store;
  @service wizard;
  @service flashMessages;
  @tracked errorBanner;
  @tracked showFileUpload = false;
  @tracked file = null;

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isNew, name, policyType } = this.args.model;
      yield this.args.model.save();
      this.flashMessages.success(
        `${policyType.toUpperCase()} policy "${name}" was successfully ${isNew ? 'created' : 'updated'}.`
      );
      if (this.wizard.featureState === 'create') {
        this.wizard.transitionFeatureMachine('create', 'CONTINUE', policyType);
      }

      // this form is sometimes used in modal, passing the model notifies
      // the parent if the save was successful
      this.args.onSave(this.args.model);
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.errorBanner = message;
    }
  }

  @action
  cancel() {
    const method = this.args.model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.args.model[method]();
    this.args.onCancel();
  }

  @action
  setModelName({ target }) {
    this.args.model.name = target.value.toLowerCase();
  }

  @action
  toggleFileUpload() {
    this.showFileUpload = !this.showFileUpload;
  }

  @action
  setPolicyFromFile(index, fileInfo) {
    let { value, fileName } = fileInfo;
    let model = this.args.model;
    model.policy = value;
    if (!model.name) {
      let trimmedFileName = trimRight(fileName, ['.json', '.txt', '.hcl', '.policy']);
      model.name = trimmedFileName;
    }
    this.showFileUpload = false;
  }
}
