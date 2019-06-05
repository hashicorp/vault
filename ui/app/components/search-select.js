import Component from '@ember/component';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { computed } from '@ember/object';
import { singularize } from 'ember-inflector';

/**
 * @module SearchSelect
 * The `SearchSelect` is an implementation of the [ember-power-select-with-create](https://github.com/poteto/ember-cli-flash) used for form elements where options come dynamically from the API.
 * @example
 * <SearchSelect @id="group-policies" @models={{["policies/acl"]}} @onChange={{onChange}} @inputValue={{get model valuePath}} @helpText="Policies associated with this group" @label="Policies" @fallbackComponent="string-list" />
 *
 * @param id {String} - The name of the form field
 * @param models {String} - An array of model types to fetch from the API.
 * @param onChange {Func} - The onchange action for this form field.
 * @param inputValue {String} -  A comma-separated string or an array of strings.
 * @param [helpText] {String} - Text to be displayed in the info tooltip for this form field
 * @param label {String} - Label for this form field
 * @param fallbackComponent {String} - name of component to be rendered if the API call 403s
 *
 */
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
  formatOptions: function(options) {
    options = options.toArray().map(option => {
      option.searchText = `${option.name} ${option.id}`;
      return option;
    });
    let formattedOptions = this.selectedOptions.map(option => {
      let matchingOption = options.findBy('id', option);
      options.removeObject(matchingOption);
      return {
        id: option,
        name: matchingOption ? matchingOption.name : option,
        searchText: matchingOption ? matchingOption.searchText : option,
      };
    });
    this.set('selectedOptions', formattedOptions);
    if (this.options) {
      options = this.options.concat(options);
    }
    this.set('options', options);
  },
  fetchOptions: task(function*() {
    for (let modelType of this.models) {
      if (modelType.includes('identity')) {
        this.set('shouldRenderName', true);
      }
      try {
        let options = yield this.store.query(modelType, {});
        this.formatOptions(options);
      } catch (err) {
        if (err.httpStatus === 404) {
          //leave options alone, it's okay
          return;
        }
        if (err.httpStatus === 403) {
          this.set('shouldUseFallback', true);
          return;
        }
        //special case for storybook
        if (this.staticOptions) {
          let options = this.staticOptions;
          this.formatOptions(options);
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
    createOption(optionId) {
      let newOption = { name: optionId, id: optionId };
      this.selectedOptions.pushObject(newOption);
      this.handleChange();
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
    constructSuggestion(id) {
      return `Add new ${singularize(this.label)}: ${id}`;
    },
    hideCreateOptionOnSameID(id) {
      let existingOption = this.options.findBy('id', id) || this.options.findBy('name', id);
      return !existingOption;
    },
  },
});
