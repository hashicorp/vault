/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ArrayProxy from '@ember/array/proxy';
import Component from '@glimmer/component';
import autosize from 'autosize';
import { action } from '@ember/object';
import { set } from '@ember/object';
import { next } from '@ember/runloop';
import { tracked } from '@glimmer/tracking';
import { addToArray } from 'vault/helpers/add-to-array';
import { removeFromArray } from 'vault/helpers/remove-from-array';

/**
 * @module StringList
 *
 * @example
 * ```js
 * <StringList @label={label} @onChange={{this.setAndBroadcast}} @inputValue={{this.valuePath}}/>
 * ```
 * @param {string} label - Text displayed in the header above all the inputs.
 * @param {function} onChange - Function called when any of the inputs change.
 * @param {string} inputValue - A string or an array of strings.
 * @param {string} helpText - Text displayed as a tooltip.
 * @param {string} type=array - Optional type for inputValue.
 * @param {string} attrName - We use this to check the type so we can modify the tooltip content.
 * @param {string} subText - Text below the label.
 */

export default class StringList extends Component {
  @tracked indicesWithComma = [];

  constructor() {
    super(...arguments);

    // inputList is type ArrayProxy, so addObject etc are fine here
    this.inputList = ArrayProxy.create({
      // trim the `value` when accessing objects
      content: [],
      objectAtContent: function (idx) {
        const obj = this.content.objectAt(idx);
        if (obj && obj.value) {
          set(obj, 'value', obj.value.trim());
        }
        return obj;
      },
    });
    this.type = this.args.type || 'array';
    this.setType();
    next(() => {
      this.toList();
      this.addInput();
    });
  }

  setType() {
    const list = this.inputList;
    if (!list) {
      return;
    }
    this.type = typeof list;
  }

  toList() {
    let input = this.args.inputValue || [];
    const inputList = this.inputList;
    if (typeof input === 'string') {
      input = input.split(',');
    }
    inputList.addObjects(input.map((value) => ({ value })));
  }

  toVal() {
    const inputs = this.inputList.filter((x) => x.value).map((x) => x.value);
    if (this.args.type === 'string') {
      return inputs.join(',');
    }
    return inputs;
  }

  @action
  autoSize(element) {
    autosize(element.querySelector('textarea'));
  }

  @action
  autoSizeUpdate(element) {
    autosize.update(element.querySelector('textarea'));
  }

  @action
  inputChanged(idx, event) {
    if (event.target.value.includes(',') && !this.indicesWithComma.includes(idx)) {
      this.indicesWithComma = addToArray(this.indicesWithComma, idx);
    }
    if (!event.target.value.includes(',')) {
      this.indicesWithComma = removeFromArray(this.indicesWithComma, idx);
    }

    const inputObj = this.inputList.objectAt(idx);
    set(inputObj, 'value', event.target.value);
    this.args.onChange(this.toVal());
  }

  @action
  addInput() {
    const [lastItem] = this.inputList.slice(-1);
    if (lastItem?.value !== '') {
      this.inputList.pushObject({ value: '' });
    }
  }

  @action
  removeInput(idx) {
    const itemToRemove = this.inputList.objectAt(idx);
    this.inputList.removeObject(itemToRemove);
    this.args.onChange(this.toVal());
  }
}
