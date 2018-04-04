import Ember from 'ember';

const { computed, set } = Ember;

export default Ember.Component.extend({
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
  inputValue: [],

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
    return Ember.ArrayProxy.create({
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
    const list = this.get('inputList');
    if (!list) {
      return;
    }
    this.set('type', typeof list);
  },

  toVal() {
    const inputs = this.get('inputList').filter(x => x.value).mapBy('value');
    if (this.get('format') === 'string') {
      return inputs.join(',');
    }
    return inputs;
  },

  toList() {
    let input = this.get('inputValue') || [];
    const inputList = this.get('inputList');
    if (typeof input === 'string') {
      input = input.split(',');
    }
    inputList.addObjects(input.map(value => ({ value })));
  },

  actions: {
    inputChanged(idx, val) {
      const inputObj = this.get('inputList').objectAt(idx);
      const onChange = this.get('onChange');
      set(inputObj, 'value', val);
      onChange(this.toVal());
    },

    addInput() {
      const inputList = this.get('inputList');
      if (inputList.get('lastObject.value') !== '') {
        inputList.pushObject({ value: '' });
      }
    },

    removeInput(idx) {
      const onChange = this.get('onChange');
      const inputs = this.get('inputList');
      inputs.removeObject(inputs.objectAt(idx));
      onChange(this.toVal());
    },
  },
});
