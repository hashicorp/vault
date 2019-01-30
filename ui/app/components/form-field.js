import Component from '@ember/component';
import { computed } from '@ember/object';
import { capitalize } from 'vault/helpers/capitalize';
import { humanize } from 'vault/helpers/humanize';
import { dasherize } from 'vault/helpers/dasherize';

export default Component.extend({
  'data-test-field': true,
  classNames: ['field'],

  /*
   * @public Function
   * called whenever a value on the model changes via the component
   *
   */
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
   * this is usually derived from ember model `attributes` lookup,
   * and all members of `attr.options` are optional
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
    const label = this.get('attr.options.label');
    const name = this.get('attr.name');
    if (label) {
      return label;
    }
    if (name) {
      return capitalize([humanize([dasherize([name])])]);
    }
  }),

  // both the path to mutate on the model, and the path to read the value from
  /*
   * @private
   * @param string
   *
   * Computed property used to set values on the passed model
   *
   */
  valuePath: computed('attr.{name,options.fieldValue}', function() {
    return this.get('attr.options.fieldValue') || this.get('attr.name');
  }),

  /*
   *
   * @public
   * @param DS.Model
   *
   * the Ember Data model that `attr` is defined on
   */
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
      const path = this.get('valuePath');
      const { value } = keyFile;
      this.get('model').set(path, value);
      this.get('onChange')(path, value);
      this.set('file', keyFile);
    },

    setAndBroadcast(path, value) {
      this.get('model').set(path, value);
      this.get('onChange')(path, value);
    },

    setAndBroadcastBool(path, trueVal, falseVal, value) {
      let valueToSet = value === true ? trueVal : falseVal;
      this.send('setAndBroadcast', path, valueToSet);
    },

    codemirrorUpdated(path, isString, value, codemirror) {
      codemirror.performLint();
      const hasErrors = codemirror.state.lint.marked.length > 0;
      let valToSet = isString ? value : JSON.parse(value);

      if (!hasErrors) {
        this.get('model').set(path, valToSet);
        this.get('onChange')(path, valToSet);
      }
    },
  },
});
