import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { singularize } from 'ember-inflector';
import { resolve } from 'rsvp';
import { filterOptions, defaultMatcher } from 'ember-power-select/utils/group-utils';

/**
 * @module SearchSelect
 * The `SearchSelect` is an implementation of the [ember-power-select](https://github.com/cibernox/ember-power-select) used for form elements where options come dynamically from the API.
 * @example
 * <SearchSelect @id="group-policies" @models={{["policies/acl"]}} @onChange={{onChange}} @selectLimit={{2}} @inputValue={{get model valuePath}} @helpText="Policies associated with this group" @label="Policies" @fallbackComponent="string-list" />
 *
 * @param {string} id - The name of the form field
 * @param {Array} model - model type to fetch from API
 * @param {function} onChange - The onchange action for this form field. ** SEE UTIL ** search-select-has-many.js if selecting models from a hasMany relationship
 * @param {string | Array} inputValue -  A comma-separated string or an array of strings -- array of ids for models.
 * @param {string} label - Label for this form field
 * @param {string} fallbackComponent - name of component to be rendered if the API call 403s
 * @param {boolean} [disallowNewItems=false] - Controls whether or not the user can add a new item if none found
 * @param {boolean} [passObject=false] - When true, the onChange callback returns an array of objects with id (string) and isNew (boolean)
 * @param {string} [helpText] - Text to be displayed in the info tooltip for this form field
 * @param {number} [selectLimit] - A number that sets the limit to how many select options they can choose
 * @param {string} [subText] - Text to be displayed below the label
 * @param {string} [subLabel] - a smaller label below the main Label
 * @param {string} [placeholder] - text you wish to replace the default "search" with
 * @param {boolean} [displayInherit] - if you need the search select component to display inherit instead of box.
 * @param {array} [removeOptions] - array of strings containing option names or model ids to filter from the dropdown (ex: 'allow_all)
 *
 * @param {Array} options - *Advanced usage* - `options` can be passed directly from the outside to the
 * power-select component. If doing this, `models` should not also be passed as that will overwrite the
 * passed value. ex: [{ name: 'namespace45', id: 'displayedName' }];
 * @param {function} search - *Advanced usage* - Customizes how the power-select component searches for matches -
 * see the power-select docs for more information.
 *
 */
export default class SearchSelectWithModal extends Component {
  @service store;

  @tracked allOptions = null; // list of options including matched
  @tracked selectedOptions = null; // list of selected options
  @tracked options = null; // all possible options
  @tracked showModal = false;
  @tracked newModelRecord = null;

  constructor() {
    super(...arguments);
    this.selectedOptions = this.inputValue;
  }

  get inputValue() {
    return this.args.inputValue || [];
  }

  get shouldRenderName() {
    return this.args.shouldRenderName || false;
  }

  get shouldUseFallback() {
    return this.args.shouldUseFallback || false;
  }

  get disallowNewItems() {
    return this.args.disallowNewItems || false;
  }

  // array of strings passed as an arg to remove from dropdown
  get removeOptions() {
    return this.args.removeOptions || null;
  }

  get passObject() {
    return this.args.passObject || false;
  }

