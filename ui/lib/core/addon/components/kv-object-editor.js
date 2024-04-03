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
 * @param {boolean} [isSingleRow = false] - when true the kv object editor will only show one row and hide the Add button
 * @param {function} [onKeyUp] - function passed in that handles the dom keyup event. Used for validation on the kv custom metadata.
 * @param {string} [label] - label displayed over key value inputs
 * @param {string} [labelClass] - override default label class in FormFieldLabel component
 * @param {string} [warning] - warning that is displayed
 * @param {string} [helpText] - helper text. In tooltip.
 * @param {string} [subText] - placed under label.
 * @param {string} [keyPlaceholder] - placeholder for key input
 * @param {string} [valuePlaceholder] - placeholder for value input
 * @param {boolean} [allowWhiteSpace = false] - when true, allows whitespace in the key input
 * @param {boolean} [warnNonStringValues = false] - when true, shows a warning if the value is a non-string
 */

export default class KvObjectEditor extends Component {
  // kvData is type ArrayProxy, so addObject etc are fine here
  @tracked kvData;

  get placeholders() {
    return {
      key: this.args.keyPlaceholder || 'key',
      value: this.args.valuePlaceholder || 'value',
    };
  }
  get hasDuplicateKeys() {
    return this.kvData.uniqBy('name').length !== this.kvData.length;
  }

  // fired on did-insert from render modifier
  @action
  createKvData(elem, [value]) {
    this.kvData = KVObject.create({ content: [] }).fromJSON(value);

    if (!this.args.isSingleRow || !value || Object.keys(value).length < 1) {
      this.addRow();
    }
  }
  @action
  addRow() {
    if (!isNone(this.kvData.find((datum) => datum.name === ''))) {
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
  showWhitespaceWarning = (name) => {
    if (this.args.allowWhiteSpace) return false;
    return new RegExp('\\s', 'g').test(name);
  };
  showNonStringWarning = (value) => {
    if (!this.args.warnNonStringValues) return false;
    try {
      JSON.parse(value);
      return true;
    } catch (e) {
      return false;
    }
  };
}
