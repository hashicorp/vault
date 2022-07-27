import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { singularize } from 'ember-inflector';
import { resolve } from 'rsvp';
import { filterOptions, defaultMatcher } from 'ember-power-select/utils/group-utils';

/**
 * @module SearchSelectWithModal
 * The `SearchSelectWithModal` is an implementation of the [ember-power-select](https://github.com/cibernox/ember-power-select) used for form elements where options come dynamically from the API. It can only accept a single model.
 * It renders a passed in form component so records can be created inline, via a modal that pops up after clicking "Create new <id>" from the dropdown menu.
 * **!! NOTE: any form passed must be able to receive an @onSave and @onCancel arg so that the modal will close properly. See `oidc/client-form.hbs` that renders a modal for the `oidc/assignment-form.hbs` as an example.
 * @example
 * <SearchSelectWithModal
 *         @id="assignments"
 *         @model="oidc/assignment"
 *         @label="assignment name"
 *         @labelClass="is-label"
 *         @subText="Search for an existing assignment, or type a new name to create it."
 *         @inputValue={{map-by "id" @model.assignments}}
 *         @onChange={{this.handleSearchSelect}}
 *         {{! since this is the "limited" radio select option we do not want to include 'allow_all' }}
 *         @excludeOptions={{array "allow_all"}}
 *         @fallbackComponent="string-list"
 *         @modalFormComponent="oidc/assignment-form"
 *         @modalSubtext="Use assignment to specify which Vault entities and groups are allowed to authenticate."
 *       />
 *
 * @param {string} id - the model's attribute for the form field, will be interpolated into create new text: `Create new ${singularize(this.args.id)}`
 * @param {Array} model - model type to fetch from API (can only be a single model)
 * @param {string} label - Label that appears above the form field
 * @param {string} [helpText] - Text to be displayed in the info tooltip for this form field
 * @param {string} [subText] - Text to be displayed below the label
 * @param {string} [placeholder] - placeholder text to override the default text of "Search"
 * @param {function} onChange - The onchange action for this form field. ** SEE UTIL ** search-select-has-many.js if selecting models from a hasMany relationship
 * @param {string | Array} inputValue -  A comma-separated string or an array of strings -- array of ids for models.
 * @param {string} fallbackComponent - name of component to be rendered if the API returns a 403s
 * @param {boolean} [passObject=false] - When true, the onChange callback returns an array of objects with id (string) and isNew (boolean)
 * @param {number} [selectLimit] - A number that sets the limit to how many select options they can choose
 * @param {array} [excludeOptions] - array of strings containing model ids to filter from the dropdown (ex: ['allow_all'])
 * @param {function} search - *Advanced usage* - Customizes how the power-select component searches for matches -
 * see the power-select docs for more information.
 *
 */
export default class SearchSelectWithModal extends Component {
  @service store;

  @tracked selectedOptions = null; // list of selected options
  @tracked allOptions = null; // all possible options
  @tracked showModal = false;
  @tracked newModelRecord = null;
  @tracked shouldUseFallback = false;

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

  get excludeOptions() {
    return this.args.excludeOptions || null;
  }

  get passObject() {
    return this.args.passObject || false;
  }

  @action
  async fetchOptions() {
    try {
      let queryOptions = {};
      let options = await this.store.query(this.args.model, queryOptions);
      this.formatOptions(options);
    } catch (err) {
      if (err.httpStatus === 404) {
        if (!this.allOptions) {
          // If the call failed but the resource has items
          // from a different namespace, this allows the
          // selected items to display
          this.allOptions = [];
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
    if (this.excludeOptions) {
      options = options.filter((o) => !this.excludeOptions.includes(o.id));
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
    this.allOptions = options;
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
    let existingOption =
      this.allOptions && (this.allOptions.findBy('id', id) || this.allOptions.findBy('name', id));
    return !existingOption;
  }
  //----- adapted from ember-power-select-with-create
  addCreateOption(term, results) {
    if (this.shouldShowCreate(term, results)) {
      const name = `Create new ${singularize(this.args.id)}: ${term}`;
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
    this.allOptions.pushObject(selected);
    this.handleChange();
  }
  // ----- adapted from ember-power-select-with-create
  @action
  searchAndSuggest(term, select) {
    if (term.length === 0) {
      return this.allOptions;
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
    const newOptions = this.filter(this.allOptions, term);
    this.addCreateOption(term, newOptions);
    return newOptions;
  }
  @action
  async selectOrCreate(selection) {
    // if creating we call handleChange in the resetModal action to ensure the model is valid and successfully created
    // before adding it to the DOM (and parent model)
    // if just selecting, then we handleChange immediately
    if (selection && selection.__isSuggestion__) {
      const name = selection.__value__;
      this.showModal = true;
      let createRecord = await this.store.createRecord(this.args.model);
      createRecord.name = name;
      this.newModelRecord = createRecord;
    } else {
      this.selectedOptions.pushObject(selection);
      this.allOptions.removeObject(selection);
      this.handleChange();
    }
  }
  // -----

  @action
  resetModal(model) {
    this.showModal = false;
    if (model && model.currentState.isSaved) {
      const { name } = model;
      this.selectedOptions.pushObject({ name, id: name });
      this.handleChange();
    }
    this.newModelRecord = null;
  }
}
