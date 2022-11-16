import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import { filterOptions, defaultMatcher } from 'ember-power-select/utils/group-utils';

/**
 * @module SearchSelectWithModal
 * The `SearchSelectWithModal` is an implementation of the [ember-power-select](https://github.com/cibernox/ember-power-select) used for form elements where options come dynamically from the API. It can only accept a single model.
 * It renders a passed in form component so records can be created inline, via a modal that pops up after clicking "Create new <id>" from the dropdown menu.
 * **!! NOTE: any form passed must be able to receive an @onSave and @onCancel arg so that the modal will close properly. See `oidc/client-form.hbs` that renders a modal for the `oidc/assignment-form.hbs` as an example.
 * @example
 * <SearchSelectWithModal
 *   @id="assignments"
 *   @models={{array "oidc/assignment"}}
 *   @label="assignment name"
 *   @subText="Search for an existing assignment, or type a new name to create it."
 *   @inputValue={{map-by "id" @model.assignments}}
 *   @onChange={{this.handleSearchSelect}}
 *   {{! since this is the "limited" radio select option we do not want to include 'allow_all' }}
 *   @excludeOptions={{array "allow_all"}}
 *   @fallbackComponent="string-list"
 *   @modalFormComponent="oidc/assignment-form"
 *   @modalSubtext="Use assignment to specify which Vault entities and groups are allowed to authenticate."
 * />
 *
 // * component functionality
 * @param {function} onChange - The onchange action for this form field. ** SEE UTIL ** search-select-has-many.js if selecting models from a hasMany relationship
 * @param {array} [inputValue] - Array of strings corresponding to the input's initial value, e.g. an array of model ids that on edit will appear as selected items below the input
 * @param {boolean} [shouldRenderName=false] - By default an item's id renders in the dropdown, `true` displays the name with its id in smaller text beside it *NOTE: the boolean flips automatically with 'identity' models
 * @param {number} [selectLimit] - Sets select limit
 * @param {array} [excludeOptions] - array of strings containing model ids to filter from the dropdown (ex: ['allow_all'])
 * @param {function} onCreate -  callback fired with search input so parent can handle creating the model
 * @param {object} formModel - model created by parent's onCreate callback, passed to form's @model argument
 
// * query params for dropdown items
 * @param {array} models - model type to fetch from API (can only be a single model)
 * @param {string} [backend] - name of the backend if the query for options needs additional information (eg. secret backend)
 * @param {object} [queryObject] - object passed as query options to this.store.query(). NOTE: will override @backend
 
 // * template only/display args
 * @param {string} id - The name of the form field
 * @param {string} [label] - Label appears above the form field
 * @param {string} [labelClass] - overwrite default label size (14px) from class="is-label"
 * @param {string} [helpText] - Text to be displayed in the info tooltip for this form field
 * @param {string} [subText] - Text to be displayed below the label
 * @param {string} fallbackComponent - name of component to be rendered if the API call 403s
 * @param {string} [placeholder] - placeholder text to override the default text of "Search"
 * @param {boolean} [displayInherit=false] - if you need the search select component to display inherit instead of box.
 */
export default class SearchSelectWithModal extends Component {
  @service store;
  @tracked shouldUseFallback = false;

  @tracked selectedOptions = []; // list of selected options
  @tracked dropdownOptions = []; // options that will render in dropdown, updates as selections are added/discarded
  @tracked showModal = false;
  @tracked nameInput = null;

  get hidePowerSelect() {
    return this.selectedOptions.length >= this.args.selectLimit;
  }

  get shouldRenderName() {
    return this.args.models?.some((model) => model.includes('identity')) || this.args.shouldRenderName
      ? true
      : false;
  }

  addSearchText(optionsToFormat) {
    // maps over array models from query
    return optionsToFormat.toArray().map((option) => {
      option.searchText = `${option.name} ${option.id}`;
      return option;
    });
  }

