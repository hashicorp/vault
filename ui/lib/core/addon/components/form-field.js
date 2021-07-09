import { or } from '@ember/object/computed';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { capitalize } from 'vault/helpers/capitalize';
import { humanize } from 'vault/helpers/humanize';
import { dasherize } from 'vault/helpers/dasherize';
import layout from '../templates/components/form-field';

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
 *
 * @param [onChange=null] {Func} - Called whenever a value on the model changes via the component.
 * @param [onKeyUp=null] {Func} - Called whenever cp-validations is being used and you need to validation on keyup.  Send name of field and value of input.
 * @param attr=null {Object} - This is usually derived from ember model `attributes` lookup, and all members of `attr.options` are optional.
 * @param model=null {DS.Model} - The Ember Data model that `attr` is defined on
 * @param [disabled=false] {Boolean} - whether the field is disabled
 * @param [showHelpText=true] {Boolean} - whether to show the tooltip with help text from OpenAPI
 * @param [subText] {String} - Text to be displayed below the label
 * @param [validationMessages] {Object} - Object of errors.  If attr.name is in object and has error message display in AlertInline.
 *
 */

export default Component.extend({
  layout,
  'data-test-field': true,
  classNames: ['field'],
  disabled: false,
  showHelpText: true,
  subText: '',
  // This is only used internally for `optional-text` editType
  showInput: false,

  init() {
    this._super(...arguments);
    const valuePath = this.attr.options?.fieldValue || this.attr.name;
    const modelValue = this.model[valuePath];
    this.set('showInput', !!modelValue);
  },

  onChange() {},

  /*
   * @public
   * @param Object
   * in the form of
   * {
   *   name: "foo",
   *   options: {
   *     label: "Foo",
   *     defaultValue: "",
   *     editType: "ttl",
   *     helpText: "This will be in a tooltip"
   *   },
   *   type: "boolean"
   * }
   *
   */
  attr: null,

  mode: null,

  /*
   * @private
   * @param string
   * Computed property used in the label element next to the form element
   *
   */
  labelString: computed('attr.{name,options.label}', function() {
    const label = this.attr.options ? this.attr.options.label : '';
    const name = this.attr.name;
    if (label) {
      return label;
    }
    if (name) {
      return capitalize([humanize([dasherize([name])])]);
    }
    return '';
  }),

  // both the path to mutate on the model, and the path to read the value from
  /*
   * @private
   * @param string
   *
   * Computed property used to set values on the passed model
   *
   */
  valuePath: or('attr.options.fieldValue', 'attr.name'),

  isReadOnly: computed('attr.options.readOnly', 'mode', function() {
    let readonly = this.attr.options ? this.attr.options.readOnly : false;
    return readonly && this.mode === 'edit';
  }),

  model: null,

  /*
   * @private
   * @param object
   *
   * Used by the pgp-file component when an attr is editType of 'file'
   */
  file: computed(function() {
    return { value: '' };
  }),
  emptyData: '{\n}',

  actions: {
    setFile(_, keyFile) {
      const path = this.valuePath;
      const { value } = keyFile;
      this.model.set(path, value);
      this.onChange(path, value);
      this.set('file', keyFile);
    },

    setAndBroadcast(path, value) {
      this.model.set(path, value);
      this.onChange(path, value);
    },

    setAndBroadcastBool(path, trueVal, falseVal, value) {
      let valueToSet = value === true ? trueVal : falseVal;
      this.send('setAndBroadcast', path, valueToSet);
    },

    setAndBroadcastTtl(path, value) {
      const alwaysSendValue = path === 'expiry' || path === 'safetyBuffer';
      let valueToSet = value.enabled === true || alwaysSendValue ? `${value.seconds}s` : 0;
      this.send('setAndBroadcast', path, `${valueToSet}`);
    },

    codemirrorUpdated(path, isString, value, codemirror) {
      codemirror.performLint();
      const hasErrors = codemirror.state.lint.marked.length > 0;
      let valToSet = isString ? value : JSON.parse(value);

      if (!hasErrors) {
        this.model.set(path, valToSet);
        this.onChange(path, valToSet);
      }
    },

    toggleShow(path) {
      const value = !this.showInput;
      this.set('showInput', value);
      if (!value) {
        this.send('setAndBroadcast', path, null);
      }
    },
    handleKeyUp(name, value) {
      if (!this.onKeyUp) {
        return;
      }
      this.onKeyUp(name, value);
    },
  },
});
