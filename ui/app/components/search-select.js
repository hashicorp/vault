import Component from '@ember/component';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { computed } from '@ember/object';

export default Component.extend({
  'data-test-component': 'search-select',
  classNames: ['field', 'search-select'],
  store: service(),

  /*
   * @public
   * @param Function
   *
   * Function called when any of the inputs change
   * accepts a single param `value`
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
  selectedOptions: null, //list of selected options
  options: null, //all possible options
  shouldUseFallback: false,
  shouldRenderName: false,
  init() {
    this._super(...arguments);
    this.set('selectedOptions', this.inputValue || []);
  },
  fetchOptions: task(function*() {
    for (let modelType of this.models) {
      if (modelType.includes('identity')) {
        this.set('shouldRenderName', true);
      }
      try {
        let options = yield this.store.query(modelType, {});
        options = options.toArray().map(option => {
          option.searchText = `${option.name} ${option.id}`;
          return option;
        });
        let formattedOptions = this.selectedOptions.map(option => {
          let matchingOption = options.findBy('id', option);
          options.removeObject(matchingOption);
          return { id: option, name: matchingOption.name, searchText: matchingOption.searchText };
        });
        this.set('selectedOptions', formattedOptions);
        if (this.options) {
          options = this.options.concat(options);
        }
        this.set('options', options);
      } catch (err) {
        if (err.httpStatus === 404) {
          //leave options alone, it's okay
          return;
        }
        if (err.httpStatus === 403) {
          this.set('shouldUseFallback', true);
          return;
        }
        throw err;
      }
    }
  }).on('didInsertElement'),
  handleChange() {
    if (this.selectedOptions.length && typeof this.selectedOptions.firstObject === 'object') {
      this.onChange(Array.from(this.selectedOptions, option => option.id));
    } else {
      this.onChange(this.selectedOptions);
    }
  },
  actions: {
    onChange(val) {
      this.onChange(val);
    },
    selectOption(option) {
      this.selectedOptions.pushObject(option);
      this.options.removeObject(option);
      this.handleChange();
    },
    discardSelection(selected) {
      this.selectedOptions.removeObject(selected);
      this.options.pushObject(selected);
      this.handleChange();
    },
  },
});