  formatInputAndUpdateDropdown(inputValues) {
    // inputValues are initially an array of strings from @inputValue
    // map over so selectedOptions are objects
    return inputValues.map((option) => {
      const matchingOption = this.dropdownOptions.findBy('id', option);
      // remove any matches from dropdown list
      this.dropdownOptions.removeObject(matchingOption);
      return {
        id: option,
        name: matchingOption ? matchingOption.name : option,
        searchText: matchingOption ? matchingOption.searchText : option,
      };
    });
  }

  @task
  *fetchOptions() {
    this.dropdownOptions = []; // reset dropdown anytime we re-fetch
    if (!this.args.models) {
      return;
    }

    for (const modelType of this.args.models) {
      try {
        const queryParams = {};
        // fetch options from the store
        let options = yield this.store.query(modelType, queryParams);
        if (this.args.excludeOptions) {
          options = options.filter((o) => !this.args.excludeOptions.includes(o.id));
        }
        // add to dropdown options
        this.dropdownOptions = [...this.dropdownOptions, ...this.addSearchText(options)];
      } catch (err) {
        if (err.httpStatus === 404) {
          // continue to query other models even if one 404s
          // and so selectedOptions will be set after for loop
          continue;
        }
        if (err.httpStatus === 403) {
          if (this.dropdownOptions && this.args.models.length > 1) continue;
          this.shouldUseFallback = true;
          return;
        }
        throw err;
      }
    }

    // after all models are queried, set selectedOptions and remove matches from dropdown list
    this.selectedOptions = this.args.inputValue
      ? this.formatInputAndUpdateDropdown(this.args.inputValue)
      : [];
  }

  @action
  handleChange() {
    if (this.selectedOptions.length && typeof this.selectedOptions.firstObject === 'object') {
      this.args.onChange(Array.from(this.selectedOptions, (option) => option.id));
    } else {
      this.args.onChange(this.selectedOptions);
    }
  }

  shouldShowCreate(id, searchResults) {
    if (searchResults && searchResults.length && searchResults.firstObject.groupName) {
      return !searchResults.some((group) => group.options.findBy('id', id));
    }
    const existingOption =
      this.dropdownOptions &&
      (this.dropdownOptions.findBy('id', id) || this.dropdownOptions.findBy('name', id));
    return !existingOption;
  }

  // ----- adapted from ember-power-select-with-create
  addCreateOption(term, results) {
    if (this.shouldShowCreate(term, results)) {
      const name = `No results found for "${term}". Click here to create it.`;
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
    this.dropdownOptions.pushObject(selected);
    this.handleChange();
  }

  // ----- adapted from ember-power-select-with-create
  @action
  searchAndSuggest(term) {
    if (term.length === 0) {
      return this.dropdownOptions;
    }
    if (this.args.models?.some((model) => model.includes('policy'))) {
      term = term.toLowerCase();
    }
    const newOptions = this.filter(this.dropdownOptions, term);
    this.addCreateOption(term, newOptions);
    return newOptions;
  }

  @action
  selectOrCreate(selection) {
    if (selection && selection.__isSuggestion__) {
      // user has clicked to create a new item
      // wait to handleChange below in resetModal
      const name = selection.__value__;
      this.showModal = true;
      this.args.onCreate({ name });
      // if firing onCreate does not generate a formModel, then pass name to form component
      // where `onCreate` is fired by user interaction (ex: selecting a policy type)
      if (!this.args.formModel) this.nameInput = name;
    } else {
      // user has selected an existing item, handleChange immediately
      this.selectedOptions.pushObject(selection);
      this.dropdownOptions.removeObject(selection);
      this.handleChange();
    }
  }
  // -----

  @action
  resetModal(model) {
    // resetModal fires when the form component calls onSave or onCancel
    this.showModal = false;
    if (model && model.currentState.isSaved) {
      const { name } = model;
      this.selectedOptions.pushObject({ name, id: name });
      this.handleChange();
    }
    this.nameInput = null;
  }
}
