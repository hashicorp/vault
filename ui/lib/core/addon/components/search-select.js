import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { singularize } from 'ember-inflector';
import { resolve } from 'rsvp';
import { filterOptions, defaultMatcher } from 'ember-power-select/utils/group-utils';
import { isWildcardString } from 'vault/helpers/is-wildcard-string';
import { assert } from '@ember/debug';
/**
 * @module SearchSelect
 * The `SearchSelect` is an implementation of the [ember-power-select](https://github.com/cibernox/ember-power-select) used for form elements where options come dynamically from the API.
 * @example
 *  <SearchSelect
 *    id="policy"
 *    models={{array "policies/acl"}}
 *    onChange={{this.onChange}}
 *    inputValue={{get @model this.valuePath}}
 *    wildcardLabel="role"
 *    fallbackComponent="string-list"
 *    selectLimit={{1}}
 *    backend={{@model.backend}}
 *    disallowNewItems={{true}}
 *    class={{if this.validationError "dropdown-has-error-border"}}
 * />
 * 
 // * component functionality
 * @param {function} onChange - The onchange action for this form field. ** SEE UTIL ** search-select-has-many.js if selecting models from a hasMany relationship
 * @param {array} [inputValue] - Array of strings corresponding to the input's initial value, e.g. an array of model ids that on edit will appear as selected items below the input
 * @param {boolean} [disallowNewItems=false] - Controls whether or not the user can add a new item if none found
 * @param {boolean} [shouldRenderName=false] - By default an item's id renders in the dropdown, `true` displays the name with its id in smaller text beside it *NOTE: the boolean flips automatically with 'identity' models or if this.idKey !== 'id'
 * @param {boolean} [resetSelection=false] - When true, the component will clear its input after "Add" is clicked.
 * @param {boolean} [passObject=false] - When true, the onChange callback returns an array of objects with id (string) and isNew (boolean) (and any params from objectKeys). By default - onChange returns an array of id strings. 
 * @param {array} [objectKeys=null] - Array of values that correlate to model attrs. When passObject=true, objectKeys are added to the passed, selected object. *NOTE: make 'id' as the first element in objectKeys if you do not want to override the default of 'id'
 * @param {number} [selectLimit] - Sets select limit
 
// * query params for dropdown items
 * @param {Array} models - An array of model types to fetch from the API.
 * @param {string} [backend] - name of the backend if the query for options needs additional information (eg. secret backend)
 * @param {object} [queryObject] - object passed as query options to this.store.query(). NOTE: will override @backend
 
 // * template only/display args
 * @param {string} id - The name of the form field
 * @param {string} [label] - Label for this form field, also used to generate `Add new ${singularize(label}` copy when adding a nonexisting option
 * @param {string} [labelClass] - overwrite default label size (14px) - NOTE within <FormField> @labelClass is "title is-4"
 * @param {string} [subText] - Text to be displayed below the label
 * @param {string} fallbackComponent - name of component to be rendered if the API call 403s
 * @param {string} [helpText] - Text to be displayed in the info tooltip for this form field
 * @param {string} [wildcardLabel] - string (singular) for rendering label tag beside a wildcard selection (i.e. 'role*'), for the number of items it includes, e.g. @wildcardLabel="role" -> "includes 4 roles"
 * @param {string} [placeholder] - text you wish to replace the default "search" with
 * @param {boolean} [displayInherit=false] - if you need the search select component to display inherit instead of box.
 * @param {boolean} [renderInfoTooltip=false] - if you want search select to render a tooltip beside a selected item if no corresponding model was returned from .query
 *
 // * advanced customization
 * @param {Array} options - array of objects passed directly to the power-select component. If doing this, `models` should not also be passed as that will overwrite the
 * passed options. ex: [{ name: 'namespace45', id: 'displayedName' }];
 * @param {function} search - Customizes how the power-select component searches for matches - see the power-select docs for more information.
 *
 */

export default class SearchSelect extends Component {
  @service store;
  @tracked selectedOptions; // array of selected options, initially array of strings, then array of objects **WHYY/WHEN?!**
  @tracked allOptions; // array id strings for ALL existing options, including matched
  @tracked searchSelectOptions; // list of options rendered in dropdown
  @tracked shouldUseFallback = false;

  get disallowNewItems() {
    return this.args.disallowNewItems || false;
  }

  get hidePowerSelect() {
    return this.selectedOptions.length >= this.args.selectLimit;
  }

  get idKey() {
    if (this.args.objectKeys) {
      assert('@objectKeys passed to <SearchSelect> must be an array', Array.isArray(this.args.objectKeys));
    }
    // if objectKeys exists, use the first element of the array as the identifier
    return this.args.objectKeys ? this.args.objectKeys[0] : 'id';
  }

  get passObject() {
    return this.args.passObject || false;
  }

  get resetSelection() {
    return this.args.resetSelection || false;
  }

  get shouldRenderName() {
    return this.args.models?.some((model) => model.includes('identity')) ||
      this.idKey !== 'id' ||
      this.args.shouldRenderName
      ? true
      : false;
  }

  constructor() {
    super(...arguments);
    this.selectedOptions = this.args.inputValue || [];
  }

  // 'format' means add searchText
  formatOptions(optionsToFormat) {
    //* `optionsToFormat` - array of objects or response from query (model Class)
    let formattedDropdownOptions = optionsToFormat.toArray().map((option) => {
      option.searchText = `${option.name} ${option[this.idKey]}`;
      return option;
    });
    this.allOptions = formattedDropdownOptions.mapBy('id'); // used by filter-wildcard helper
    return formattedDropdownOptions;
  }

