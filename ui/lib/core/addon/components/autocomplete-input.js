/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module AutocompleteInput
 * AutocompleteInput components are used as standard string inputs or optionally select options to append to input value
 *
 * @example
 * <AutocompleteInput @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 *
 * @callback inputChangeCallback
 * @param {string} value - input value
 * @param {inputChangeCallback} onChange - fires when input value changes to mutate value param by caller
 * @param {string} [optionsTrigger] - display options dropdown when trigger character is input
 * @param {Object[]} [options] - array of { label, value } objects where label is displayed in options dropdown and value is appended to input value
 * @param {string} [label] - label to display above input
 * @param {string} [subText] - text to display below label
 * @param {string} [placeholder] - input placeholder
 */

export default class AutocompleteInput extends Component {
  dropdownAPI;
  inputElement;

  @action
  setElement(element) {
    this.inputElement = element.querySelector('.input');
  }
  @action
  setDropdownAPI(dropdownAPI) {
    this.dropdownAPI = dropdownAPI;
  }
  @action
  onInput(event) {
    const { options = [], optionsTrigger } = this.args;
    if (optionsTrigger && options.length) {
      const method = event.data === optionsTrigger ? 'open' : 'close';
      this.dropdownAPI.actions[method]();
    }
    this.args.onChange(event.target.value);
  }
  @action
  selectOption(value) {
    // if trigger character is at start of value it needs to be trimmed
    const appendValue = value.startsWith(this.args.optionsTrigger) ? value.slice(1) : value;
    const newValue = this.args.value + appendValue;
    this.args.onChange(newValue);
    this.dropdownAPI.actions.close();
    this.inputElement.focus();
  }
}
