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
 * @param attr=null {Object} - This is usually derived from ember model `attributes` lookup, and all members of `attr.options` are optional.
 * @param model=null {DS.Model} - The Ember Data model that `attr` is defined on
 * @param [disabled=false] {Boolean} - whether the field is disabled
 * @param [showHelpText=true] {Boolean} - whether to show the tooltip with help text from OpenAPI
 * @param [subText] {String} - Text to be displayed below the label
 *
 */

/*
  ATTR def:
  name: String
  type: One of: boolean, string, number, object
  options?: {
    editType: One of: boolean, optionalText, searchSelect, mountAccessor, kv, file, ttl, stringArray, json, textarea, password
  }
*/
export default Component.extend({
  layout,
  'data-test-field': true,
  classNames: ['field'],
  disabled: false,
  showHelpText: true,
  subText: '',

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
   *     editType: "ttl", // boolean, optionalText, searchSelect, mountAccessor, kv, file, ttl, stringArray, json, textarea, password
   *     helpText: "This will be in a tooltip",
   *     subText: "This will be under the label",
   *     warning: "This will be in stringlist, mountAccessor, kv, file",
   *     sensitive: true, // MaskedInput
   *     theme: "hashi", // editType json
   *     characterLimit: 30, // type string
   *     validationAttr: '?', // regular text input
   *     invalidMessage: 'Goes with validationAttr', // regular text input
   *     possibleValues: ['select', 'dropdown', 'options'],
   *     trueValue: 1, // overrides boolean true/false value
   *     falseValue: 0, // overrides boolean true/false value
   *     models: ['ssh-role'], // used for searchSelect lookup
   *     wildcardLabel: 'role', // used for searchSelect lookup
   *     subLabel: 'This goes in search select',
   *     fallbackComponent: 'some-component', // searchSelect
   *     selectLimit: 2, // limits how many in searchSelect can be chosen
   *     onlyAllowExisting: false, // do not allow create new on searchSelect
   *     setDefault: '45m', // used in ttl initial value
   *   },
   *   type: "boolean"
   * }
   *
   */
  attr: null,

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

  model: null,

  // This is only used internally for `optional-text` editType
  showInput: false,

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
  },
});
