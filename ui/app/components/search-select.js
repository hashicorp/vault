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
  shouldUseFallback: false,
  init() {
    this._super(...arguments);
    this.set('selectedOptions', this.inputValue || []);
  },
  fetchOptions: task(function*() {
    for (let modelType of this.models) {
      yield this.store
        .adapterFor(modelType)
        .query(null, { modelName: modelType })
        .then(resp => {
          let options = [];
          let data = resp.data;
          switch (modelType) {
            case 'identity/group':
            case 'identity/entity':
              data = data.key_info;
              Object.keys(data).forEach(id => {
                if (this.selectedOptions.includes(id)) {
                  this.selectedOptions.removeObject(id);
                  this.selectedOptions.addObject({ key: id, name: data[id].name });
                } else {
                  options.addObject({ key: id, name: data[id].name });
                }
              });
              break;
            default:
              options = data.keys;
              break;
          }
          options.removeObjects(this.selectedOptions);
          if (this.options) {
            options = this.options.concat(options);
          }
          this.set('options', options);
        })
        .catch(err => {
          if (err.httpStatus === 403) {
            this.set('shouldUseFallback', true);
          }
        });
    }
  }).on('didInsertElement'),
  handleChange() {
    if (this.selectedOptions.length && typeof this.selectedOptions.firstObject === 'object') {
      this.onChange(Array.from(this.selectedOptions, option => option.key));
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
