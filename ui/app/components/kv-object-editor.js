/**
 * @module KvObjectEditor
 * KvObjectEditor components are called in FormFields when the editType on the model is kv.  They are used to show a key-value input field.
 *
 * @example
 * ```js
 * <KvObjectEditor
 *  @value={{get model valuePath}}
 *  @onChange={{action "setAndBroadcast" valuePath }}
 *  @label="some label"
   />
 * ```
 * @param {string} value - the value is captured from the model.
 * @param {function} onChange - function that captures the value on change
 * @param {function} onKeyUp - function passed in that handles the dom keyup event. Used for validation on the kv custom metadata.
 * @param {string} [label] - label displayed over key value inputs
 * @param {string} [warning] - warning that is displayed
 * @param {string} [helpText] - helper text. In tooltip.
 * @param {string} [subText] - placed under label.
 * @param {boolean} [small-label]- change label size.
 * @param {boolean} [formSection] - if false the component is meant to live outside of a form, like in the customMetadata which is nested already inside a form-section.
 */

import { isNone } from '@ember/utils';
import { assert } from '@ember/debug';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { guidFor } from '@ember/object/internals';
import KVObject from 'vault/lib/kv-object';

export default Component.extend({
  'data-test-component': 'kv-object-editor',
  classNames: ['field'],
  classNameBindings: ['formSection:form-section'],
  formSection: true,
  // public API
  // Ember Object to mutate
  value: null,
  label: null,
  helpText: null,
  subText: null,
  // onChange will be called with the changed Value
  onChange() {},

  init() {
    this._super(...arguments);
    const data = KVObject.create({ content: [] }).fromJSON(this.value);
    this.set('kvData', data);
    this.addRow();
  },

  kvData: null,

  kvDataAsJSON: computed('kvData', 'kvData.[]', function() {
    return this.kvData.toJSON();
  }),

  kvDataIsAdvanced: computed('kvData', 'kvData.[]', function() {
    return this.kvData.isAdvanced();
  }),

  kvHasDuplicateKeys: computed('kvData', 'kvData.@each.name', function() {
    let data = this.kvData;
    return data.uniqBy('name').length !== data.get('length');
  }),

  addRow() {
    let data = this.kvData;
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
      let data = this.kvData;
      this.onChange(data.toJSON());
    },

    deleteRow(object, index) {
      let data = this.kvData;
      let oldObj = data.objectAt(index);

      assert('object guids match', guidFor(oldObj) === guidFor(object));
      data.removeAt(index);
      this.onChange(data.toJSON());
    },

    handleKeyUp(name, value) {
      if (!this.onKeyUp) {
        return;
      }
      this.onKeyUp(name, value);
    },
  },
});
