import Component from '@ember/component';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { computed } from '@ember/object';
import { singularize } from 'ember-inflector';
import { resolve } from 'rsvp';
import { filterOptions, defaultMatcher } from 'ember-power-select/utils/group-utils';
import layout from '../templates/components/search-select';

/**
 * @module SearchSelect
 * The `SearchSelect` is an implementation of the [ember-power-select](https://github.com/cibernox/ember-power-select) used for form elements where options come dynamically from the API.
 * @example
 * <SearchSelect @id="group-policies" @models={{["policies/acl"]}} @onChange={{onChange}} @selectLimit={{2}} @inputValue={{get model valuePath}} @helpText="Policies associated with this group" @label="Policies" @fallbackComponent="string-list" />
 *
 * @param {string} id - The name of the form field
 * @param {Array} models - An array of model types to fetch from the API.
 * @param {object} queryObject - query object to pass as query options to as second argument to this.store.query(modelType, {})
 * @param {function} onChange - The onchange action for this form field. ** SEE UTIL ** search-select-has-many.js if selecting models from a hasMany relationship
 * @param {string | Array} inputValue -  A comma-separated string or an array of strings -- array of ids for models.
 * @param {string} label - Label for this form field
 * @param {string} fallbackComponent - name of component to be rendered if the API call 403s
 * @param {string} [backend] - name of the backend if the query for options needs additional information (eg. secret backend)
 * @param {boolean} [disallowNewItems=false] - Controls whether or not the user can add a new item if none found
 * @param {boolean} [passObject=false] - When true, the onChange callback returns an array of objects with id (string) and isNew (boolean) (instead of an array of id strings)
 * @param {array} [objectKeys=null] - Array of values that correlate to model attrs. When passObject=true, objectKeys are added to the passed object. NOTE: make 'id' as the first element in objectKeys if you do not want to override the default of 'id'
 * @param {string} [helpText] - Text to be displayed in the info tooltip for this form field
 * @param {number} [selectLimit] - A number that sets the limit to how many select options they can choose
 * @param {string} [subText] - Text to be displayed below the label
 * @param {string} [subLabel] - a smaller label below the main Label
 * @param {string} [wildcardLabel] - when you want the searchSelect component to return a count on the model for options returned when using a wildcard you must provide a label of the count e.g. role.  Should be singular.
 * @param {string} [placeholder] - text you wish to replace the default "search" with
 * @param {boolean} [displayInherit] - if you need the search select component to display inherit instead of box.
 *
 * @param {Array} options - *Advanced usage* - `options` can be passed directly from the outside to the
 * power-select component. If doing this, `models` should not also be passed as that will overwrite the
 * passed value. ex: [{ name: 'namespace45', id: 'displayedName' }];
 * @param {function} search - *Advanced usage* - Customizes how the power-select component searches for matches -
 * see the power-select docs for more information.
 *
 */
