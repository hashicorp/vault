/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import layout from '../templates/components/select';

/**
 * @module Select
 * Select components are used to render a dropdown.
 *
 * @example
 * ```js
 * <Select @label='Date Range' @options={{[{ value: 'berry', label: 'Berry' }]}} @onChange={{onChange}}/>
 * ```
 *
 * @param {string} [label=null] - The label for the select element.
 * @param {Array} [options=null] - A list of items that the user will select from. This can be an array of strings or objects.
 * @param {string} [selectedValue=null] - The currently selected item. Can also be used to set the default selected item. This should correspond to the `value` of one of the `<option>`s.
 * @param {string} [name = null] - The name of the select, used for the test selector.
 * @param {string} [valueAttribute = value]- When `options` is an array objects, the key to check for when assigning the option elements value.
 * @param {string} [labelAttribute = label] - When `options` is an array objects, the key to check for when assigning the option elements' inner text.
 * @param {boolean} [isInline = false] - Whether or not the select should be displayed as inline-block or block.
 * @param {boolean} [isFullwidth = false] - Whether or not the select should take up the full width of the parent element.
 * @param {boolean} [noDefault = false] - shows Select One with empty value as first option
 * @param {Func} [onChange] - The action to take once the user has selected an item. This method will be passed the `value` of the select.
 * @param {string} [ariaLabel] - pass when label is defined elsewhere to ensure the select input has a valid label
 */

export default Component.extend({
  layout,
  classNames: ['field'],
  label: null,
  selectedValue: null,
  name: null,
  options: null,
  valueAttribute: 'value',
  labelAttribute: 'label',
  isInline: false,
  isFullwidth: false,
  noDefault: false,
  onChange() {},
  ariaLabel: null,
});
