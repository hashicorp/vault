/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ArrayProxy from '@ember/array/proxy';
import Component from '@glimmer/component';
import autosize from 'autosize';
import { action } from '@ember/object';
import { set } from '@ember/object';
import { next } from '@ember/runloop';

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
 * @param {string} warning - Text displayed as a warning.
 * @param {string} helpText - Text displayed as a tooltip.
 * @param {string} type=array - Optional type for inputValue.
 * @param {string} attrName - We use this to check the type so we can modify the tooltip content.
 * @param {string} subText - Text below the label.
 */

export default class StringList extends Component {
  constructor() {
    super(...arguments);

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
    const inputs = this.inputList.filter((x) => x.value).mapBy('value');
    if (this.args.type === 'string') {
      return inputs.join(',');
    }
    return inputs;
  }

  get helpText() {
    if (this.args.attrName === 'tokenBoundCidrs') {
      return 'Specifies the blocks of IP addresses which are allowed to use the generated token. One entry per row.';
    } else {
      return this.args.helpText;
    }
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
    const inputObj = this.inputList.objectAt(idx);
    const onChange = this.args.onChange;
    set(inputObj, 'value', event.target.value);
    onChange(this.toVal());
  }

  @action
  addInput() {
    const inputList = this.inputList;
    if (inputList.get('lastObject.value') !== '') {
      inputList.pushObject({ value: '' });
    }
  }

  @action
  removeInput(idx) {
    const onChange = this.args.onChange;
    const inputs = this.inputList;
    inputs.removeObject(inputs.objectAt(idx));
    onChange(this.toVal());
  }
}
