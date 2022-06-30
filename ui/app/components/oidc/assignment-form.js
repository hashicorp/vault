import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import handleHasManySelection from 'core/utils/search-select-has-many';

/**
 * @module OidcAssignmentForm
 * OidcAssignmentForm components are used to...
 *
 * @example
 * ```js
 * <OidcAssignmentForm @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export default class OidcAssignmentForm extends Component {
  @service store;
  @service flashMessages;

  @tracked modelErrors;

  get errors() {
    return this.args.modelErrors || this.modelErrors;
  }

  @task
  *save() {
    this.modelErrors = {};
    // check validity state first and abort if invalid
    const { isValid, state } = this.args.model.validate();
    if (!isValid) {
      this.modelErrors = state;
    } else {
      try {
        yield this.args.model.save();
        this.args.onSave();
      } catch (error) {
        console.log(error, 'error');
        const message = error.errors ? error.errors.join('. ') : error.message;
        this.flashMessages.danger(message);
      }
    }
  }

  @action
  cancel() {
    // revert model changes
    this.args.model.rollbackAttributes();
    this.args.onClose();
  }

  @action
  handleOperation(e) {
    let value = e.target.value;
    this.args.model.name = value;
  }

  @action
  onEntitiesSelect(selectedIds) {
    const entityIds = this.args.model.entityIds;
    handleHasManySelection(selectedIds, entityIds, this.store, 'identity/entity');
  }

  @action
  onGroupsSelect(selectedIds) {
    const groupIds = this.args.model.GroupIds;
    handleHasManySelection(selectedIds, groupIds, this.store, 'identity/group');
  }
}
