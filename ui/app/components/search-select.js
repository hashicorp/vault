import Component from '@ember/component';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { computed } from '@ember/object';

export default Component.extend({
  'data-test-component': 'search-select',
  classNames: ['field', 'search-select', 'form-section'],
  store: service(),
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
  selectedOption: null,
  selectedOptions: null,
  options: null,
  init() {
    this._super(...arguments);
    this.set('selectedOptions', this.inputValue);
  },
  fetchOptions: task(function*() {
    yield this.store
      .adapterFor(this.modelType)
      .query(null, { modelName: this.modelType }, { findAll: true })
      .then(resp => {
        let options = resp.data.keys;
        options.removeObjects(this.selectedOptions);
        this.set('options', resp.data.keys);
      });
  }).on('didInsertElement'),

  actions: {
    selectOption(option) {
      this.selectedOptions.pushObject(option);
      this.options.removeObject(option);
      if (!this.isList) {
        this.set('selectedOption', option);
      }
      this.onChange(this.selectedOptions);
    },
    discardSelection(selected) {
      this.selectedOptions.removeObject(selected);
      this.options.pushObject(selected);
      this.onChange(this.selectedOptions);
    },
  },
});
