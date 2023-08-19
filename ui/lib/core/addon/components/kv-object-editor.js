/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { isNone } from '@ember/utils';
import { assert } from '@ember/debug';
import { action } from '@ember/object';
import { guidFor } from '@ember/object/internals';
import KVObject from 'vault/lib/kv-object';

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
 * @param {boolean} [isMasked = false] - when true the <MaskedInput> renders instead of the default <textarea> to input the value portion of the key/value object 
 * @param {function} [onKeyUp] - function passed in that handles the dom keyup event. Used for validation on the kv custom metadata.
 * @param {string} [label] - label displayed over key value inputs
 * @param {string} [labelClass] - override default label class in FormFieldLabel component
 * @param {string} [warning] - warning that is displayed
 * @param {string} [helpText] - helper text. In tooltip.
 * @param {string} [subText] - placed under label.
 * @param {string} [keyPlaceholder] - placeholder for key input
 * @param {string} [valuePlaceholder] - placeholder for value input
 */

export default class KvObjectEditor extends Component {
  @tracked kvData;

  get placeholders() {
    return {
      key: this.args.keyPlaceholder || 'key',
      value: this.args.valuePlaceholder || 'value',
    };
  }
  get hasDuplicateKeys() {
    return this.kvData.uniqBy('name').length !== this.kvData.get('length');
  }

  // fired on did-insert from render modifier
  @action
  createKvData(elem, [value]) {
    this.kvData = KVObject.create({ content: [] }).fromJSON(value);
    this.addRow();
  }
  @action
  addRow() {
    if (!isNone(this.kvData.findBy('name', ''))) {
      return;
    }
    const newObj = { name: '', value: '' };
    guidFor(newObj);
    this.kvData.addObject(newObj);
  }
  @action
  updateRow() {
    this.args.onChange(this.kvData.toJSON());
  }
  @action
  deleteRow(object, index) {
    const oldObj = this.kvData.objectAt(index);
    assert('object guids match', guidFor(oldObj) === guidFor(object));
    this.kvData.removeAt(index);
    this.args.onChange(this.kvData.toJSON());
  }
  @action
  handleKeyUp(event) {
    if (this.args.onKeyUp) {
      this.args.onKeyUp(event.target.value);
    }
  }
}
