import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import handleHasManySelection from 'core/utils/search-select-has-many';

/**
 * @module Oidc::AssignmentForm
 * Oidc::AssignmentForm components are used to display the create view for OIDC providers assignments.
 *
 * @example
 * ```js
 * <Oidc::AssignmentForm @model={this.model}
 * @onCancel={transition-to "vault.cluster.access.oidc.assignment"} @param1={{param1}}
 * @onSave={transition-to "vault.cluster.access.oidc.assignments.assignment.details" this.model.name}
 * />
 * ```
 * @callback onCancel
 * @callback onSave
 * @param {object} model - The parent's model
 * @param {string} onCancel - callback triggered when cancel button is clicked
 * @param {string} onSave - callback triggered when save button is clicked
 */

export default class OidcAssignmentFormComponent extends Component {
  @service store;
  @service flashMessages;

  targetTypes = [
    { label: 'Entity', type: 'identity/entity', key: 'entity_ids' },
    { label: 'Group', type: 'identity/group', key: 'groups_ids' },
  ];

  @tracked modelValidations;
  @tracked targets = [];

  constructor() {
    super(...arguments);
    // aggregate different target array properties on model into flat list
    this.flattenTargets();
    // eagerly fetch identity groups and entities for use as search select options
    this.resetTargetState();
  }

  async flattenTargets() {
    for (let { label, key } of this.targetTypes) {
      const targetArray = await this.args.model[key];
      if (typeof targetArray !== 'object') {
        return;
      }
      const targets = targetArray.map((value) => ({ label, key, value }));
      this.targets.addObjects(targets);
    }
  }
  async resetTargetState() {
    this.selectedTargetValue = null;
    const options = this.searchSelectOptions || {};
    if (!this.searchSelectOptions) {
      const types = ['identity/group', 'identity/entity'];
      for (const type of types) {
        try {
          options[type] = (await this.store.query(type, {})).toArray();
        } catch (error) {
          options[type] = [];
        }
      }
      this.searchSelectOptions = options;
    }
    // if (this.selectedTargetType.includes('identity')) {
    //   this.searchSelect = {
    //     selected: [],
    //     options: [...options[this.selectedTargetType]],
    //   };
    // }
  }

  get selectedTarget() {
    return this.targetTypes.findBy('type', this.selectedTargetType);
  }

  @task
  *save() {
    event.preventDefault();
    try {
      const { isValid, state } = this.args.model.validate();
      this.modelValidations = isValid ? null : state;
      if (isValid) {
        yield this.args.model.save();
        this.flashMessages.success('Successfully created an assignment');
        this.args.onSave();
      }
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.flashMessages.danger(message);
    }
  }

  @action
  cancel() {
    const method = this.args.model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.args.model[method]();
    this.args.onCancel();
  }

  @action
  handleOperation(e) {
    let value = e.target.value;
    this.args.model.name = value;
  }

  @action
  onEntitiesSelect(selectedIds) {
    const entityIds = this.args.model.entity_ids;
    handleHasManySelection(selectedIds, entityIds, this.store, 'identity/entity');
  }

  @action
  onGroupsSelect(selectedIds) {
    const groupIds = this.args.model.groupIds;
    handleHasManySelection(selectedIds, groupIds, this.store, 'identity/group');
  }
}
