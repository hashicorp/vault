/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { capitalize } from 'vault/helpers/capitalize';
import { humanize } from 'vault/helpers/humanize';
import { dasherize } from 'vault/helpers/dasherize';
import { assert } from '@ember/debug';
import { addToArray } from 'vault/helpers/add-to-array';
import { removeFromArray } from 'vault/helpers/remove-from-array';
import { isEmpty } from '@ember/utils';
import { presence } from 'vault/utils/forms/validators';
import { get } from '@ember/object';

/**
 * @module FormField
 * FormField components are field elements associated with a particular model.
 * @description
 * ```
 * sample attr shape:
 * attr = {
 * name: "foo", // name of attribute -- used to populate various fields and pull value from model
 * type: "boolean" // type of attribute value -- string, boolean, etc.
 * options: {
 *  label: "To do task", // custom label to be shown, otherwise attr.name will be displayed
 *  defaultValue: "", // default value to display if model value is not present
 *  fieldValue: "toDo", // used for value lookup on model over attr.name
 *  editType: "boolean", type of field to use. List of editTypes:boolean, file, json, kv, optionalText, mountAccessor, password, radio, regex, searchSelect, stringArray, textarea, ttl, yield.
 *  helpText: "This will be in a tooltip",
 *  readOnly: true
 *  },
 * }
 * ```
 *
 * @example
 * {{#each @model.fields as |attr|}}
 *  <FormField data-test-field @attr={{attr}} @model={{@model}} />
 * {{/each}}
 * <FormField @attr={{hash name="toDo" options=(hash label="To do task" helpText="helpful text" editType="boolean")}} @model={{hash toDo=true}} />
 *
 * @param {Object} attr - usually derived from ember model `attributes` lookup, and all members of `attr.options` are optional
 * @param {Model} model - Ember Data model that `attr` is defined on
 * @param {boolean} [disabled=false] - whether the field is disabled
 * @param {boolean} [showHelpText=true] - whether to show the tooltip with help text from OpenAPI
 * @param {string} [mode] - used when editType is 'kv'
 * @param {object} [modelValidations] - Object of errors.  If attr.name is in object and has error message display in AlertInline.
 * @param {function} [onChange] - called whenever a value on the model changes via the component
 * @param {function} [onKeyUp] - function passed through into MaskedInput to handle validation. It is also handled for certain form-field types here in the action handleKeyUp.
 *
 */

export default class FormFieldComponent extends Component {
  emptyData = '{\n}';
  shouldHideLabel = [
    'file',
    'json',
    'kv',
    'mountAccessor',
    'optionalText',
    'regex',
    'searchSelect',
    'stringArray',
    'ttl',
    'toggleButton',
  ];
  @tracked showToggleTextInput = false;
  @tracked toggleInputEnabled = false;

  radioValue = (item) => (isEmpty(item.value) ? item : item.value);

  constructor() {
    super(...arguments);
    const { attr, model } = this.args;
    const valuePath = attr.options?.fieldValue || attr.name;
    assert(
      'Form is attempting to modify an ID. Ember-data does not allow this.',
      valuePath.toLowerCase() !== 'id'
    );
    assert('@name is required', presence(attr.name));
    assert('@model (or resource object being updated) is required', presence(model));
    const modelValue = get(model, valuePath);
    this.showToggleTextInput = !!modelValue;
    this.toggleInputEnabled = !!modelValue;
  }

  // ---------------------------------------------------------------
  // IMPORTANT:
  // this controls the top-level logic in the associated template
  // to decide which type of form field to render (Vault or HDS)
  // ---------------------------------------------------------------
  //
  get isHdsFormField() {
    const { type, options } = this.args.attr;

    // here we replicate the logic in the `form-field.hbs` template, as it was at the beginning of the migation to HDS
    // to make sure we don't change the order in which the "ifs" are evaluated
    // see: https://github.com/hashicorp/vault/blob/e99c06a0249f1ecde02c5b48990cb0e91e4ec575/ui/lib/core/addon/components/form-field.hbs
    if (options?.possibleValues?.length > 0) {
      return true;
    } else {
      if (options?.editType === 'dateTimeLocal') {
        return true;
      } else if (
        options?.editType === 'searchSelect' ||
        options?.editType === 'mountAccessor' ||
        options?.editType === 'kv' ||
        options?.editType === 'file' ||
        options?.editType === 'ttl' ||
        options?.editType === 'regex' ||
        options?.editType === 'toggleButton' ||
        options?.editType === 'optionalText' ||
        options?.editType === 'stringArray' ||
        options?.sensitive === true
      ) {
        return false;
      } else if (type === 'number' || type === 'string') {
        if (options?.editType === 'textarea' || options?.editType === 'password') {
          return true;
        } else {
          return false;
        }
      } else if (type === 'boolean' || options?.editType === 'boolean') {
        return true;
      } else if (type === 'object' || options?.editType === 'yield') {
        return false;
      } else {
        // fallback, just in case
        return false;
      }
    }
  }