  @action
  async fetchOptions() {
    let { oldOptions, options, selectedOptions } = this;
    let hasFormattedInput = typeof selectedOptions.firstObject !== 'string';
    if (options && !oldOptions && !hasFormattedInput) {
      // this is the first time they've been set, so we need to format them
      this.formatOptions(options);
    }
    this.oldOptions = options;
    if (!this.args.model) {
      if (this.options) {
        this.formatOptions(this.options);
      }
      return;
    }

    try {
      let queryOptions = {};
      let options = await this.store.query(this.args.model, queryOptions);
      this.formatOptions(options);
    } catch (err) {
      if (err.httpStatus === 404) {
        if (!this.options) {
          // If the call failed but the resource has items
          // from a different namespace, this allows the
          // selected items to display
          this.options = [];
        }

        return;
      }
      if (err.httpStatus === 403) {
        this.shouldUseFallback = true;
        return;
      }
      throw err;
    }
  }
  formatOptions(options) {
    options = options.toArray();
    if (this.removeOptions && this.removeOptions.length > 0) {
      options = options.filter((o) => !this.removeOptions.includes(o.id));
    }
    options = options.map((option) => {
      option.searchText = `${option.name} ${option.id}`;
      return option;
    });

    if (this.selectedOptions.length > 0) {
      this.selectedOptions = this.selectedOptions.map((option) => {
        let matchingOption = options.findBy('id', option);
        options.removeObject(matchingOption);
        return {
          id: option,
          name: matchingOption ? matchingOption.name : option,
          searchText: matchingOption ? matchingOption.searchText : option,
        };
      });
    }
    if (this.args.options) {
      if (this.removeOptions && this.removeOptions.length > 0) {
        options = this.options.filter((o) => !this.removeOptions.includes(o.id));
      }
      options = this.options.concat(options).uniq();
    }
    this.options = options;
  }

  handleChange() {
    if (this.selectedOptions.length && typeof this.selectedOptions.firstObject === 'object') {
      if (this.passObject) {
        this.args.onChange(
          Array.from(this.selectedOptions, (option) => ({ id: option.id, isNew: !!option.new }))
        );
      } else {
        this.args.onChange(Array.from(this.selectedOptions, (option) => option.id));
      }
    } else {
      this.args.onChange(this.selectedOptions);
    }
  }
  shouldShowCreate(id, options) {
    if (options && options.length && options.firstObject.groupName) {
      return !options.some((group) => group.options.findBy('id', id));
    }
    let existingOption = this.options && (this.options.findBy('id', id) || this.options.findBy('name', id));
    if (this.disallowNewItems && !existingOption) {
      return false;
    }
    return !existingOption;
  }
  //----- adapted from ember-power-select-with-create
  addCreateOption(term, results) {
    if (this.shouldShowCreate(term, results)) {
      const name = `Add new ${singularize(this.args.label)}: ${term}`;
      const suggestion = {
        __isSuggestion__: true,
        __value__: term,
        name,
        id: name,
      };
      results.unshift(suggestion);
    }
  }
  filter(options, searchText) {
    const matcher = (option, text) => defaultMatcher(option.searchText, text);
    return filterOptions(options || [], searchText, matcher);
  }
  // -----

  @action
  discardSelection(selected) {
    this.selectedOptions.removeObject(selected);
    // fire off getSelectedValue action higher up in get-credentials-card component
    if (!selected.new) {
      this.options.pushObject(selected);
    }
    this.handleChange();
  }
  // ----- adapted from ember-power-select-with-create
  @action
  searchAndSuggest(term, select) {
    if (term.length === 0) {
      return this.options;
    }
    if (this.search) {
      return resolve(this.search(term, select)).then((results) => {
        if (results.toArray) {
          results = results.toArray();
        }
        this.addCreateOption(term, results);
        return results;
      });
    }
    const newOptions = this.filter(this.options, term);
    this.addCreateOption(term, newOptions);
    return newOptions;
  }
  @action
  async selectOrCreate(selection) {
    if (selection && selection.__isSuggestion__) {
      const name = selection.__value__;
      this.showModal = true;
      let createRecord = await this.store.createRecord(this.args.model);
      createRecord.name = name;
      this.newModelRecord = createRecord;
    } else {
      this.selectedOptions.pushObject(selection);
      this.options.removeObject(selection);
    }
    this.handleChange();
  }
  // -----

  @action
  resetModal(model) {
    this.showModal = false;
    if (model && model.currentState.isSaved) {
      const { name } = model;
      this.selectedOptions.pushObject({ name, id: name });
    }
    this.newModelRecord = null;
  }
}
