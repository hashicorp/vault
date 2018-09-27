import { isNone } from '@ember/utils';
import { assert } from '@ember/debug';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { guidFor } from '@ember/object/internals';
import KVObject from 'vault/lib/kv-object';

export default Component.extend({
  'data-test-component': 'kv-object-editor',
  classNames: ['field', 'form-section'],
  // public API
  // Ember Object to mutate
  value: null,
  label: null,
  helpText: null,
  // onChange will be called with the changed Value
  onChange() {},

  init() {
    this._super(...arguments);
    const data = KVObject.create({ content: [] }).fromJSON(this.get('value'));
    this.set('kvData', data);
    this.addRow();
  },

  kvData: null,

  kvDataAsJSON: computed('kvData', 'kvData.[]', function() {
    return this.get('kvData').toJSON();
  }),

  kvDataIsAdvanced: computed('kvData', 'kvData.[]', function() {
    return this.get('kvData').isAdvanced();
  }),

  kvHasDuplicateKeys: computed('kvData', 'kvData.@each.name', function() {
    let data = this.get('kvData');
    return data.uniqBy('name').length !== data.get('length');
  }),

  addRow() {
    let data = this.get('kvData');
    let newObj = { name: '', value: '' };
    if (!isNone(data.findBy('name', ''))) {
      return;
    }
    guidFor(newObj);
    data.addObject(newObj);
  },
  actions: {
    addRow() {
      this.addRow();
    },

    updateRow() {
      let data = this.get('kvData');
      this.get('onChange')(data.toJSON());
    },

    deleteRow(object, index) {
      let data = this.get('kvData');
      let oldObj = data.objectAt(index);

      assert('object guids match', guidFor(oldObj) === guidFor(object));
      data.removeAt(index);
      this.get('onChange')(data.toJSON());
    },
  },
});
