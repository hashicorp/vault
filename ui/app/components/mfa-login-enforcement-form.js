import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { methods } from 'vault/helpers/mountable-auth-methods';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';

/**
 * @module MfaLoginEnforcementForm
 * MfaLoginEnforcementForm components are used to create and edit login enforcements
 *
 * @example
 * ```js
 * <MfaLoginEnforcementForm @model={{this.model}} @mfaMethod={{this.mfaMethod}} @hasActions={{true}} @onSave={{this.onSave}} @onClose={{this.onClose}} />
 * ```
 * @callback onSave
 * @callback onClose
 * @param {Object} model - login enforcement model
 * @param {Object} [mfaMethod] - provide when creating a login enforcement for a selected method -- otherwise search selector is displayed
 * @param {boolean} [hasActions] - whether the save and cancel actions will be displayed and handled internally or not
 * @param {onSave} [onSave] - triggered on save success
 * @param {onClose} [onClose] - triggered on cancel
 */

export default class MfaLoginEnforcementForm extends Component {
  @service store;
  @service flashMessages;

  targetTypes = [
    { label: 'Authentication mount', type: 'accessor', key: 'auth_method_accessors' },
    { label: 'Authentication method', type: 'method', key: 'auth_method_types' },
    { label: 'Group', type: 'identity/group', key: 'identity_groups' },
    { label: 'Entity', type: 'identity/entity', key: 'identity_entities' },
  ];
  authMethods = methods();
  searchSelectOptions = null;

  @tracked name;
  @tracked targets = [];
  @tracked selectedTargetType = 'accessor';
  @tracked selectedTargetValue = null;
  @tracked searchSelect = {
    options: [],
    selected: [],
  };

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
      const targets = targetArray.map((value) => ({ label, key, value }));
      this.targets.addObjects(targets);
    }
  }
  async resetTargetState() {
    this.selectedTargetValue = null;
    this.searchSelect.selected = [];
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
    if (this.selectedTargetType.includes('identity')) {
      this.searchSelect.options = [...options[this.selectedTargetType]];
    }
  }

  get selectedTarget() {
    return this.targetTypes.findBy('type', this.selectedTargetType);
  }

  @task
  *save() {
    try {
      yield this.args.model.save();
      this.args.onSave();
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.flashMessages.danger(message);
    }
  }

  @action
  async onMethodChange(selectedIds) {
    const methods = await this.args.model.mfa_methods;
    // first check for existing methods that have been removed from selection
    methods.forEach((method) => {
      if (!selectedIds.includes(method.id)) {
        methods.removeObject(method);
      }
    });
    // now check for selected items that don't exist and add them to the model
    const methodIds = methods.mapBy('id');
    selectedIds.forEach((id) => {
      if (!methodIds.includes(id)) {
        const model = this.store.peekRecord('mfa-method', id);
        methods.addObject(model);
      }
    });
  }
  @action
  onTargetSelect(type) {
    this.selectedTargetType = type;
    this.resetTargetState();
  }
  @action
  setTargetValue(selected) {
    const { type } = this.selectedTarget;
    if (type.includes('identity')) {
      // for identity groups and entities grab model from store as value
      this.selectedTargetValue = this.store.peekRecord(type, selected[0]);
    } else {
      this.selectedTargetValue = selected;
    }
  }
  @action
  addTarget() {
    const { label, key } = this.selectedTarget;
    const value = this.selectedTargetValue;
    this.targets.addObject({ label, value, key });
    // add target to appropriate model property
    this.args.model[key].addObject(value);
    this.selectedTargetValue = null;
    this.resetTargetState();
  }
  @action
  removeTarget(target) {
    this.targets.removeObject(target);
    // remove target from appropriate model property
    this.args.model[target.key].addObject(target.value);
  }
  @action
  cancel() {
    // revert model changes
    this.args.model.rollbackAttributes();
    this.args.onClose();
  }
}