  get hasRadioSubText() {
    // for 'radio' editType, check to see if any of the possibleValues has a subText
    return this.args?.attr?.options?.possibleValues?.any((v) => v.subText);
  }

  get hasRadioHelpText() {
    // for 'radio' editType, check to see if any of the possibleValues has a helpText
    return this.args?.attr?.options?.possibleValues?.any((v) => v.helpText);
  }

  get hideLabel() {
    const { type, options } = this.args.attr;
    if (type === 'object' || options?.isSectionHeader) {
      return true;
    }
    // falsey values render a <FormFieldLabel>
    return this.shouldHideLabel.includes(options?.editType);
  }

  get disabled() {
    return this.args.disabled || false;
  }

  get helpTextString() {
    const helpText = this.args.attr?.options?.helpText;
    if (this.args.showHelpText !== false && helpText) {
      return helpText;
    }
    return '';
  }

  // used in the label element next to the form element
  get labelString() {
    const label = this.args.attr.options?.label || '';
    return label ? label : capitalize([humanize([dasherize([this.args.attr.name])])]);
  }

  // both the path to mutate on the model, and the path to read the value from
  get valuePath() {
    return this.args.attr.options?.fieldValue || this.args.attr.name;
  }

  get isReadOnly() {
    const readonly = this.args.attr.options?.readOnly || false;
    return readonly && this.args.mode === 'edit';
  }

  get validationError() {
    const validations = this.args.modelValidations || {};
    const state = validations[this.valuePath];
    return state && !state.isValid ? state.errors.join(' ') : null;
  }

  get validationWarning() {
    const validations = this.args.modelValidations || {};
    const state = validations[this.valuePath];
    return state?.warnings?.length ? state.warnings.join(' ') : null;
  }

  onChange() {
    if (this.args.onChange) {
      this.args.onChange(...arguments);
    }
  }

  @action
  setFile(keyFile) {
    const path = this.valuePath;
    const { value } = keyFile;
    this.args.model.set(path, value);
    this.onChange(path, value);
  }
  @action
  setAndBroadcast(value) {
    this.args.model.set(this.valuePath, value);
    this.onChange(this.valuePath, value);
  }
  @action
  setAndBroadcastBool(trueVal, falseVal, event) {
    const valueToSet = event.target.checked === true ? trueVal : falseVal;
    this.setAndBroadcast(valueToSet);
  }
  @action
  setAndBroadcastChecklist(event) {
    let updatedValue = this.args.model[this.valuePath];
    if (event.target.checked) {
      updatedValue = addToArray(updatedValue, event.target.value);
    } else {
      updatedValue = removeFromArray(updatedValue, event.target.value);
    }
    this.setAndBroadcast(updatedValue);
  }
  @action
  setAndBroadcastRadio(item) {
    // we want to read the original value instead of `event.target.value` so we have `false` (boolean) and not `"false"` (string)
    const valueToSet = this.radioValue(item);
    this.setAndBroadcast(valueToSet);
  }
  @action
  setAndBroadcastTtl(value) {
    const alwaysSendValue = this.valuePath === 'expiry' || this.valuePath === 'safetyBuffer';
    const attrOptions = this.args.attr.options || {};
    let valueToSet = 0;
    if (value.enabled || alwaysSendValue) {
      valueToSet = `${value.seconds}s`;
    } else if (Object.keys(attrOptions).includes('ttlOffValue')) {
      valueToSet = attrOptions.ttlOffValue;
    }
    this.setAndBroadcast(`${valueToSet}`);
  }
  @action
  codemirrorUpdated(isString, value, codemirror) {
    codemirror.performLint();
    const hasErrors = codemirror.state.lint.marked.length > 0;
    const valToSet = isString ? value : JSON.parse(value);

    if (!hasErrors) {
      this.args.model.set(this.valuePath, valToSet);
      this.onChange(this.valuePath, valToSet);
    }
  }
  @action
  toggleTextShow() {
    const value = !this.showToggleTextInput;
    this.showToggleTextInput = value;
    if (!value) {
      this.setAndBroadcast(null);
    }
  }
  @action
  toggleButton() {
    this.toggleInputEnabled = !this.toggleInputEnabled;
    this.setAndBroadcast(this.toggleInputEnabled);
  }
  @action
  handleKeyUp(maybeEvent) {
    const value = typeof maybeEvent === 'object' ? maybeEvent.target.value : maybeEvent;
    if (!this.args.onKeyUp) {
      return;
    }
    this.args.onKeyUp(this.valuePath, value);
  }
  @action
  onChangeWithEvent(event) {
    const prop = event.target.type === 'checkbox' ? 'checked' : 'value';
    this.setAndBroadcast(event.target[prop]);
  }
}
