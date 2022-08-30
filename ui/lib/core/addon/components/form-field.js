import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { capitalize } from 'vault/helpers/capitalize';
import { humanize } from 'vault/helpers/humanize';
import { dasherize } from 'vault/helpers/dasherize';
/**
 * @module FormField
 * `FormField` components are field elements associated with a particular model.
 *
 * @example
 * ```js
 * {{#each @model.fields as |attr|}}
 *  <FormField data-test-field @attr={{attr}} @model={{this.model}} />
 * {{/each}}
 * ```
 * example attr object:
 * attr = {
 *   name: "foo", // name of attribute -- used to populate various fields and pull value from model
 *   options: {
 *     label: "Foo", // custom label to be shown, otherwise attr.name will be displayed
 *     defaultValue: "", // default value to display if model value is not present
 *     fieldValue: "bar", // used for value lookup on model over attr.name
 *     editType: "ttl", type of field to use -- example boolean, searchSelect, etc.
 *     helpText: "This will be in a tooltip",
 *     readOnly: true
 *   },
 *   type: "boolean" // type of attribute value -- string, boolean, etc.
 * }
 * @param {Object} attr - usually derived from ember model `attributes` lookup, and all members of `attr.options` are optional
 * @param {Model} model - Ember Data model that `attr` is defined on
 * @param {boolean} [disabled=false] - whether the field is disabled
 * @param {boolean} [showHelpText=true] - whether to show the tooltip with help text from OpenAPI
 * @param {string} [subText] - text to be displayed below the label
 * @param {string} [mode] - used when editType is 'kv'
 * @param {ModelValidations} [modelValidations] - Object of errors.  If attr.name is in object and has error message display in AlertInline.
 * @callback onChangeCallback
 * @param {onChangeCallback} [onChange] - called whenever a value on the model changes via the component
 * @callback onKeyUpCallback
 * @param {onKeyUpCallback} [onKeyUp] - function passed through into MaskedInput to handle validation. It is also handled for certain form-field types here in the action handleKeyUp.
 *
 */

export default class FormFieldComponent extends Component {
  @tracked showInput = false;
  @tracked file = { value: '' }; // used by the pgp-file component when an attr is editType of 'file'
  emptyData = '{\n}';

  constructor() {
    super(...arguments);
    const { attr, model } = this.args;
    const valuePath = attr.options?.fieldValue || attr.name;
    const modelValue = model[valuePath];
    this.showInput = !!modelValue;
  }

  get disabled() {
    return this.args.disabled || false;
  }
  get showHelpText() {
    return this.args.showHelpText || true;
  }
  get subText() {
    return this.args.subText || '';
  }
  // used in the label element next to the form element
  get labelString() {
    const label = this.args.attr.options?.label || '';
    if (label) {
      return label;
    }
    if (this.args.attr.name) {
      return capitalize([humanize([dasherize([this.args.attr.name])])]);
    }
    return '';
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

  onChange() {
    if (this.args.onChange) {
      this.args.onChange(...arguments);
    }
  }

  @action
  setFile(_, keyFile) {
    const path = this.valuePath;
    const { value } = keyFile;
    this.args.model.set(path, value);
    this.onChange(path, value);
    this.file = keyFile;
  }
  @action
  setAndBroadcast(value) {
    this.args.model.set(this.valuePath, value);
    this.onChange(this.valuePath, value);
  }
  @action
  setAndBroadcastBool(trueVal, falseVal, event) {
    let valueToSet = event.target.checked === true ? trueVal : falseVal;
    this.setAndBroadcast(valueToSet);
  }
  @action
  setAndBroadcastTtl(value) {
    const alwaysSendValue = this.valuePath === 'expiry' || this.valuePath === 'safetyBuffer';
    let valueToSet = value.enabled === true || alwaysSendValue ? `${value.seconds}s` : 0;
    this.setAndBroadcast(`${valueToSet}`);
  }
  @action
  codemirrorUpdated(isString, value, codemirror) {
    codemirror.performLint();
    const hasErrors = codemirror.state.lint.marked.length > 0;
    let valToSet = isString ? value : JSON.parse(value);

    if (!hasErrors) {
      this.args.model.set(this.valuePath, valToSet);
      this.onChange(this.valuePath, valToSet);
    }
  }
  @action
  toggleShow() {
    const value = !this.showInput;
    this.showInput = value;
    if (!value) {
      this.setAndBroadcast(null);
    }
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
