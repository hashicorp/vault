/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * @module RegexValidator
 * RegexValidator components are used to provide input forms for regex values, along with a toggle-able validation input which does not get saved to the model.
 *
 * @example
 * ```js
 * const attrExample = {
 *    name: 'valName',
 *    options: {
 *      helpText: 'Shows in tooltip',
 *      subText: 'Shows underneath label',
 *      docLink: 'Adds docs link to subText if present',
 *      defaultValue: 'hello', // Shows if no value on model
 *    }
 * }
 * <RegexValidator @onChange={action 'myAction'} @attr={attrExample} @labelString="Label String" @value="initial value" />
 * ```
 * @param {string} value - the value of the main input which will be updated in onChange
 * @param {func} [onChange] - the action that should trigger when pattern input is changed. Required when attr is provided.
 * @param {string} [labelString] - Form label. Anticipated from form-field. Required when attr is provided.
 * @param {object} [attr] - attribute from model. Anticipated from form-field. Example of attribute shape above. When not provided pattern input is hidden
 * @param {string} [testInputLabel] - label for test input
 * @param {string} [testInputSubText] - sub text for test input
 * @param {boolean} [showGroups] - show groupings based on pattern and test input
 * @param {func} [onValidate] - action triggered every time the test string is validated against the regex -- passes testValue and captureGroups
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class RegexValidator extends Component {
  @tracked testValue = '';
  @tracked showTestValue = false;

  constructor() {
    super(...arguments);
    this.showTestValue = !this.args.attr;
  }

  get testInputLabel() {
    return this.args.testInputLabel || 'Test string';
  }
  get regex() {
    return new RegExp(this.args.value, 'g');
  }
  get regexError() {
    const testString = this.testValue;
    if (!testString || !this.args.value) return false;
    const matchArray = testString.toString().match(this.regex);
    if (this.args.onValidate) {
      this.args.onValidate(this.testValue, this.captureGroups);
    }
    return testString !== matchArray?.join('');
  }
  get captureGroups() {
    const result = this.regex.exec(this.testValue);
    if (result) {
      // first item is full string match but we are only interested in the captured groups
      const [fullMatch, ...matches] = result; // eslint-disable-line
      const groups = matches.map((m, index) => ({ position: `$${index + 1}`, value: m }));
      // push named capture groups into array -> eg (<lastFour>\d{4})
      if (result.groups) {
        for (const key in result.groups) {
          groups.push({ position: `$${key}`, value: result.groups[key] });
        }
      }
      return groups;
    }
    return [];
  }

  @action
  updateTestValue(evt) {
    this.testValue = evt.target.value;
  }

  @action
  toggleTestValue() {
    this.showTestValue = !this.showTestValue;
  }
}