  // remove any selected items from dropdown and configure selected object returned to parent
  updateDropdownOptions(selectedOptions) {
    // iterate selectedOptions, initially set by inputValue, and format
    // remove any already selected items from array of ALL options (formattedDropdownOptions)
    this.selectedOptions = selectedOptions.map((option) => {
      let matchingOption = this.searchSelectOptions.findBy(this.idKey, option);
      // an undefined matchingOption means a selectedOption, on edit, didn't match a model returned from the query
      let addTooltip = matchingOption || isWildcardString([option]) ? false : true; // add tooltip to let user know the selection may not exist
      this.searchSelectOptions.removeObject(matchingOption);
      return {
        id: option,
        name: matchingOption ? matchingOption.name : option,
        searchText: matchingOption ? matchingOption.searchText : option,
        addTooltip,
        // conditionally spread configured object if we're using the dynamic idKey
        ...(this.idKey !== 'id' && this.customizeObject(matchingOption)),
      };
    });
  }

  @task
  *fetchOptions() {
    if (this.args.resetSelection) {
      this.selectedOptions = [];
    }
    this.searchSelectOptions = [];

    if (!this.args.models) {
      if (this.args.options) {
        if (this.args.search) {
          this.searchSelectOptions = this.args.options;
        } else {
          this.searchSelectOptions = [...this.formatOptions(this.args.options)];
          this.updateDropdownOptions(this.selectedOptions);
        }
      }
      return;
    }
    for (let modelType of this.args.models) {
      try {
        let queryOptions = {};
        if (this.args.backend) {
          queryOptions = { backend: this.args.backend };
        }
        if (this.args.queryObject) {
          queryOptions = this.args.queryObject;
        }
        let options = yield this.store.query(modelType, queryOptions);
        this.searchSelectOptions = [...this.searchSelectOptions, ...this.formatOptions(options)];
        this.updateDropdownOptions(this.selectedOptions); // remove options from inputValue, format selected options
      } catch (err) {
        if (err.httpStatus === 404) {
          if (!this.args.options) {
            // If the call failed but the resource has items
            // from a different namespace, this allows the
            // selected items to display
            return;
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
  }

  @action
  handleChange() {
    if (this.selectedOptions.length && typeof this.selectedOptions.firstObject === 'object') {
      this.onChange(Array.from(this.selectedOptions, (option) => this.customizeObject(option)));
    } else {
      this.onChange(this.selectedOptions);
    }
  }

  shouldShowCreate(id, dropdownOptions) {
    if (dropdownOptions && dropdownOptions.length && dropdownOptions.firstObject.groupName) {
      return !dropdownOptions.some((group) => group.options.findBy('id', id));
    }
    let existingOption =
      this.searchSelectOptions &&
      (this.searchSelectOptions.findBy('id', id) || this.searchSelectOptions.findBy('name', id));
    if (this.disallowNewItems && !existingOption) {
      return false;
    }
    return !existingOption;
  }

  //----- adapted from ember-power-select-with-create
  addCreateOption(term, results) {
    if (this.shouldShowCreate(term, results)) {
      const name = `Add new ${singularize(this.args.label.toLowerCase() || 'item')}: ${term}`;
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
  customizeObject(option) {
    if (!option) return;
    // if passObject=true return object, otherwise return string of option id
    if (this.passObject) {
      let additionalKeys;
      if (this.args.objectKeys) {
        // pull attrs corresponding to objectKeys from model record, add to the selected option (object) and send to the parent
        additionalKeys = Object.fromEntries(this.args.objectKeys.map((key) => [key, option[key]]));
        // filter any undefined attrs, which means the model did not have a value for that attr
        // no value could mean the model was not hydrated, the record is new or the model doesn't have that attribute
        Object.keys(additionalKeys).forEach((key) => {
          if (additionalKeys[key] === undefined) {
            delete additionalKeys[key];
          }
        });
      }
      return {
        id: option.id,
        isNew: !!option.new,
        ...additionalKeys,
      };
    }
    return option.id;
  }

  @action
  onChange(val) {
    this.args.onChange(val);
  }

  @action
  discardSelection(selected) {
    this.selectedOptions.removeObject(selected);
    // fire off getSelectedValue action higher up in get-credentials-card component
    if (!selected.new) {
      this.searchSelectOptions.pushObject(selected);
    }
    this.handleChange();
  }

  // ----- adapted from ember-power-select-with-create
  @action
  searchAndSuggest(term, select) {
    if (term.length === 0) {
      return this.searchSelectOptions;
    }
    if (this.args.search) {
      return resolve(this.args.search(term, select)).then((results) => {
        if (results.toArray) {
          results = results.toArray();
        }
        this.addCreateOption(term, results);
        return results;
      });
    }
    const newOptions = this.filter(this.searchSelectOptions, term);
    this.addCreateOption(term, newOptions);
    return newOptions;
  }

  @action
  selectOrCreate(selection) {
    if (selection && selection.__isSuggestion__) {
      const name = selection.__value__;
      this.selectedOptions.pushObject({ name, id: name, new: true });
    } else {
      this.selectedOptions.pushObject(selection);
      this.searchSelectOptions.removeObject(selection);
    }
    this.handleChange();
  }
  // -----
}