export default Component.extend({
  layout,
  'data-test-component': 'search-select',
  attributeBindings: ['data-test-component'],
  classNameBindings: ['displayInherit:display-inherit'],
  classNames: ['field', 'search-select'],
  store: service(),
  onChange: () => {},
  inputValue: computed(function () {
    return [];
  }),
  allOptions: null, // list of options including matched
  selectedOptions: null, // list of selected options
  options: null, // all possible options
  queryObject: null,
  shouldUseFallback: false,
  shouldRenderName: false,
  disallowNewItems: false,
  passObject: false,
  objectKeys: null,
  idKey: computed('objectKeys', function () {
    // if objectKeys exists, then use the first element of the array as the identifier
    return this.objectKeys ? this.objectKeys[0] : 'id';
  }),
  init() {
    this._super(...arguments);
    this.set('selectedOptions', this.inputValue || []);
  },
  didRender() {
    this._super(...arguments);
    let { oldOptions, options, selectedOptions } = this;
    let hasFormattedInput = typeof selectedOptions.firstObject !== 'string';
    if (options && !oldOptions && !hasFormattedInput) {
      // this is the first time they've been set, so we need to format them
      this.formatOptions(options);
    }
    this.set('oldOptions', options);
  },
  formatOptions: function (options) {
    options = options.toArray().map((option) => {
      option.searchText = `${option.name} ${option[this.idKey]}`;
      return option;
    });
    let allOptions = options.toArray().map((option) => {
      return option.id;
    });
    this.set('allOptions', allOptions); // used by filter-wildcard helper
    let formattedOptions = this.selectedOptions.map((option) => {
      let matchingOption = options.findBy(this.idKey, option);
      options.removeObject(matchingOption);
      return {
        id: option,
        name: matchingOption ? matchingOption.name : option,
        searchText: matchingOption ? matchingOption.searchText : option,
        // conditionally spread configured object if we're using the dynamic idKey
        ...(this.idKey !== 'id' && this.customizeObject(matchingOption)),
      };
    });
    this.set('selectedOptions', formattedOptions);
    if (this.options) {
      options = this.options.concat(options).uniq();
    }
    this.set('options', options);
  },
  fetchOptions: task(function* () {
    if (!this.models) {
      if (this.options) {
        this.formatOptions(this.options);
      }
      return;
    }
    if (this.idKey !== 'id') {
      // if passing a dynamic idKey, then display it in the dropdown beside the name
      this.set('shouldRenderName', true);
    }
    for (let modelType of this.models) {
      if (modelType.includes('identity')) {
        this.set('shouldRenderName', true);
      }
      try {
        let queryOptions = {};
        if (this.backend) {
          queryOptions = { backend: this.backend };
        }
        if (this.queryObject) {
          queryOptions = this.queryObject;
        }
        let options = yield this.store.query(modelType, queryOptions);
        console.log(options, 'OPTIIONS');
        this.formatOptions(options);
      } catch (err) {
        if (err.httpStatus === 404) {
          if (!this.options) {
            // If the call failed but the resource has items
            // from a different namespace, this allows the
            // selected items to display
            this.set('options', []);
          }

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
      this.onChange(Array.from(this.selectedOptions, (option) => this.customizeObject(option)));
    } else {
      this.onChange(this.selectedOptions);
    }
  },
  shouldShowCreate(id, options) {
    if (options && options.length && options.firstObject.groupName) {
      return !options.some((group) => group.options.findBy('id', id));
    }
    let existingOption = this.options && (this.options.findBy('id', id) || this.options.findBy('name', id));
    if (this.disallowNewItems && !existingOption) {
      return false;
    }
    return !existingOption;
  },
  //----- adapted from ember-power-select-with-create
  addCreateOption(term, results) {
    if (this.shouldShowCreate(term, results)) {
      const name = `Add new ${singularize(this.label || 'item')}: ${term}`;
      const suggestion = {
        __isSuggestion__: true,
        __value__: term,
        name,
        id: name,
      };
      results.unshift(suggestion);
    }
  },
  filter(options, searchText) {
    const matcher = (option, text) => defaultMatcher(option.searchText, text);
    return filterOptions(options || [], searchText, matcher);
  },
  // -----
  customizeObject(option) {
    // if passObject=true return object, otherwise return string of option id
    if (this.passObject) {
      let additionalKeys;
      if (this.objectKeys) {
        // pull attrs corresponding to objectKeys from model record, add to the selected option (object) and send to the parent
        additionalKeys = Object.fromEntries(this.objectKeys.map((key) => [key, option[key]]));
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
  },
  actions: {
    onChange(val) {
      this.onChange(val);
    },
    discardSelection(selected) {
      this.selectedOptions.removeObject(selected);
      // fire off getSelectedValue action higher up in get-credentials-card component
      if (!selected.new) {
        this.options.pushObject(selected);
      }
      this.handleChange();
    },
    // ----- adapted from ember-power-select-with-create
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
    },
    selectOrCreate(selection) {
      if (selection && selection.__isSuggestion__) {
        const name = selection.__value__;
        this.selectedOptions.pushObject({ name, id: name, new: true });
      } else {
        this.selectedOptions.pushObject(selection);
        this.options.removeObject(selection);
      }
      this.handleChange();
    },
    // -----
  },
});
