import ArrayProxy from '@ember/array/proxy';
import Component from '@ember/component';
import { set, computed } from '@ember/object';

export default Component.extend({
  'data-test-component': 'string-list',
  classNames: ['field', 'string-list', 'form-section'],

  /*
   * @public
   * @param String
   *
   * Optional - Text displayed in the header above all of the inputs
   *
   */
  label: null,

  /*
   * @public
   * @param Function
   *
   * Function called when any of the inputs change
   * accepts a single param `value` that is the
   * result of calling `toVal()`.
   *
   */
  onChange: () => {},

  /*
   * @public
   * @param String | Array
   * A comma-separated string or an array of strings.
   * Defaults to an empty array.
   *
   */
  inputValue: computed(function() {
    return [];
  }),

  /*
    *
    * @public
    * @param String - ['array'|'string]
    *
    * Optional type for `inputValue` - defaults to `'array'`
    * Needs to match type of `inputValue` because it is set by the component on init.
    *
    */
  type: 'array',

  /*
    *
    * @private
    * @param Ember.ArrayProxy
    *
    * mutable array that contains objects in the form of
    * {
    *   value: 'somestring',
    * }
    *
    * used to track the state of values bound to the various inputs
    *
    */
  inputList: computed(function() {
    return ArrayProxy.create({
      content: [],
      // trim the `value` when accessing objects
      objectAtContent: function(idx) {
        const obj = this.get('content').objectAt(idx);
        if (obj && obj.value) {
          set(obj, 'value', obj.value.trim());
        }
        return obj;
      },
    });
  }),

  init() {
    this._super(...arguments);
    this.setType();
    this.toList();
    this.send('addInput');
  },

  setType() {
    const list = this.inputList;
    if (!list) {
      return;
    }
    this.set('type', typeof list);
  },

  toVal() {
    const inputs = this.inputList.filter(x => x.value).mapBy('value');
    if (this.get('format') === 'string') {
      return inputs.join(',');
    }
    return inputs;
  },

  toList() {
    let input = this.inputValue || [];
    const inputList = this.inputList;
    if (typeof input === 'string') {
      input = input.split(',');
    }
    inputList.addObjects(input.map(value => ({ value })));
  },

  actions: {
    inputChanged(idx, val) {
      const inputObj = this.inputList.objectAt(idx);
      const onChange = this.onChange;
      set(inputObj, 'value', val);
      onChange(this.toVal());
    },

    addInput() {
      const inputList = this.inputList;
      if (inputList.get('lastObject.value') !== '') {
        inputList.pushObject({ value: '' });
      }
    },

    removeInput(idx) {
      const onChange = this.onChange;
      const inputs = this.inputList;
      inputs.removeObject(inputs.objectAt(idx));
      onChange(this.toVal());
    },
  },
});
