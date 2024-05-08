/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { resolve } from 'rsvp';
import { filterOptions, defaultMatcher } from 'ember-power-select/utils/group-utils';
import { removeFromArray } from 'vault/helpers/remove-from-array';
import { addToArray } from 'vault/helpers/add-to-array';
import { assert } from '@ember/debug';
/**
 * @module SearchSelect
 * The `SearchSelect` is an implementation of the [ember-power-select](https://github.com/cibernox/ember-power-select) used for form elements where options come dynamically from the API.
 * @example
 *  <SearchSelect
 *    @id="policy"
 *    @models={{array "policies/acl"}}
 *    @onChange={{this.onChange}}
 *    @inputValue={{get @model this.valuePath}}
 *    @wildcardLabel="role"
 *    @fallbackComponent="string-list"
 *    @selectLimit={{1}}
 *    @backend={{@model.backend}}
 *    @disallowNewItems={{true}}
 *    class={{if this.validationError "dropdown-has-error-border"}}
 * />
 *
 // * component functionality
 * @param {function} onChange - The onchange action for this form field. ** SEE EXAMPLE ** mfa-login-enforcement-form.js (onMethodChange) for example when selecting models from a hasMany relationship
 * @param {array} [inputValue] - Array of strings corresponding to the input's initial value, e.g. an array of model ids that on edit will appear as selected items below the input
 * @param {boolean} [disallowNewItems=false] - Controls whether or not the user can add a new item if none found
 * @param {boolean} [shouldRenderName=false] - By default an item's id renders in the dropdown, `true` displays the name with its id in smaller text beside it *NOTE: the boolean flips automatically with 'identity' models or if this.idKey !== 'id'
 * @param {string} [nameKey="name"] - if shouldRenderName=true, you can use this arg to specify which key to use for the rendered name. Defaults to "name".
 * @param {array} [parentManageSelected] - Array of selected items if the parent is keeping track of selections, see mfa-login-enforcement-form.js
 * @param {boolean} [passObject=false] - When true, the onChange callback returns an array of objects with id (string) and isNew (boolean) (and any params from objectKeys). By default - onChange returns an array of id strings.
 * @param {array} [objectKeys] - Array of values that correlate to model attrs. Used to render attr other than 'id' beside the name if shouldRenderName=true. If passObject=true, objectKeys are added to the passed, selected object.
 * @param {number} [selectLimit] - Sets select limit

// * query params for dropdown items
 * @param {Array} models - An array of model types to fetch from the API.
 * @param {string} [backend] - name of the backend if the query for options needs additional information (eg. secret backend)
 * @param {object} [queryObject] - object passed as query options to this.store.query(). NOTE: will override @backend

 // * template only/display args
 * @param {string} id - The name of the form field
 * @param {string} [label] - Label for this form field
 * @param {string} [labelClass] - overwrite default label size (14px) from class="is-label"
 * @param {string} [ariaLabel] - fallback accessible label if label is not provided
 * @param {string} [subText] - Text to be displayed below the label
 * @param {string} fallbackComponent - name of component to be rendered if the API call 403s
 * @param {string} [helpText] - Text to be displayed in the info tooltip for this form field
 * @param {string} [wildcardLabel] - string (singular) for rendering label tag beside a wildcard selection (i.e. 'role*'), for the number of items it includes, e.g. @wildcardLabel="role" -> "includes 4 roles"
 * @param {string} [placeholder] - text you wish to replace the default "search" with
 * @param {boolean} [displayInherit=false] - if you need the search select component to display inherit instead of box.
 * @param {boolean} [renderInPlace] - pass `true` when power select renders in a modal
 * @param {function} [renderInfoTooltip] - receives each inputValue string and list of dropdownOptions as args, so parent can determine when to render a tooltip beside a selectedOption and the tooltip text. see 'oidc/provider-form.js'
 * @param {boolean} [disabled] - if true sets the disabled property on the ember-power-select component and makes it unusable.
 *
 // * advanced customization
 * @param {Array} options - array of objects passed directly to the power-select component. If doing this, `models` should not also be passed as that will overwrite the
 * passed options. ex: [{ name: 'namespace45', id: 'displayedName' }]. It's recommended the parent should manage the array of selected items if manually passing in options.
 * @param {function} search - Customizes how the power-select component searches for matches - see the power-select docs for more information.
 *
 */

export default class SearchSelect extends Component {
  @service store;
  @tracked shouldUseFallback = false;
  @tracked selectedOptions = []; // array of selected options (initially set by @inputValue)
  @tracked dropdownOptions = []; // options that will render in dropdown, updates as selections are added/discarded
  @tracked allOptions = []; // both selected and unselected options, used for wildcard filter

  constructor() {
    super(...arguments);
    assert(
      'one of @id, @label, or @ariaLabel must be passed to search-select component',
      this.args.id || this.args.label || this.args.ariaLabel
    );
  }

  get hidePowerSelect() {
    return this.selectedOptions.length >= this.args.selectLimit;
  }

  get idKey() {
    // if objectKeys exists, use the first element of the array as the identifier
    // make 'id' as the first element in objectKeys if you do not want to override the default of 'id'
    return this.args.objectKeys ? this.args.objectKeys[0] : 'id';
  }

  get shouldRenderName() {
    return this.args.models?.some((model) => model.includes('identity')) ||
      this.idKey !== 'id' ||
      this.args.shouldRenderName
      ? true
      : false;
  }

  get nameKey() {
    return this.args.nameKey || 'name';
  }

  get searchEnabled() {
    if (typeof this.args.searchEnabled === 'boolean') return this.args.searchEnabled;
    return true;
  }

  addSearchText(optionsToFormat) {
    // maps over array of objects or response from query
    return optionsToFormat.map((option) => {
      const id = option[this.idKey] ? option[this.idKey] : option.id;
      option.searchText = `${option[this.nameKey]} ${id}`;
      return option;
    });
  }

  formatInputAndUpdateDropdown(inputValues) {
    // inputValues are initially an array of strings from @inputValue
    // map over so selectedOptions are objects
    return inputValues.map((option) => {
      const matchingOption = this.dropdownOptions.find((opt) => opt[this.idKey] === option);
      // tooltip text comes from return of parent function
      const addTooltip = this.args.renderInfoTooltip
        ? this.args.renderInfoTooltip(option, this.dropdownOptions)
        : false;

      // remove any matches from dropdown list
      this.dropdownOptions = removeFromArray(this.dropdownOptions, matchingOption);
      return {
        id: option,
        name: matchingOption ? matchingOption[this.nameKey] : option,
        searchText: matchingOption ? matchingOption.searchText : option,
        addTooltip,
        // add additional attrs if we're using a dynamic idKey
        ...(this.idKey !== 'id' && this.customizeObject(matchingOption)),
      };
    });
  }

  @task
  *fetchOptions() {
    this.dropdownOptions = []; // reset dropdown anytime we re-fetch

    if (this.args.parentManageSelected) {
      // works in tandem with parent passing in @options directly
      this.selectedOptions = this.args.parentManageSelected;
    }

    if (!this.args.models) {
      if (Array.isArray(this.args.options)) {
        const { options } = this.args;
        // if options are nested, let parent handle formatting - see path-filter-config-list.js
        this.dropdownOptions = options.some((e) => Object.keys(e).includes('groupName'))
          ? options
          : [...this.addSearchText(options)];

        if (!this.args.parentManageSelected) {
          //  set selectedOptions and remove matches from dropdown list
          this.selectedOptions = this.args.inputValue
            ? this.formatInputAndUpdateDropdown(this.args.inputValue)
            : [];
        }
      }
      return;
    }

    for (const modelType of this.args.models) {
      try {
        let queryParams = {};
        if (this.args.backend) {
          queryParams = { backend: this.args.backend };
        }
        if (this.args.queryObject) {
          queryParams = this.args.queryObject;
        }
        // fetch options from the store
        const options = yield this.store.query(modelType, queryParams);

        // store both select + unselected options in tracked property used by wildcard filter
        this.allOptions = [...this.allOptions, ...options.map((option) => option.id)];

        // add to dropdown options
        this.dropdownOptions = [...this.dropdownOptions, ...this.addSearchText(options)];
      } catch (err) {
        if (err.httpStatus === 404) {
          // continue to query other models even if one 404s
          // and so selectedOptions will be set after for loop
          continue;
        }
        if (err.httpStatus === 403) {
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
    if (this.selectedOptions.length && typeof this.selectedOptions[0] === 'object') {
      this.args.onChange(
        Array.from(this.selectedOptions, (option) =>
          this.args.passObject ? this.customizeObject(option) : option.id
        )
      );
    } else {
      this.args.onChange(this.selectedOptions);
    }
  }

  shouldShowCreate(id, searchResults) {
    if (searchResults && searchResults.length && searchResults[0].groupName) {
      return !searchResults.some((group) => group.options.find((opt) => opt.id === id));
    }
    const existingOption =
      this.dropdownOptions && this.dropdownOptions.find((opt) => opt.id === id || opt.name === id);
    if (this.args.disallowNewItems && !existingOption) {
      return false;
    }
    return !existingOption;
  }

  // ----- adapted from ember-power-select-with-create
  addCreateOption(term, results) {
    if (this.shouldShowCreate(term, results)) {
      const name = `Click to add new item: ${term}`;
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

    let additionalKeys;
    if (this.args.objectKeys) {
      // pull attrs corresponding to objectKeys from model record, add to the selection
      additionalKeys = Object.fromEntries(this.args.objectKeys.map((key) => [key, option[key]]));
      // filter any undefined attrs, which could mean the model was not hydrated,
      // the record is new or the model doesn't have that attribute
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

  @action
  discardSelection(selected) {
    this.selectedOptions = removeFromArray(this.selectedOptions, selected);
    if (!selected.new) {
      this.dropdownOptions = addToArray(this.dropdownOptions, selected);
    }
    this.handleChange();
  }

  // ----- adapted from ember-power-select-with-create
  @action
  searchAndSuggest(term, select) {
    if (term.length === 0) {
      return this.dropdownOptions;
    }
    if (this.args.search) {
      return resolve(this.args.search(term, select)).then((results) => {
        this.addCreateOption(term, results);
        return results;
      });
    }
    const newOptions = this.filter(this.dropdownOptions, term);
    this.addCreateOption(term, newOptions);
    return newOptions;
  }

  @action
  selectOrCreate(selection) {
    if (selection && selection.__isSuggestion__) {
      const name = selection.__value__;
      this.selectedOptions = addToArray(this.selectedOptions, { name, id: name, new: true });
    } else {
      this.selectedOptions = addToArray(this.selectedOptions, selection);
      this.dropdownOptions = removeFromArray(this.dropdownOptions, selection);
    }
    this.handleChange();
  }

  // -----
}
